package process

import (
	"fmt"
	"time"
	"watchAlert/alert/mute"
	"watchAlert/alert/storage"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/tools"
)

func BuildEvent(rule models.AlertRule) models.AlertCurEvent {
	return models.AlertCurEvent{
		TenantId:             rule.TenantId,
		DatasourceType:       rule.DatasourceType,
		RuleId:               rule.RuleId,
		RuleName:             rule.RuleName,
		EvalInterval:         rule.EvalInterval,
		ForDuration:          rule.PrometheusConfig.ForDuration,
		IsRecovered:          false,
		RepeatNoticeInterval: rule.RepeatNoticeInterval,
		Severity:             rule.Severity,
		EffectiveTime:        rule.EffectiveTime,
		FaultCenterId:        rule.FaultCenterId,
		Status:               1,
	}
}

// isSilencedEvent 静默检查
func isSilencedEvent(event *models.AlertCurEvent) bool {
	return mute.IsSilence(mute.MuteParams{
		EffectiveTime: event.EffectiveTime,
		IsRecovered:   event.IsRecovered,
		TenantId:      event.TenantId,
		Metrics:       event.Metric,
		FaultCenterId: event.FaultCenterId,
	})
}

func PushEventToFaultCenter(ctx *ctx.Context, event *models.AlertCurEvent) {
	ctx.Mux.Lock()
	defer ctx.Mux.Unlock()
	if len(event.TenantId) <= 0 || len(event.Fingerprint) <= 0 {
		return
	}

	eventOpt := ctx.Redis.Event()
	event.FirstTriggerTime = eventOpt.GetFirstTimeForFaultCenter(event.TenantId, event.FaultCenterId, event.Fingerprint)
	event.LastEvalTime = eventOpt.GetLastEvalTimeForFaultCenter()
	event.LastSendTime = eventOpt.GetLastSendTimeForFaultCenter(event.TenantId, event.FaultCenterId, event.Fingerprint)

	if event.IsArriveForDuration() {
		event.Status = 1
	}
	if isSilencedEvent(event) {
		event.Status = 2
	}
	if event.IsRecovered {
		event.Status = 3
	}

	eventOpt.PushEventToFaultCenter(event)
}

func GcRecoverWaitCache(ctx *ctx.Context, alarmRecoverStore storage.AlarmRecoverWaitStore, rule models.AlertRule, curKeys []string) {
	// 获取等待恢复告警的keys
	recoverWaitKeys := getRecoverWaitList(alarmRecoverStore, rule)
	// 删除正常告警的key
	firingKeys := tools.GetSliceSame(recoverWaitKeys, curKeys)
	deleteFiringKeys(ctx, alarmRecoverStore, firingKeys)
}

func getRecoverWaitList(recoverStore storage.AlarmRecoverWaitStore, rule models.AlertRule) []string {
	var keys []string
	for _, dsId := range rule.DatasourceIdList {
		keyPrefix := fmt.Sprintf("%s", models.FiringAlertCachePrefix+rule.RuleId+"-"+dsId+"-")
		keys = append(keys, recoverStore.Search(keyPrefix)...)
	}
	return keys
}

func deleteFiringKeys(ctx *ctx.Context, recoverStore storage.AlarmRecoverWaitStore, keys []string) {
	ctx.Mux.Lock()
	defer ctx.Mux.Unlock()

	for _, key := range keys {
		recoverStore.Remove(key)
	}
}

// GetNoticeGroupId 获取告警分组的通知ID
func GetNoticeGroupId(alert *models.AlertCurEvent, faultCenter models.FaultCenter) string {
	if len(faultCenter.NoticeGroup) != 0 {
		var noticeGroup []map[string]string
		for _, v := range faultCenter.NoticeGroup {
			noticeGroup = append(noticeGroup, map[string]string{
				v["key"]:   v["value"],
				"noticeId": v["noticeId"],
			})
		}

		// 从Metric中获取Key/Value
		for metricKey, metricValue := range alert.Metric {
			// 如果配置分组的Key/Value 和 Metric中的Key/Value 一致，则使用分组的 noticeId，匹配不到则用默认的。
			for _, noticeInfo := range noticeGroup {
				value, ok := noticeInfo[metricKey]
				if ok && metricValue == value {
					noticeId := noticeInfo["noticeId"]
					return noticeId
				}
			}
		}
	}

	return faultCenter.NoticeId
}

func GetDutyUser(ctx *ctx.Context, noticeData models.AlertNotice) string {
	user, ok := ctx.DB.DutyCalendar().GetDutyUserInfo(noticeData.DutyId, time.Now().Format("2006-1-2"))
	if ok {
		switch noticeData.NoticeType {
		case "FeiShu":
			return fmt.Sprintf("<at id=%s></at>", user.DutyUserId)
		case "DingDing":
			return fmt.Sprintf("%s", user.DutyUserId)
		case "Email", "WeChat", "CustomHook":
			return fmt.Sprintf("@%s", user.UserName)
		}
	}

	return "暂无"
}

// GetDutyUserPhoneNumber 获取当班人员手机号
func GetDutyUserPhoneNumber(ctx *ctx.Context, noticeData models.AlertNotice) []string {
	user, ok := ctx.DB.DutyCalendar().GetDutyUserInfo(noticeData.DutyId, time.Now().Format("2006-1-2"))
	if ok {
		switch noticeData.NoticeType {
		case "PhoneCall":
			if len(user.DutyUserId) > 1 {
				return []string{user.Phone}
			}
		}
	}
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
		Metric:           alert.Metric,
		EvalInterval:     alert.EvalInterval,
		Annotations:      alert.Annotations,
		IsRecovered:      true,
		FirstTriggerTime: alert.FirstTriggerTime,
		LastEvalTime:     alert.LastEvalTime,
		LastSendTime:     alert.LastSendTime,
		RecoverTime:      alert.RecoverTime,
		FaultCenterId:    alert.FaultCenterId,
	}

	err := ctx.DB.Event().CreateHistoryEvent(hisData)
	if err != nil {
		return fmt.Errorf("RecordAlertHisEvent -> %s", err)
	}

	return nil
}
