package consumer

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"
	"watchAlert/alert/mute"
	"watchAlert/alert/process"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"

	"github.com/zeromicro/go-zero/core/logc"
	"golang.org/x/sync/errgroup"
)

const (
	// 任务通道缓冲区大小
	TaskChannelBufferSize = 1

	// 默认处理时间间隔，2s
	DefaultProcessTime = 2

	// 状态前缀
	RecoverStatePrefix = "Recover_"
	FiringStatePrefix  = "Firing_"

	// 消费者 panic 重启策略
	ConsumerRestartMaxAttempts = 5
	ConsumerRestartBaseDelay   = time.Second
	ConsumerRestartMaxDelay    = 30 * time.Second
	ConsumerRestartResetAfter  = time.Minute
)

type consumerTask struct {
	cancel context.CancelFunc
	done   chan struct{}
}

type (
	ConsumeInterface interface {
		Submit(faultCenter models.FaultCenter)
		Stop(faultCenterId string)
		Watch(ctx context.Context, faultCenter models.FaultCenter)
		RestartAllConsumers()
		StopAllConsumers()
	}

	Consume struct {
		ctx *ctx.Context
		sync.RWMutex
		taskMu        sync.Mutex
		tasks         map[string]*consumerTask
		restartPolicy consumerRestartPolicy
	}

	EventsGroup struct {
		NoticeID string // 通知组 ID
		Events   []*models.AlertCurEvent
	}

	RulesGroup struct {
		RuleID string // 规则组 ID
		Groups map[string]EventsGroup
	}

	AlertGroups struct {
		Rules map[string]RulesGroup
		lock  sync.RWMutex
	}
)

// AddAlert 添加告警
func (ag *AlertGroups) AddAlert(stateId string, alert *models.AlertCurEvent, faultCenter models.FaultCenter) {
	ag.lock.Lock()
	defer ag.lock.Unlock()

	// 获取通知对象 ID 列表 用于事件分组
	noticeObjIds := ag.getNoticeId(alert, faultCenter)
	if len(noticeObjIds) == 0 {
		return // 如果没有通知对象ID，则跳过
	}

	for _, noticeObjId := range noticeObjIds {
		// 检查 Rule 是否存在
		rule, exists := ag.Rules[stateId]
		if !exists {
			// 创建新的 RuleGroup
			ag.Rules[stateId] = RulesGroup{
				RuleID: stateId,
				Groups: make(map[string]EventsGroup),
			}
			rule = ag.Rules[stateId]
		}

		// 检查 Group 是否存在
		group, groupExists := rule.Groups[noticeObjId]
		if !groupExists {
			// 创建新的 EventsGroup
			rule.Groups[noticeObjId] = EventsGroup{
				NoticeID: noticeObjId,
				Events:   []*models.AlertCurEvent{},
			}
			group = rule.Groups[noticeObjId]
		}

		// 添加事件到对应组
		group.Events = append(group.Events, alert)
		// 更新 group 映射
		ag.Rules[stateId].Groups[noticeObjId] = group
	}
}

// getNoticeId 从告警路由中获取该事件匹配的通知对象
func (ag *AlertGroups) getNoticeId(alert *models.AlertCurEvent, faultCenter models.FaultCenter) []string {
	if len(faultCenter.NoticeRoutes) > 0 {
		labels := alert.Labels

		for _, route := range faultCenter.NoticeRoutes {
			if evalCondition(labels, route.NoticeLabels) {
				return route.NoticeIds
			}
		}
	}

	return faultCenter.NoticeIds
}

func NewConsumerEngine(ctx *ctx.Context) ConsumeInterface {
	return &Consume{
		ctx:   ctx,
		tasks: make(map[string]*consumerTask),
	}
}

func (c *Consume) Submit(faultCenter models.FaultCenter) {
	c.taskMu.Lock()
	defer c.taskMu.Unlock()

	c.stopLocked(faultCenter.ID)
	c.startLocked(faultCenter)
}

func (c *Consume) startLocked(faultCenter models.FaultCenter) {
	watchCtx, cancel := context.WithCancel(c.ctx.Ctx)
	task := &consumerTask{cancel: cancel, done: make(chan struct{})}
	c.tasks[faultCenter.ID] = task
	go func() {
		shouldRestart := false
		defer func() {
			if r := recover(); r != nil {
				if watchCtx.Err() == nil {
					stack := debug.Stack()
					logc.Error(c.ctx.Ctx, fmt.Sprintf("Recovered from consumer watch goroutine panic: %s, FaultCenterName: %s, Id: %s\n%s", r, faultCenter.Name, faultCenter.ID, stack))
					shouldRestart = true
				}
			}

			close(task.done)
			if shouldRestart {
				c.restartAfterPanic(watchCtx, faultCenter, task)
			}
		}()
		c.Watch(watchCtx, faultCenter)
	}()
}

func (c *Consume) stopLocked(faultCenterID string) {
	task, exists := c.tasks[faultCenterID]
	if !exists {
		return
	}

	delete(c.tasks, faultCenterID)
	task.cancel()
	<-task.done
}

