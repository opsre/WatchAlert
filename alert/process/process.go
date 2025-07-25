package process

import (
	"fmt"
	"time"
	"watchAlert/alert/mute"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
)

func BuildEvent(rule models.AlertRule, labels func() map[string]interface{}) models.AlertCurEvent {
	return models.AlertCurEvent{
		TenantId:             rule.TenantId,
		DatasourceType:       rule.DatasourceType,
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
	ctx.Mux.Lock()
	defer ctx.Mux.Unlock()
	if len(event.TenantId) <= 0 || len(event.Fingerprint) <= 0 {
		return
	}

	cache := ctx.Redis

	// 获取基础信息
	event.FirstTriggerTime = cache.Alert().GetFirstTime(event.TenantId, event.FaultCenterId, event.Fingerprint)
	event.LastEvalTime = cache.Alert().GetLastEvalTime()
	event.LastSendTime = cache.Alert().GetLastSendTime(event.TenantId, event.FaultCenterId, event.Fingerprint)
	event.UpgradeState = cache.Alert().GetLastUpgradeState(event.TenantId, event.FaultCenterId, event.Fingerprint)
	event.FaultCenter = cache.FaultCenter().GetFaultCenterInfo(models.BuildFaultCenterInfoCacheKey(event.TenantId, event.FaultCenterId))

	// 获取当前缓存中的状态
	currentStatus := cache.Alert().GetEventStatus(event.TenantId, event.FaultCenterId, event.Fingerprint)

	// 如果是新的告警事件，设置为 StatePreAlert
	if currentStatus == "" {
		event.Status = models.StatePreAlert
	} else {
		event.Status = currentStatus
	}

	// 检查是否处于静默状态
	isSilenced := IsSilencedEvent(event)

	// 根据不同情况处理状态转换
	switch event.Status {
	case models.StatePreAlert:
		// 如果需要静默
		if isSilenced {
			event.TransitionStatus(models.StateSilenced)
		} else if event.IsArriveForDuration() {
			// 如果达到持续时间，转为告警状态
			event.TransitionStatus(models.StateAlerting)
		}
	case models.StateAlerting:
		// 如果需要静默
		if isSilenced {
			event.TransitionStatus(models.StateSilenced)
		}
	case models.StateSilenced:
		// 如果不再静默，转换回预告警状态
		if !isSilenced {
			event.TransitionStatus(models.StatePreAlert)
		}
	}

	// 更新缓存
	cache.Alert().PushAlertEvent(event)
}

// IsSilencedEvent 静默检查
func IsSilencedEvent(event *models.AlertCurEvent) bool {
	return mute.IsSilence(mute.MuteParams{
		EffectiveTime: event.EffectiveTime,
		IsRecovered:   event.IsRecovered,
		TenantId:      event.TenantId,
		Labels:        event.Labels,
		FaultCenterId: event.FaultCenterId,
	})
}

func GetDutyUsers(ctx *ctx.Context, noticeData models.AlertNotice) []string {
	var us []string
	users, ok := ctx.DB.DutyCalendar().GetDutyUserInfo(*noticeData.GetDutyId(), time.Now().Format("2006-1-2"))
	if ok {
		switch noticeData.NoticeType {
		case "FeiShu":
			for _, user := range users {
				us = append(us, fmt.Sprintf("<at id=%s></at>", user.DutyUserId))
			}
			return us
		case "DingDing":
			for _, user := range users {
				us = append(us, fmt.Sprintf("@%s", user.DutyUserId))
			}
			return us
		case "Email", "WeChat", "CustomHook":
			for _, user := range users {
				us = append(us, fmt.Sprintf("@%s", user.UserName))
			}
			return us
		case "Slack":
			for _, user := range users {
				us = append(us, fmt.Sprintf("<@%s>", user.DutyUserId))
			}
			return us
		}
	}

	return []string{"暂无"}
}

// GetDutyUserPhoneNumber 获取当班人员手机号
func GetDutyUserPhoneNumber(ctx *ctx.Context, noticeData models.AlertNotice) []string {
	//user, ok := ctx.DB.DutyCalendar().GetDutyUserInfo(*noticeData.GetDutyId(), time.Now().Format("2006-1-2"))
	//if ok {
	//	switch noticeData.NoticeType {
	//	case "PhoneCall":
	//		if len(user.DutyUserId) > 1 {
	//			return []string{user.Phone}
	//		}
	//	}
	//}
	return []string{}
}

// RecordAlertHisEvent 记录历史告警
func RecordAlertHisEvent(ctx *ctx.Context, alert models.AlertCurEvent) error {
	hisData := models.AlertHisEvent{
		TenantId:         alert.TenantId,
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
		UpgradeState:     alert.UpgradeState,
		AlarmDuration:    alert.RecoverTime - alert.FirstTriggerTime,
	}

	err := ctx.DB.Event().CreateHistoryEvent(hisData)
	if err != nil {
		return fmt.Errorf("RecordAlertHisEvent -> %s", err)
	}

	return nil
}
