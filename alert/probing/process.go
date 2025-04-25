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

func SaveProbingEndpointEvent(ctx *ctx.Context, event models.ProbingEvent) {
	firingKey := models.BuildProbingEventCacheKey(event.TenantId, event.RuleId)
	cache := ctx.Redis.Probing()
	event.FirstTriggerTime = cache.GetProbingEventFirstTime(firingKey)
	event.LastEvalTime = cache.GetProbingEventLastEvalTime(firingKey)
	event.LastSendTime = cache.GetProbingEventLastSendTime(firingKey)

	cache.SetProbingEventCache(event, 0)
}

func SetProbingValueMap(key models.ProbingValueCacheKey, m map[string]any) error {
	for k, v := range m {
		ctx.DO().Redis.Redis().HSet(string(key), k, v)
	}
	return nil
}

func GetProbingValueMap(key models.ProbingValueCacheKey) map[string]string {
	result := ctx.DO().Redis.Redis().HGetAll(string(key)).Val()
	return result
}
