package consumer

import (
	"context"
	"fmt"
	"regexp"
	"runtime/debug"
	"sync"
	"time"
	"watchAlert/alert/process"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"

	"github.com/zeromicro/go-zero/core/logc"
	"golang.org/x/sync/errgroup"
)

const (
	// 任务通道缓冲区大小
	TaskChannelBufferSize = 1

	// 默认处理时间间隔
	DefaultProcessTime = 1

	// 状态前缀
	RecoverStatePrefix = "Recover_"
	FiringStatePrefix  = "Firing_"
)

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
			val, ok := labels[route.Key].(string)
			if !ok {
				continue
			}

			if regexp.MustCompile(route.Value).MatchString(val) {
				return route.NoticeIds
			}
		}
	}

	return faultCenter.NoticeIds
}

func NewConsumerWork(ctx *ctx.Context) ConsumeInterface {
	return &Consume{
		ctx: ctx,
	}
}

func (c *Consume) Submit(faultCenter models.FaultCenter) {
	c.ctx.Mux.Lock()
	defer c.ctx.Mux.Unlock()

	withCtx, cancel := context.WithCancel(context.Background())
	c.ctx.ContextMap[faultCenter.ID] = cancel
	go c.Watch(withCtx, faultCenter)
}

func (c *Consume) Stop(faultCenterId string) {
	c.ctx.Mux.Lock()
	defer c.ctx.Mux.Unlock()

	if cancel, exists := c.ctx.ContextMap[faultCenterId]; exists {
		cancel()
		delete(c.ctx.ContextMap, faultCenterId)
	}
}

func (c *Consume) Restart(faultCenter models.FaultCenter) {
	c.Stop(faultCenter.ID)
	c.Submit(faultCenter)
}

// Watch 启动 Consumer Watch 进程
func (c *Consume) Watch(ctx context.Context, faultCenter models.FaultCenter) {
	taskChan := make(chan struct{}, TaskChannelBufferSize)
	timer := time.NewTicker(time.Second * time.Duration(DefaultProcessTime))
	defer func() {
		timer.Stop()
		if r := recover(); r != nil {
			// 获取调用栈信息
			stack := debug.Stack()
			logc.Error(c.ctx.Ctx, fmt.Sprintf("Recovered from consumer watch goroutine panic: %s, FaultCenterName: %s, Id: %s\n%s", r, faultCenter.Name, faultCenter.ID, stack))
			c.Restart(faultCenter)
		}
	}()

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
	// 处理静默规则
	c.processSilenceRule(faultCenter)
	// 获取故障中心的所有告警事件
	data, err := c.ctx.Redis.Alert().GetAllEvents(models.BuildAlertEventCacheKey(faultCenter.TenantId, faultCenter.ID))
	if err != nil {
		logc.Errorf(c.ctx.Ctx, "从 Redis 中获取事件信息错误, faultCenterKey: %s, err: %s", models.BuildAlertEventCacheKey(faultCenter.TenantId, faultCenter.ID), err.Error())
		return
	}

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

// filterAlertEvents 过滤告警事件
func (c *Consume) filterAlertEvents(faultCenter models.FaultCenter, alerts map[string]*models.AlertCurEvent) []*models.AlertCurEvent {
	var newEvents []*models.AlertCurEvent

	for _, event := range alerts {
		if event.Fingerprint == "" {
			continue
		}

		// 过滤掉 非告警中 状态的事件
		if event.Status != models.StateAlerting {
			if event.IsRecovered {
				c.removeAlertFromCache(event)
				if err := process.RecordAlertHisEvent(c.ctx, *event); err != nil {
					logc.Error(c.ctx.Ctx, fmt.Sprintf("Failed to record alert history: %v", err))
				}
			}

			continue
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
		event.LastEvalTime >= event.LastSendTime+faultCenter.RepeatNoticeInterval*60
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
	g.Go(func() error { return c.handleSubscribe(alerts) })
	g.Go(func() error { return handleAlert(c.ctx, "alarm", faultCenter, noticeId, alerts) })

	if err := g.Wait(); err != nil {
		logc.Errorf(c.ctx.Ctx, "Alert group processing failed: %v", err)
	}
}

// handleSubscribe 处理订阅逻辑
func (c *Consume) handleSubscribe(alerts []*models.AlertCurEvent) error {
	g := new(errgroup.Group)
	for _, event := range alerts {
		event := event
		g.Go(func() error {
			if err := processSubscribe(c.ctx, event); err != nil {
				return fmt.Errorf("failed to process subscribe: %v", err)
			}

			return nil
		})
	}

	return g.Wait()
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

func (c *Consume) processSilenceRule(faultCenter models.FaultCenter) {
	currentTime := time.Now().Unix()
	silenceCtx := c.ctx.Redis.Silence()
	// 获取静默列表中所有的id
	silenceIds, err := silenceCtx.GetAlertMutes(faultCenter.TenantId, faultCenter.ID)
	if err != nil {
		logc.Error(ctx.Ctx, err.Error())
		return
	}

	// 根据ID获取到详细的静默规则
	for _, silenceId := range silenceIds {
		muteRule, err := silenceCtx.WithIdGetMuteFromCache(faultCenter.TenantId, faultCenter.ID, silenceId)
		if err != nil {
			logc.Error(ctx.Ctx, err.Error())
			return
		}

		// 如果当前状态为「未生效」，并且生效时间大于等于当前时间，则标记为「生效中」状态
		if muteRule.Status == 0 && currentTime >= muteRule.StartsAt {
			muteRule.Status = 1
			err := c.ctx.DB.Silence().Update(*muteRule)
			if err != nil {
				logc.Error(c.ctx.Ctx, fmt.Sprintf("Update silence rule failed, err: %s", err.Error()))
				return
			}
		}

		// 如果到达失效日期，则标记「已失效」状态
		if muteRule.EndsAt <= currentTime {
			muteRule.Status = 2
			err := c.ctx.DB.Silence().Update(*muteRule)
			if err != nil {
				logc.Error(c.ctx.Ctx, fmt.Sprintf("Update silence rule failed, err: %s", err.Error()))
				return
			}
		}

		silenceCtx.PushAlertMute(*muteRule)
	}
}

// StopAllConsumers 停止所有消费者
func (c *Consume) StopAllConsumers() {
	c.ctx.Mux.Lock()
	defer c.ctx.Mux.Unlock()

	count := len(c.ctx.ContextMap)
	if count == 0 {
		return
	}

	logc.Infof(c.ctx.Ctx, "停止 %d 个故障中心消费者...", count)

	// 取消所有消费任务
	for fcId, cancel := range c.ctx.ContextMap {
		cancel()
		delete(c.ctx.ContextMap, fcId)
	}

	logc.Infof(c.ctx.Ctx, "所有故障中心消费者已停止")
}
