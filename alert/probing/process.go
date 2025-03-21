package probing

import (
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
)

func (t *ProductProbing) buildEvent(rule models.ProbingRule) models.ProbingEvent {
	return models.ProbingEvent{
		TenantId:              rule.TenantId,
		RuleId:                rule.RuleId,
		RuleName:              rule.RuleName,
		RuleType:              rule.RuleType,
		NoticeId:              rule.NoticeId,
		Severity:              rule.Severity,
		IsRecovered:           false,
		RepeatNoticeInterval:  rule.RepeatNoticeInterval,
		RecoverNotify:         rule.GetRecoverNotify(),
		ProbingEndpointConfig: rule.ProbingEndpointConfig,
	}
}

func SaveProbingEndpointEvent(event models.ProbingEvent) {
	firingKey := event.GetFiringAlertCacheKey()
	cache := ctx.DO().Redis.Event()
	resFiring, _ := cache.GetProbingEventCache(firingKey)
	event.FirstTriggerTime = cache.GetProbingEventFirstTime(firingKey)
	event.LastEvalTime = cache.GetProbingEventLastEvalTime(firingKey)
	event.LastSendTime = resFiring.LastSendTime
	cache.SetProbingEventCache(event, 0)
}

func SetProbingValueMap(key string, m map[string]any) error {
	for k, v := range m {
		ctx.DO().Redis.Redis().HSet(key, k, v)
	}
	return nil
}

func GetProbingValueMap(key string) map[string]string {
	result := ctx.DO().Redis.Redis().HGetAll(key).Val()
	return result
}
