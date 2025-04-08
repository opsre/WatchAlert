package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"golang.org/x/sync/errgroup"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"time"
	"watchAlert/alert/mute"
	"watchAlert/alert/process"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/sender"
	"watchAlert/pkg/templates"
	"watchAlert/pkg/tools"
)

type (
	ConsumeInterface interface {
		Submit(faultCenter models.FaultCenter)
		Stop(faultCenterId string)
		Watch(ctx context.Context, faultCenter models.FaultCenter)
		RestartAllConsumers()
	}

	Consume struct {
		ctx *ctx.Context
		sync.RWMutex
	}

	EventsGroup struct {
		ID     string // 事件组 ID
		Events []*models.AlertCurEvent
	}

	RulesGroup struct {
		RuleID string // 规则组 ID
		Groups []EventsGroup
	}

	AlertGroups struct {
		Rules []RulesGroup // 告警事件列表, 根据规则划分组
		lock  sync.Mutex
	}
)

func (ag *AlertGroups) AddAlert(stateId string, alert *models.AlertCurEvent, faultCenter models.FaultCenter) {
	// 获取通知对象 ID 列表 用于事件分组
	noticeObjIds := ag.getNoticeId(alert, faultCenter)

	ag.lock.Lock()
	defer ag.lock.Unlock()

	for _, noticeObjId := range noticeObjIds {
		// 查找 Rule 位置
		rulePos := ag.getRuleNodePos(stateId)

		// Rule 存在时的处理，找到对应的规则组
		if rulePos < len(ag.Rules) && ag.Rules[rulePos].RuleID == stateId {
			rule := &ag.Rules[rulePos]

			// 查找 Group 位置
			groupPos := ag.getGroupNodePos(rule, noticeObjId)

			if groupPos < len(rule.Groups) && (rule.Groups)[groupPos].ID == noticeObjId {
				// 追加事件
				(rule.Groups)[groupPos].Events = append((rule.Groups)[groupPos].Events, alert)
			} else {
				// 插入新数据
				rule.Groups = append(rule.Groups, EventsGroup{
					ID:     noticeObjId,
					Events: []*models.AlertCurEvent{alert},
				})
			}
			return
		}

		// 插入新Rule
		ag.Rules = append(ag.Rules, RulesGroup{
			RuleID: stateId,
			Groups: []EventsGroup{
				{
					ID:     noticeObjId,
					Events: []*models.AlertCurEvent{alert},
				},
			},
		})
	}
}

// getNoticeId 从告警路由中获取该事件匹配的通知对象
func (ag *AlertGroups) getNoticeId(alert *models.AlertCurEvent, faultCenter models.FaultCenter) []string {
	if len(faultCenter.NoticeRoutes) > 0 {
		metrics := alert.Metric

		for _, route := range faultCenter.NoticeRoutes {
			if metrics[route.Key] == route.Value {
				return route.NoticeIds
			}
		}
	}

	return faultCenter.NoticeIds
}

// getRuleNodePos 获取 Rule 点位
func (ag *AlertGroups) getRuleNodePos(ruleId string) int {
	// Rules 切片排序
	sort.Slice(ag.Rules, func(i, j int) bool {
		return ag.Rules[i].RuleID < ag.Rules[j].RuleID
	})

	// 查找Rule位置
	return sort.Search(len(ag.Rules), func(i int) bool {
		return ag.Rules[i].RuleID >= ruleId
	})
}

func (ag *AlertGroups) getGroupNodePos(rule *RulesGroup, groupId string) int {
	// Groups 切片排序
	sort.Slice(rule.Groups, func(i, j int) bool {
		return rule.Groups[i].ID < rule.Groups[j].ID
	})

	// 查找Group位置
	return sort.Search(len(rule.Groups), func(i int) bool {
		return (rule.Groups)[i].ID >= groupId
	})
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
	c.ctx.ConsumerContextMap[faultCenter.ID] = cancel
	go c.Watch(withCtx, faultCenter)
}