func (c *Consume) stopWithoutWaitLocked(faultCenterID string) {
	if task, exists := c.tasks[faultCenterID]; exists {
		delete(c.tasks, faultCenterID)
		task.cancel()
	}
}

func (c *Consume) Stop(faultCenterId string) {
	c.taskMu.Lock()
	defer c.taskMu.Unlock()
	c.stopLocked(faultCenterId)
}

func (c *Consume) Restart(faultCenter models.FaultCenter) {
	// Panic recovery runs from the Consumer task itself. Waiting here would wait
	// for the current task, so cancel the old task without waiting.
	c.taskMu.Lock()
	defer c.taskMu.Unlock()

	c.stopWithoutWaitLocked(faultCenter.ID)
	c.startLocked(faultCenter)
}

func (c *Consume) restartAfterPanic(watchCtx context.Context, faultCenter models.FaultCenter, task *consumerTask) {
	if watchCtx.Err() != nil {
		return
	}

	delay, ok := c.restartPolicy.next(faultCenter.ID, time.Now())
	if !ok {
		logc.Errorf(c.ctx.Ctx, "consumer restart limit reached, FaultCenterName: %s, Id: %s", faultCenter.Name, faultCenter.ID)
		return
	}

	logc.Errorf(c.ctx.Ctx, "consumer will restart after %s, FaultCenterName: %s, Id: %s", delay, faultCenter.Name, faultCenter.ID)
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-timer.C:
		c.taskMu.Lock()
		defer c.taskMu.Unlock()
		if watchCtx.Err() != nil {
			return
		}
		current, exists := c.tasks[faultCenter.ID]
		if !exists || current != task {
			return
		}
		c.stopWithoutWaitLocked(faultCenter.ID)
		c.startLocked(faultCenter)
	case <-watchCtx.Done():
		return
	}
}

// Watch 启动 Consumer Watch 进程
func (c *Consume) Watch(ctx context.Context, faultCenter models.FaultCenter) {
	taskChan := make(chan struct{}, TaskChannelBufferSize)
	timer := time.NewTicker(time.Second * time.Duration(DefaultProcessTime))
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			// 处理任务信号量
			taskChan <- struct{}{}
			c.executeTask(faultCenter, taskChan)
		case <-ctx.Done():
			return
		}
	}
}

// executeTask 执行具体的任务逻辑
func (c *Consume) executeTask(faultCenter models.FaultCenter, taskChan chan struct{}) {
	defer func() {
		// 释放任务信号量
		<-taskChan
	}()

	// 获取故障中心的所有告警事件
	data, err := c.ctx.Redis.Alert().GetAllEvents(models.BuildAlertEventCacheKey(faultCenter.TenantId, faultCenter.ID))
	if err != nil {
		logc.Errorf(c.ctx.Ctx, "从 Redis 中获取事件信息错误, faultCenterKey: %s, err: %s", models.BuildAlertEventCacheKey(faultCenter.TenantId, faultCenter.ID), err.Error())
		return
	}

	// 同步状态
	c.syncStates(data)

	// 事件过滤
	filterEvents := c.filterAlertEvents(faultCenter, data)
	// 事件分组
	alertGroups := AlertGroups{
		Rules: make(map[string]RulesGroup),
	}
	c.alarmGrouping(faultCenter, &alertGroups, filterEvents)
	// 发送事件
	c.sendAlerts(faultCenter, &alertGroups)
	// 处理告警升级
	err = alarmUpgrade(c.ctx, faultCenter, data)
	if err != nil {
		logc.Error(c.ctx.Ctx, fmt.Sprintf("process alarm upgeade fail, err: %s", err.Error()))
	}
}

func (c *Consume) syncStates(alerts map[string]*models.AlertCurEvent) {
	for _, event := range alerts {
		if event == nil || event.IsRecovered {
			continue
		}

		currentStatus := event.GetEventStatus()
		event.Status = currentStatus
		muted := mute.IsSilence(mute.MuteParams{
			TenantId: event.TenantId,
			Labels:   event.Labels,
		})

		if currentStatus == models.StateMuting {
			if muted {
				continue
			}

			if err := event.TransitionStatus(models.StateAlerting); err != nil {
				logc.Errorf(c.ctx.Ctx, "解除静默状态失败, fingerprint=%s, err=%v", event.Fingerprint, err)
				continue
			}
			c.ctx.Redis.Alert().PushAlertEvent(event)
			continue
		}

		if !muted {
			continue
		}

		switch currentStatus {
		case models.StateAlerting:
			if err := event.TransitionStatus(models.StateMuting); err != nil {
				logc.Errorf(c.ctx.Ctx, "设置静默状态失败, fingerprint=%s, status=%s, err=%v", event.Fingerprint, currentStatus, err)
				continue
			}
			c.ctx.Redis.Alert().PushAlertEvent(event)
		}
	}
}

