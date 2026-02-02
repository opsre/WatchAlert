package process

import (
	"fmt"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"

	"github.com/zeromicro/go-zero/core/logc"
)

func BuildEvent(rule models.AlertRule, labels func() map[string]interface{}) models.AlertCurEvent {
	return models.AlertCurEvent{
		TenantId:             rule.TenantId,
		DatasourceType:       rule.DatasourceType,
		RuleGroupId:          rule.RuleGroupId,
		RuleId:               rule.RuleId,
		RuleName:             rule.RuleName,
		Labels:               labels(),
		EvalInterval:         rule.EvalInterval,
		IsRecovered:          false,
		RepeatNoticeInterval: rule.RepeatNoticeInterval,
		Severity:             rule.Severity,
		EffectiveTime:        rule.EffectiveTime,
		FaultCenterId:        rule.FaultCenterId,
	}
}

func PushEventToFaultCenter(ctx *ctx.Context, event *models.AlertCurEvent) {
	if event == nil {
		return
	}

	ctx.Mux.Lock()
	defer ctx.Mux.Unlock()
	if len(event.TenantId) <= 0 || len(event.Fingerprint) <= 0 {
		return
	}

	cache := ctx.Redis
	cacheEvent, _ := cache.Alert().GetEventFromCache(event.TenantId, event.FaultCenterId, event.Fingerprint)

	// 获取基础信息
	event.FirstTriggerTime = cacheEvent.GetFirstTime()
	event.LastEvalTime = cacheEvent.GetLastEvalTime()
	event.LastSendTime = cacheEvent.GetLastSendTime()
	event.ConfirmState = cacheEvent.GetLastConfirmState()
	event.EventId = cacheEvent.GetEventId()
	event.FaultCenter = cache.FaultCenter().GetFaultCenterInfo(models.BuildFaultCenterInfoCacheKey(event.TenantId, event.FaultCenterId))

	// 获取当前缓存中的状态
	currentStatus := cacheEvent.GetEventStatus()

	// 如果是新的告警事件，设置为 StatePreAlert
	if currentStatus == "" {
		event.Status = models.StatePreAlert
	} else {
		event.Status = currentStatus
	}

	// 根据不同情况处理状态转换
	switch event.Status {
	case models.StatePreAlert:
		if event.IsArriveForDuration() {
			// 如果达到持续时间，转为告警状态
			event.TransitionStatus(models.StateAlerting)
		}
	}

	// 最终再次校验 fingerprint 非空，避免 push 时使用空 key
	if event.Fingerprint == "" {
		logc.Errorf(ctx.Ctx, "PushEventToFaultCenter: fingerprint became empty before PushAlertEvent, tenant=%s, rule=%s(%s)", event.TenantId, event.RuleName, event.RuleId)
		return
	}

	// 更新缓存
	cache.Alert().PushAlertEvent(event)
}

// RecordAlertHisEvent 记录历史告警
func RecordAlertHisEvent(ctx *ctx.Context, alert models.AlertCurEvent) error {
	hisData := models.AlertHisEvent{
		TenantId:         alert.TenantId,
		RuleGroupId:      alert.RuleGroupId,
		EventId:          alert.EventId,
		DatasourceType:   alert.DatasourceType,
		DatasourceId:     alert.DatasourceId,
		Fingerprint:      alert.Fingerprint,
		RuleId:           alert.RuleId,
		RuleName:         alert.RuleName,
		Severity:         alert.Severity,
		Labels:           alert.Labels,
		EvalInterval:     alert.EvalInterval,
		Annotations:      alert.Annotations,
		FirstTriggerTime: alert.FirstTriggerTime,
		LastEvalTime:     alert.LastEvalTime,
		LastSendTime:     alert.LastSendTime,
		RecoverTime:      alert.RecoverTime,
		FaultCenterId:    alert.FaultCenterId,
		ConfirmState:     alert.ConfirmState,
		AlarmDuration:    alert.RecoverTime - alert.FirstTriggerTime,
		SearchQL:         alert.SearchQL,
	}

	err := ctx.DB.Event().CreateHistoryEvent(hisData)
	if err != nil {
		return fmt.Errorf("RecordAlertHisEvent, 恢复告警记录失败, err: %s", err)
	}

	return nil
}
