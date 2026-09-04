package process

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"

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
	pushEventToFaultCenter(ctx, event, "")
}

// PushEventToFaultCenterWithStatus persists an explicit state transition instead of
// restoring the event status from cache. Metric evaluation uses it for suppression
// and for promoting a previously suppressed series back to alerting.
func PushEventToFaultCenterWithStatus(ctx *ctx.Context, event *models.AlertCurEvent, targetStatus models.AlertStatus) {
	pushEventToFaultCenter(ctx, event, targetStatus)
}

func pushEventToFaultCenter(ctx *ctx.Context, event *models.AlertCurEvent, targetStatus models.AlertStatus) {
	if event == nil {
		return
	}

	if NotInTheEffectiveTime(event.EffectiveTime) {
		return
	}

	ctx.Mux.Lock()
	defer ctx.Mux.Unlock()
	if len(event.TenantId) <= 0 || len(event.Fingerprint) <= 0 {
		return
	}

	cache := ctx.Redis
	cacheEvent, _ := cache.Alert().GetEventFromCache(event.TenantId, event.FaultCenterId, event.Fingerprint)
	// 非零值由调用方显式设置（例如高优先级抑制低优先级时），不能被缓存覆盖。
	explicitLastSendTime := event.LastSendTime

	// 获取基础信息
	event.FirstTriggerTime = cacheEvent.GetFirstTime()
	event.LastEvalTime = cacheEvent.GetLastEvalTime()
	event.LastSendTime = cacheEvent.GetLastSendTime()
	if explicitLastSendTime != 0 {
		event.LastSendTime = explicitLastSendTime
	}
	event.ConfirmState = cacheEvent.GetLastConfirmState()
	event.EventId = cacheEvent.GetEventId()
	event.FaultCenter = cache.FaultCenter().GetFaultCenterInfo(models.BuildFaultCenterInfoCacheKey(event.TenantId, event.FaultCenterId))

	// 以缓存状态为转换起点，确保显式转换遵循状态机并最终持久化。
	event.Status = cacheEvent.GetEventStatus()
	if targetStatus != "" {
		if err := event.TransitionStatus(targetStatus); err != nil {
			logc.Errorf(ctx.Ctx, "PushEventToFaultCenter: transition status failed, fingerprint=%s, from=%s, to=%s, err=%v", event.Fingerprint, event.Status, targetStatus, err)
			return
		}
	} else if event.Status == models.StatePreAlert && event.IsArriveForDuration() {
		// 如果达到持续时间，转为告警状态。
		if err := event.TransitionStatus(models.StateAlerting); err != nil {
			logc.Errorf(ctx.Ctx, "PushEventToFaultCenter: transition alerting failed, fingerprint=%s, err=%v", event.Fingerprint, err)
			return
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

// NotInTheEffectiveTime 判断是否不在生效时间内
func NotInTheEffectiveTime(et models.EffectiveTime) bool {
	// 如果没有配置有效星期，则认为始终有效
	if len(et.Week) == 0 {
		return false
	}

	// 当前日期
	currentTime := time.Now()
	currentWeekday := tools.TimeTransformToWeek(currentTime)

	// 检查当前星期是否在有效范围内
	var isInValidWeekday bool
	for _, weekday := range et.Week {
		if currentWeekday == weekday {
			isInValidWeekday = true
			break
		}
	}

	// 如果当前星期不在有效范围内，直接返回 true
	if !isInValidWeekday {
		return true
	}

	// 如果开始时间和结束时间都为0，表示全天有效
	if et.StartTime == 0 && et.EndTime == 0 {
		return false
	}

	// 检查当前时间是否在指定的时间段内
	currentTimeSeconds := tools.TimeTransformToSeconds(currentTime)
	isInValidTimeRange := currentTimeSeconds >= et.StartTime && currentTimeSeconds <= et.EndTime

	// 返回是否 不在有效时间范围内的结果
	return !isInValidTimeRange
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

// ProcessRuleExpr 处理规则表达式
func ProcessRuleExpr(ruleExpr string) (operator string, value float64, err error) {
	var supportedOperators = []string{">=", "<=", "==", "!=", ">", "<", "="}

	// 去除表达式两端的空白字符
	trimmedExpr := strings.TrimSpace(ruleExpr)

	// 遍历操作符列表。
	for _, op := range supportedOperators {
		if strings.HasPrefix(trimmedExpr, op) {
			// 提取数值
			valueStr := strings.TrimPrefix(trimmedExpr, op)
			value, err = strconv.ParseFloat(strings.TrimSpace(valueStr), 64)
			if err != nil {
				return "", 0, fmt.Errorf("无法解析数值 '%s': %w", valueStr, err)
			}

			return op, value, nil
		}
	}

	return "", 0, fmt.Errorf("无效的表达式，未找到有效的操作符: %s", ruleExpr)
}