// filterAlertEvents 过滤告警事件
func (c *Consume) filterAlertEvents(faultCenter models.FaultCenter, alerts map[string]*models.AlertCurEvent) []*models.AlertCurEvent {
	var newEvents []*models.AlertCurEvent

	for _, event := range alerts {
		if event.Fingerprint == "" {
			continue
		}

		// 过滤掉 非告警中, 非恢复 状态的事件
		if event.Status != models.StateAlerting && event.Status != models.StateRecovered {
			continue
		}

		// 记录恢复状态的事件
		if event.IsRecovered {
			c.removeAlertFromCache(event)
			if err := process.RecordAlertHisEvent(c.ctx, *event); err != nil {
				logc.Error(c.ctx.Ctx, fmt.Sprintf("Failed to record alert history: %v", err))
			}
		}

		if valid := c.validateEvent(event, faultCenter); valid {
			newEvents = append(newEvents, event)
		}
	}

	return newEvents
}

// validateEvent 事件验证
func (c *Consume) validateEvent(event *models.AlertCurEvent, faultCenter models.FaultCenter) bool {
	return event.IsRecovered || event.LastSendTime == 0 ||
		event.LastEvalTime >= event.LastSendTime+int64(faultCenter.GetRepeatNoticeInterval(event.Severity))*60
}

// alarmGrouping 告警分组
// 会进行两次分组
// 第一次是状态+规则，避免不同状态及不同规则的事件分到一级组。
// 第二次时告警路由，与告警路由中 KV 匹配的事件分到二级组。
func (c *Consume) alarmGrouping(faultCenter models.FaultCenter, alertGroups *AlertGroups, alerts []*models.AlertCurEvent) {
	if len(alerts) == 0 {
		return
	}

	for _, alert := range alerts {
		// 状态+规则 = 状态 ID
		var stateId string
		if alert.IsRecovered {
			stateId = RecoverStatePrefix + alert.RuleId
		} else {
			stateId = FiringStatePrefix + alert.RuleId
		}

		alertGroups.AddAlert(stateId, alert, faultCenter)
	}
}

// sendAlerts 发送告警
func (c *Consume) sendAlerts(faultCenter models.FaultCenter, aggEvents *AlertGroups) {
	c.RLock()
	defer c.RUnlock()

	for _, rule := range aggEvents.Rules {
		for _, groups := range rule.Groups {
			c.processAlertGroup(faultCenter, groups.NoticeID, groups.Events)
		}
	}
}

// processAlertGroup 处理告警组
func (c *Consume) processAlertGroup(faultCenter models.FaultCenter, noticeId string, alerts []*models.AlertCurEvent) {
	g := new(errgroup.Group)
	g.Go(func() error { return handleAlert(c.ctx, "alarm", faultCenter, noticeId, alerts) })

	if err := g.Wait(); err != nil {
		logc.Errorf(c.ctx.Ctx, "Alert group processing failed: %v", err)
	}
}

// removeAlertFromCache 从缓存中删除告警
func (c *Consume) removeAlertFromCache(alert *models.AlertCurEvent) {
	c.ctx.Redis.Alert().RemoveAlertEvent(alert.TenantId, alert.FaultCenterId, alert.Fingerprint)
}

// RestartAllConsumers 重启消费进程
func (c *Consume) RestartAllConsumers() {
	list, err := ctx.DB.FaultCenter().List("", "")
	if err != nil {
		logc.Error(ctx.Ctx, fmt.Sprintf("获取故障中心列表错误, err: %s", err.Error()))
		return
	}
	for _, fc := range list {
		c.ctx.Redis.FaultCenter().PushFaultCenterInfo(fc)
		c.Submit(fc)
	}
}

// StopAllConsumers 停止所有消费者
func (c *Consume) StopAllConsumers() {
	c.taskMu.Lock()
	defer c.taskMu.Unlock()

	count := len(c.tasks)
	if count == 0 {
		return
	}

	logc.Infof(c.ctx.Ctx, "停止 %d 个故障中心消费者...", count)
	for faultCenterID := range c.tasks {
		c.stopLocked(faultCenterID)
	}

	logc.Infof(c.ctx.Ctx, "所有故障中心消费者已停止")
}

func evalCondition(metrics map[string]interface{}, noticeLabels []models.NoticeLabels) bool {
	for _, noticeLabel := range noticeLabels {
		value, exists := metrics[noticeLabel.Key]
		if !exists {
			return false
		}

		val, ok := value.(string)
		if !ok {
			continue
		}

		var matched bool
		switch noticeLabel.Operator {
		case "==", "=":
			matched = (val == noticeLabel.Value)
		case "!=":
			matched = (val != noticeLabel.Value)
		case "=~":
			matched, _ = tools.MatchRegexp(noticeLabel.Value, val)
		case "!~":
			if regexpMatched, err := tools.MatchRegexp(noticeLabel.Value, val); err == nil {
				matched = !regexpMatched
			}
		default:
			matched = false
		}

		if !matched {
			return false
		}
	}

	return true
}
