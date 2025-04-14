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
		IsRecovered:           false,
		RepeatNoticeInterval:  rule.RepeatNoticeInterval,
		RecoverNotify:         rule.GetRecoverNotify(),
		ProbingEndpointConfig: rule.ProbingEndpointConfig,
	}
}

func SaveProbingEndpointEvent(event models.ProbingEvent) {
	firingKey := event.GetFiringAlertCacheKey()
	cache := ctx.DO().Cache.Event()
	resFiring, _ := cache.GetProbingEventCache(firingKey)
	event.FirstTriggerTime = cache.GetProbingEventFirstTime(firingKey)
	event.LastEvalTime = cache.GetProbingEventLastEvalTime(firingKey)
	event.LastSendTime = resFiring.LastSendTime
	cache.SetProbingEventCache(event, 0)
}

func SetProbingValueMap(key string, m map[string]any) error {
	for k, v := range m {
		ctx.DO().Cache.Cache().SetHashAny(key, k, v)
	}
	return nil
}

func GetProbingValueMap(key string) map[string]string {
	result, _ := ctx.DO().Cache.Cache().GetHashAll(key)
	return result
}
