package probing

import (
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
)

func (t *ProductProbing) buildEvent(rule models.ProbingRule) *models.ProbingEvent {
	return &models.ProbingEvent{
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
