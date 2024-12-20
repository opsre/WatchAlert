package process

import (
	"fmt"
	"time"
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
		Labels:               rule.Labels,
		EvalInterval:         rule.EvalInterval,
		ForDuration:          rule.PrometheusConfig.ForDuration,
		NoticeId:             rule.NoticeId,
		NoticeGroup:          rule.NoticeGroup,
		IsRecovered:          false,
		RepeatNoticeInterval: rule.RepeatNoticeInterval,
		Severity:             rule.Severity,
		EffectiveTime:        rule.EffectiveTime,
		RecoverNotify:        rule.GetRecoverNotify(),
		AlarmAggregation:     rule.GetAlarmAggregation(),
	}
}

func SaveEventCache(ctx *ctx.Context, event models.AlertCurEvent) {
	ctx.Mux.Lock()
	defer ctx.Mux.Unlock()

	firingKey := event.GetFiringAlertCacheKey()
	pendingKey := event.GetPendingAlertCacheKey()

	// 判断改事件是否是Firing状态, 如果不是Firing状态 则标记Pending状态
	resFiring := ctx.Redis.Event().GetCache(firingKey)
	if resFiring.Fingerprint != "" {
		event.FirstTriggerTime = resFiring.FirstTriggerTime
		event.LastEvalTime = ctx.Redis.Event().GetLastEvalTime(firingKey)
		event.LastSendTime = resFiring.LastSendTime
	} else {
		event.FirstTriggerTime = ctx.Redis.Event().GetFirstTime(pendingKey)
		event.LastEvalTime = ctx.Redis.Event().GetLastEvalTime(pendingKey)
		event.LastSendTime = ctx.Redis.Event().GetLastSendTime(pendingKey)
		ctx.Redis.Event().SetCache("Pending", event, 0)
	}

	// 初次告警需要比对持续时间
	if resFiring.LastSendTime == 0 {
		if event.LastEvalTime-event.FirstTriggerTime < event.ForDuration {
			return
		}
	}

	ctx.Redis.Event().SetCache("Firing", event, 0)
	ctx.Redis.Event().DelCache(pendingKey)

}

/*
GcPendingCache
清理 Pending 数据的缓存.
场景: 第一次查询到有异常的指标会写入 Pending 缓存, 当该指标持续 Pending 到达持续时间后才会写入 Firing 缓存,
那么未到达持续时间并且该指标恢复正常, 那么就需要清理该指标的 Pending 数据.
*/
func GcPendingCache(ctx *ctx.Context, rule models.AlertRule, curKeys []string) {
	pendingKeys, err := ctx.Redis.Rule().GetAlertPendingCacheKeys(models.AlertRuleQuery{
		TenantId:         rule.TenantId,
		RuleId:           rule.RuleId,
		RuleGroupId:      rule.RuleGroupId,
		DatasourceIdList: rule.DatasourceIdList,
	})
	if err != nil {
		return
	}

	gcPendingKeys := tools.GetSliceDifference(pendingKeys, curKeys)
	for _, key := range gcPendingKeys {
		ctx.Redis.Event().DelCache(key)
	}
}

func GcRecoverWaitCache(ctx *ctx.Context, alarmRecoverStore storage.AlarmRecoverWaitStore, rule models.AlertRule, curKeys []string) {
	// 获取等待恢复告警的keys
	recoverWaitKeys := getRecoverWaitList(alarmRecoverStore, rule)
	// 删除正常告警的key
	firingKeys := tools.GetSliceSame(recoverWaitKeys, curKeys)
	deleteFiringKeys(ctx, alarmRecoverStore, firingKeys)
}

func getRecoverWaitList(recoverStore storage.AlarmRecoverWaitStore, rule models.AlertRule) []string {
	keyPrefix := fmt.Sprintf("%s", models.FiringAlertCachePrefix+rule.RuleId+"-"+rule.DatasourceIdList[0]+"-")
	return recoverStore.Search(keyPrefix)
}

func deleteFiringKeys(ctx *ctx.Context, recoverStore storage.AlarmRecoverWaitStore, keys []string) {
	ctx.Mux.Lock()
	defer ctx.Mux.Unlock()

	for _, key := range keys {
		recoverStore.Remove(key)
	}
}

// GetRedisFiringKeys 获取缓存所有Firing的Keys
func GetRedisFiringKeys(ctx *ctx.Context) []string {
	var keys []string
	cursor := uint64(0)
	pattern := "*" + ":" + models.FiringAlertCachePrefix + "*"
	// 每次获取的键数量
	count := int64(100)

	for {
		var curKeys []string
		var err error

		curKeys, cursor, err = ctx.Redis.Redis().Scan(cursor, pattern, count).Result()
		if err != nil {
			break
		}

		keys = append(keys, curKeys...)

		if cursor == 0 {
			break
		}
	}

	return keys
}

// GetNoticeGroupId 获取告警分组的通知ID
func GetNoticeGroupId(alert models.AlertCurEvent) string {
	if len(alert.NoticeGroup) != 0 {
		var noticeGroup []map[string]string
		for _, v := range alert.NoticeGroup {
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

	return alert.NoticeId
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
	}

	err := ctx.DB.Event().CreateHistoryEvent(hisData)
	if err != nil {
		return fmt.Errorf("RecordAlertHisEvent -> %s", err)
	}

	return nil
}
