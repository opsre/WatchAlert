package process

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"watchAlert/alert/queue"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
)

// GetSliceDifference 获取差异key. 当slice1中存在, slice2不存在则标记为可恢复告警
func GetSliceDifference(slice1 []string, slice2 []string) []string {
	difference := []string{}

	// 遍历缓存
	for _, item1 := range slice1 {
		found := false
		// 遍历当前key
		for _, item2 := range slice2 {
			if item1 == item2 {
				found = true
				break
			}
		}
		// 添加到差异切片中
		if !found {
			difference = append(difference, item1)
		}
	}

	return difference
}

// GetSliceSame 获取相同key, 当slice1中存在, slice2也存在则标记为正在告警中撤销告警恢复
func GetSliceSame(slice1 []string, slice2 []string) []string {
	same := []string{}
	for _, item1 := range slice1 {
		for _, item2 := range slice2 {
			if item1 == item2 {
				same = append(same, item1)
			}
		}
	}
	return same
}

func ParserDefaultEvent(rule models.AlertRule) models.AlertCurEvent {

	event := models.AlertCurEvent{
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
		DutyUser:             "暂无", // 默认暂无值班人员, 渲染模版时会实际判断 Notice 是否存在值班人员
		Severity:             rule.Severity,
		EffectiveTime:        rule.EffectiveTime,
		RecoverNotify:        rule.RecoverNotify,
	}

	return event

}

func SaveEventCache(ctx *ctx.Context, event models.AlertCurEvent) {
	ctx.Lock()
	defer ctx.Unlock()

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

// ParserDuration 获取时间区间的开始时间
func ParserDuration(curTime time.Time, logScope int, timeType string) time.Time {

	duration, err := time.ParseDuration(strconv.Itoa(logScope) + timeType)
	if err != nil {
		global.Logger.Sugar().Error("解析相对时间失败 ->", err.Error())
		return time.Time{}
	}
	startsAt := curTime.Add(-duration)

	return startsAt

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

	gcPendingKeys := GetSliceDifference(pendingKeys, curKeys)
	for _, key := range gcPendingKeys {
		ctx.Redis.Event().DelCache(key)
	}
}

func GcRecoverWaitCache(rule models.AlertRule, curKeys []string) {
	// 获取等待恢复告警的keys
	recoverWaitKeys := getRecoverWaitList(queue.RecoverWaitMap, rule)
	// 删除正常告警的key
	firingKeys := GetSliceSame(recoverWaitKeys, curKeys)
	for _, key := range firingKeys {
		delete(queue.RecoverWaitMap, key)
	}
}

func getRecoverWaitList(m map[string]int64, rule models.AlertRule) []string {
	var l []string
	for k, _ := range m {
		// 只获取当前规则组的告警。
		keyPrefix := fmt.Sprintf("%s", models.FiringAlertCachePrefix+rule.RuleId+"-"+rule.DatasourceIdList[0]+"-")
		if strings.HasPrefix(k, keyPrefix) {
			l = append(l, k)
		}
	}
	return l
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
	user := ctx.DB.DutyCalendar().GetDutyUserInfo(noticeData.DutyId, time.Now().Format("2006-1-2"))
	switch noticeData.NoticeType {
	case "FeiShu":
		// 判断是否有安排值班人员
		if len(user.DutyUserId) > 1 {
			return fmt.Sprintf("<at id=%s></at>", user.DutyUserId)
		}
	case "DingDing":
		if len(user.DutyUserId) > 1 {
			return fmt.Sprintf("%s", user.DutyUserId)
		}
	}

	return ""
}

// GetDutyUserPhoneNumber 获取当班人员手机号
func GetDutyUserPhoneNumber(ctx *ctx.Context, noticeData models.AlertNotice) []string {
	user := ctx.DB.DutyCalendar().GetDutyUserInfo(noticeData.DutyId, time.Now().Format("2006-1-2"))
	switch noticeData.NoticeType {
	case "PhoneCall":
		if len(user.DutyUserId) > 1 {
			return []string{user.Phone}
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
	}

	err := ctx.DB.Event().CreateHistoryEvent(hisData)
	if err != nil {
		return fmt.Errorf("RecordAlertHisEvent -> %s", err)
	}

	return nil
}
