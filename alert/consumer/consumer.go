package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"golang.org/x/sync/errgroup"
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
		preStoreAlertEvents map[string][]*models.AlertCurEvent
	}
)

func NewConsumerWork(ctx *ctx.Context) ConsumeInterface {
	return &Consume{
		ctx:                 ctx,
		preStoreAlertEvents: make(map[string][]*models.AlertCurEvent),
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
	timer := time.NewTicker(time.Second * time.Duration(1))
	defer func() {
		timer.Stop()
		if r := recover(); r != nil {
			logc.Error(c.ctx.Ctx, fmt.Sprintf("Recovered from consumer watch goroutine panic: %s, FaultCenterName: %s, Id: %s", r, faultCenter.Name, faultCenter.ID))
		}
	}()

	for {
		select {
		case <-timer.C:
			c.processSilenceRule(faultCenter)
			// 获取故障中心的所有告警事件
			data, err := c.ctx.Redis.Redis().HGetAll(faultCenter.GetFaultCenterKey()).Result()
			if err != nil {
				logc.Error(c.ctx.Ctx, fmt.Sprintf("从 Redis 中获取事件信息错误, faultCenterKey: %s, err: %s", faultCenter.GetFaultCenterKey(), err.Error()))
				return
			}
			c.fireAlertEvent(faultCenter, c.filterAlertEvents(faultCenter, data))
			c.clear()
		case <-ctx.Done():
			return
		}
	}
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

		if !event.IsRecovered && c.isMutedEvent(event, faultCenter) {
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
	if event.State == "Pending" {
		return false
	}

	return event.IsRecovered || event.LastSendTime == 0 ||
		event.LastEvalTime >= event.LastSendTime+faultCenter.RepeatNoticeInterval*60
}

// addAlertToGroup 告警分组
func (c *Consume) addAlertToGroup(alert *models.AlertCurEvent, noticeGroupMap []map[string]string) {
	c.Lock()
	defer c.Unlock()

	groupId := alert.RuleId
	if len(noticeGroupMap) > 0 {
		for key, value := range alert.Metric {
			for _, noticeGroup := range noticeGroupMap {
				if noticeGroup["key"] == key && noticeGroup["value"] == value.(string) {
					groupId = tools.WithKVCalculateHash(key, value.(string)) + "_" + alert.RuleId
					break
				}
			}
		}
	}

	c.preStoreAlertEvents[groupId] = append(c.preStoreAlertEvents[groupId], alert)
}

// fireAlertEvent 触发告警事件
func (c *Consume) fireAlertEvent(faultCenter models.FaultCenter, alerts []*models.AlertCurEvent) {
	if len(alerts) == 0 {
		return
	}

	for _, alert := range alerts {
		c.addAlertToGroup(alert, faultCenter.NoticeGroup)
		if alert.IsRecovered {
			c.removeAlertFromCache(alert)
			if err := process.RecordAlertHisEvent(c.ctx, *alert); err != nil {
				logc.Error(c.ctx.Ctx, fmt.Sprintf("Failed to record alert history: %v", err))
			}
		}
	}

	c.sendAlerts(faultCenter, c.preStoreAlertEvents)
}

// sendAlerts 发送告警
func (c *Consume) sendAlerts(faultCenter models.FaultCenter, alertMapping map[string][]*models.AlertCurEvent) {
	c.RLock()
	defer c.RUnlock()

	for key, alerts := range alertMapping {
		ruleId := key
		if strings.Contains(key, "_") {
			ruleId = strings.Split(key, "_")[1]
		}

		rule := c.ctx.DB.Rule().GetRuleObject(ruleId)
		if rule.RuleId == "" || len(alerts) == 0 {
			continue
		}

		c.processAlertGroup(faultCenter, alerts)
	}
}

// processAlertGroup 处理告警组
func (c *Consume) processAlertGroup(faultCenter models.FaultCenter, alerts []*models.AlertCurEvent) {
	g := new(errgroup.Group)
	g.Go(func() error { return c.handleSubscribe(faultCenter, alerts) })
	g.Go(func() error { return c.handleAlert(faultCenter, alerts) })

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
			noticeId := process.GetNoticeGroupId(event, faultCenter)
			noticeData, err := c.getNoticeData(event.TenantId, noticeId)
			if err != nil {
				return fmt.Errorf("failed to get notice data: %v", err)
			}

			if err := processSubscribe(c.ctx, event, noticeData); err != nil {
				return fmt.Errorf("failed to process subscribe: %v", err)
			}

			return nil
		})
	}

	return g.Wait()
}

// handleAlert 处理告警逻辑
func (c *Consume) handleAlert(faultCenter models.FaultCenter, alerts []*models.AlertCurEvent) error {
	curTime := time.Now().Unix()
	switch faultCenter.GetAlarmAggregationType() {
	case "Rule":
		alerts = c.withRuleGroupByAlerts(curTime, alerts)
	default:
	}

	g := new(errgroup.Group)
	for _, alert := range alerts {
		g.Go(func() error {
			if alert == nil {
				return nil
			}
			noticeId := process.GetNoticeGroupId(alert, faultCenter)
			noticeData, err := c.getNoticeData(alert.TenantId, noticeId)
			if err != nil {
				logc.Error(c.ctx.Ctx, fmt.Sprintf("Failed to get notice data: %v", err))
				return err
			}

			if !alert.IsRecovered {
				alert.LastSendTime = curTime
				c.ctx.Redis.Event().PushEventToFaultCenter(alert)
			}

			phoneNumber := func() []string {
				if len(alert.DutyUserPhoneNumber) > 0 {
					return alert.DutyUserPhoneNumber
				}
				if len(noticeData.PhoneNumber) > 0 {
					return noticeData.PhoneNumber
				}
				return []string{}
			}()

			alert.DutyUser = process.GetDutyUser(c.ctx, noticeData)
			alert.DutyUserPhoneNumber = process.GetDutyUserPhoneNumber(c.ctx, noticeData)
			content := c.generateAlertContent(alert, noticeData)
			return sender.Sender(c.ctx, sender.SendParams{
				TenantId:    alert.TenantId,
				RuleName:    alert.RuleName,
				Severity:    alert.Severity,
				NoticeType:  noticeData.NoticeType,
				NoticeId:    noticeId,
				NoticeName:  noticeData.Name,
				IsRecovered: alert.IsRecovered,
				Hook:        noticeData.Hook,
				Email:       noticeData.Email,
				Content:     content,
				Event:       nil,
				PhoneNumber: phoneNumber,
			})
		})
	}

	return g.Wait()
}

// generateAlertContent 生成告警内容
func (c *Consume) generateAlertContent(alert *models.AlertCurEvent, noticeData models.AlertNotice) string {
	alert.FaultCenter = c.ctx.Redis.FaultCenter().GetFaultCenterInfo(models.BuildCacheInfoKey(alert.TenantId, alert.FaultCenterId))
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
		alert.Annotations += fmt.Sprintf("\n聚合 %d 条告警\n", len(alerts))
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

// clear 清楚本地缓存
func (c *Consume) clear() {
	c.Lock()
	defer c.Unlock()

	c.preStoreAlertEvents = make(map[string][]*models.AlertCurEvent)
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