func (c *Consume) Stop(faultCenterId string) {
	c.ctx.Mux.Lock()
	defer c.ctx.Mux.Unlock()

	if cancel, exists := c.ctx.ConsumerContextMap[faultCenterId]; exists {
		cancel()
		delete(c.ctx.ConsumerContextMap, faultCenterId)
	}
}

// Watch 启动 Consumer Watch 进程
func (c *Consume) Watch(ctx context.Context, faultCenter models.FaultCenter) {
	taskChan := make(chan struct{}, 1)
	timer := time.NewTicker(time.Second * time.Duration(1))
	defer func() {
		timer.Stop()
		if r := recover(); r != nil {
			// 获取调用栈信息
			stack := debug.Stack()
			logc.Error(c.ctx.Ctx, fmt.Sprintf("Recovered from consumer watch goroutine panic: %s, FaultCenterName: %s, Id: %s\n%s", r, faultCenter.Name, faultCenter.ID, stack))
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
	data, err := c.ctx.Redis.Redis().HGetAll(faultCenter.GetFaultCenterKey()).Result()
	if err != nil {
		logc.Error(c.ctx.Ctx, fmt.Sprintf("从 Redis 中获取事件信息错误, faultCenterKey: %s, err: %s", faultCenter.GetFaultCenterKey(), err.Error()))
		return
	}

	// 事件过滤
	filterEvents := c.filterAlertEvents(faultCenter, data)
	// 事件分组
	var alertGroups AlertGroups
	c.alarmGrouping(faultCenter, &alertGroups, filterEvents)
	// 发送事件
	c.sendAlerts(faultCenter, &alertGroups)
}

// filterAlertEvents 过滤告警事件
func (c *Consume) filterAlertEvents(faultCenter models.FaultCenter, alerts map[string]string) []*models.AlertCurEvent {
	var newEvents []*models.AlertCurEvent

	for _, alert := range alerts {
		var event *models.AlertCurEvent
		if err := json.Unmarshal([]byte(alert), &event); err != nil {
			logc.Error(c.ctx.Ctx, fmt.Sprintf("Failed to unmarshal alert: %v", err))
			continue
		}

		// 过滤掉 预告警, 待恢复 状态的事件
		if event.Status == 0 || event.Status == 3 {
			continue
		}

		if c.isMutedEvent(event, faultCenter) {
			continue
		}

		if valid := c.validateEvent(event, faultCenter); valid {
			newEvents = append(newEvents, event)
		}
	}

	return newEvents
}

// isMutedEvent 静默检查
func (c *Consume) isMutedEvent(event *models.AlertCurEvent, faultCenter models.FaultCenter) bool {
	return mute.IsMuted(mute.MuteParams{
		EffectiveTime: event.EffectiveTime,
		IsRecovered:   event.IsRecovered,
		TenantId:      event.TenantId,
		Metrics:       event.Metric,
		FaultCenterId: event.FaultCenterId,
		RecoverNotify: faultCenter.RecoverNotify,
	})
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
		switch alert.IsRecovered {
		case true:
			stateId = "Recover_" + alert.RuleId
		case false:
			stateId = "Firing_" + alert.RuleId
		default:
			stateId = "Unknown_" + alert.RuleId
		}

		alertGroups.AddAlert(stateId, alert, faultCenter)
		if alert.IsRecovered {
			c.removeAlertFromCache(alert)
			if err := process.RecordAlertHisEvent(c.ctx, *alert); err != nil {
				logc.Error(c.ctx.Ctx, fmt.Sprintf("Failed to record alert history: %v", err))
			}
		}
	}
}

// alarmAggregation 告警聚合
func (c *Consume) alarmAggregation(faultCenter models.FaultCenter, alertGroups map[string][]*models.AlertCurEvent) map[string][]*models.AlertCurEvent {
	curTime := time.Now().Unix()
	newAlertGroups := alertGroups
	switch faultCenter.GetAlarmAggregationType() {
	case "Rule":
		for severity, events := range alertGroups {
			newAlertGroups[severity] = c.withRuleGroupByAlerts(curTime, events)
		}
	default:
		return alertGroups
	}

	return newAlertGroups
}

// sendAlerts 发送告警
func (c *Consume) sendAlerts(faultCenter models.FaultCenter, aggEvents *AlertGroups) {
	c.RLock()
	defer c.RUnlock()

	for _, rule := range aggEvents.Rules {
		for _, groups := range rule.Groups {
			c.processAlertGroup(faultCenter, groups.ID, groups.Events)
		}
	}
}

// processAlertGroup 处理告警组
func (c *Consume) processAlertGroup(faultCenter models.FaultCenter, noticeId string, alerts []*models.AlertCurEvent) {
	g := new(errgroup.Group)
	g.Go(func() error { return c.handleSubscribe(faultCenter, alerts) })
	g.Go(func() error { return c.handleAlert(faultCenter, noticeId, alerts) })

	if err := g.Wait(); err != nil {
		logc.Error(c.ctx.Ctx, fmt.Sprintf("Alert group processing failed: %v", err))
	}
}

// handleSubscribe 处理订阅逻辑
func (c *Consume) handleSubscribe(faultCenter models.FaultCenter, alerts []*models.AlertCurEvent) error {
	g := new(errgroup.Group)
	for _, event := range alerts {
		event := event
		g.Go(func() error {
			event.FaultCenter = faultCenter
			if err := processSubscribe(c.ctx, event); err != nil {
				return fmt.Errorf("failed to process subscribe: %v", err)
			}

			return nil
		})
	}

	return g.Wait()
}

// handleAlert 处理告警逻辑
func (c *Consume) handleAlert(faultCenter models.FaultCenter, noticeId string, alerts []*models.AlertCurEvent) error {
	curTime := time.Now().Unix()
	g := new(errgroup.Group)

	// 获取通知对象详细信息
	noticeData, err := c.getNoticeData(faultCenter.TenantId, noticeId)
	if err != nil {
		logc.Error(c.ctx.Ctx, fmt.Sprintf("Failed to get notice data: %v", err))
		return err
	}

	// 按告警等级分组
	severityGroups := make(map[string][]*models.AlertCurEvent)
	for _, alert := range alerts {
		severityGroups[alert.Severity] = append(severityGroups[alert.Severity], alert)
	}

	// 告警聚合
	aggregationEvents := c.alarmAggregation(faultCenter, severityGroups)
	for severity, events := range aggregationEvents {
		g.Go(func() error {
			if events == nil {
				return nil
			}

			// 获取当前事件等级对应的 Hook 和 Sign
			Hook, Sign := c.getNoticeHookUrlAndSign(noticeData, severity)

			for _, event := range events {
				if !event.IsRecovered {
					event.LastSendTime = curTime
					c.ctx.Redis.Event().PushEventToFaultCenter(event)
				}

				phoneNumber := func() []string {
					if len(event.DutyUserPhoneNumber) > 0 {
						return event.DutyUserPhoneNumber
					}
					if len(noticeData.PhoneNumber) > 0 {
						return noticeData.PhoneNumber
					}
					return []string{}
				}()

				event.DutyUser = process.GetDutyUser(c.ctx, noticeData)
				event.DutyUserPhoneNumber = process.GetDutyUserPhoneNumber(c.ctx, noticeData)
				event.FaultCenter = faultCenter
				content := c.generateAlertContent(event, noticeData)
				return sender.Sender(c.ctx, sender.SendParams{
					TenantId:    event.TenantId,
					RuleName:    event.RuleName,
					Severity:    event.Severity,
					NoticeType:  noticeData.NoticeType,
					NoticeId:    noticeId,
					NoticeName:  noticeData.Name,
					IsRecovered: event.IsRecovered,
					Hook:        Hook,
					Email:       c.getNoticeEmail(noticeData, severity),
					Content:     content,
					Event:       nil,
					PhoneNumber: phoneNumber,
					Sign:        Sign,
				})
			}
			return nil
		})
	}

	return g.Wait()
}

// getNoticeHookUrlAndSign 获取事件等级对应的 Hook 和 Sign
func (c *Consume) getNoticeHookUrlAndSign(notice models.AlertNotice, severity string) (string, string) {
	if notice.Routes != nil {
		for _, hook := range notice.Routes {
			if hook.Severity == severity {
				return hook.Hook, hook.Sign
			}
		}
	}
	return notice.DefaultHook, notice.DefaultSign
}

// getNoticeEmail 获取事件等级对应的 Email
func (c *Consume) getNoticeEmail(notice models.AlertNotice, severity string) models.Email {
	if notice.Routes != nil {
		for _, route := range notice.Routes {
			if route.Severity == severity {
				return models.Email{
					Subject: notice.Email.Subject,
					To:      route.To,
					CC:      route.CC,
				}
			}
		}
	}
	return notice.Email
}

// generateAlertContent 生成告警内容
func (c *Consume) generateAlertContent(alert *models.AlertCurEvent, noticeData models.AlertNotice) string {
	if noticeData.NoticeType == "CustomHook" {
		return tools.JsonMarshal(alert)
	}
	return templates.NewTemplate(c.ctx, *alert, noticeData).CardContentMsg
}

// withRuleGroupByAlerts 聚合告警
func (c *Consume) withRuleGroupByAlerts(timeInt int64, alerts []*models.AlertCurEvent) []*models.AlertCurEvent {
	if len(alerts) <= 1 {
		return alerts
	}

	var aggregatedAlert *models.AlertCurEvent
	for i := range alerts {
		alert := alerts[i]
		if !strings.Contains(alert.Annotations, "聚合") {
			alert.Annotations += fmt.Sprintf("\n聚合 %d 条告警\n", len(alerts))
		}
		aggregatedAlert = alert

		if !alert.IsRecovered {
			alert.LastSendTime = timeInt
			c.ctx.Redis.Event().PushEventToFaultCenter(alert)
		}
	}

	return []*models.AlertCurEvent{aggregatedAlert}
}

// removeAlertFromCache 从缓存中删除告警
func (c *Consume) removeAlertFromCache(alert *models.AlertCurEvent) {
	c.ctx.Redis.Event().RemoveEventFromFaultCenter(alert.TenantId, alert.FaultCenterId, alert.Fingerprint)
}

// getNoticeData 获取 Notice 数据
func (c *Consume) getNoticeData(tenantId, noticeId string) (models.AlertNotice, error) {
	return c.ctx.DB.Notice().Get(models.NoticeQuery{
		TenantId: tenantId,
		Uuid:     noticeId,
	})
}

// RestartAllConsumers 重启消费进程
func (c *Consume) RestartAllConsumers() {
	list, err := ctx.DB.FaultCenter().List(models.FaultCenterQuery{})
	if err != nil {
		logc.Error(ctx.Ctx, fmt.Sprintf("获取故障中心列表错误, err: %s", err.Error()))
		return
	}
	for _, fc := range list {
		c.Submit(fc)
	}
}

func (c *Consume) processSilenceRule(faultCenter models.FaultCenter) {
	currentTime := time.Now().Unix()
	silenceCtx := c.ctx.Redis.Silence()
	// 获取静默列表中所有的id
	silenceIds, err := silenceCtx.GetMutesForFaultCenter(faultCenter.TenantId, faultCenter.ID)
	if err != nil {
		logc.Errorf(ctx.Ctx, err.Error())
		return
	}

	// 根据ID获取到详细的静默规则
	for _, silenceId := range silenceIds {
		muteRule, err := silenceCtx.WithIdGetMuteFromCache(faultCenter.TenantId, faultCenter.ID, silenceId)
		if err != nil {
			logc.Errorf(ctx.Ctx, err.Error())
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

		silenceCtx.PushMuteToFaultCenter(*muteRule)
	}
}
