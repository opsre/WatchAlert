package eval

import (
	"fmt"
	"time"
	"watchAlert/alert/process"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/provider"
	"watchAlert/pkg/tools"

	"github.com/zeromicro/go-zero/core/logc"
)

// Traces 包含 Jaeger 数据源
func traces(ctx *ctx.Context, datasourceId, datasourceType string, rule models.AlertRule) []string {
	var (
		queryRes       []provider.Traces
		externalLabels map[string]interface{}
	)

	pools := ctx.Redis.ProviderPools()
	switch datasourceType {
	case provider.JaegerDsProviderName:
		curAt := time.Now().UTC()
		startsAt := tools.ParserDuration(curAt, rule.JaegerConfig.Scope, "m")

		cli, err := pools.GetClient(datasourceId)
		if err != nil {
			logc.Errorf(ctx.Ctx, "获取Jaeger数据源客户端失败, 规则ID: %s, 规则名称: %s, 数据源ID: %s, 错误: %v", rule.RuleId, rule.RuleName, datasourceId, err)
			return []string{}
		}

		queryOptions := provider.TraceQueryOptions{
			Tags:    rule.JaegerConfig.Tags,
			Service: rule.JaegerConfig.Service,
			StartAt: startsAt.UnixMicro(),
			EndAt:   curAt.UnixMicro(),
		}
		queryRes, err = cli.(provider.JaegerDsProvider).Query(queryOptions)
		if err != nil {
			logc.Errorf(ctx.Ctx, "Jaeger查询失败, 规则ID: %s, 规则名称: %s, 数据源ID: %s, 服务: %s, 错误: %v", rule.RuleId, rule.RuleName, datasourceId, rule.JaegerConfig.Service, err)
			return []string{}
		}

		externalLabels = cli.(provider.JaegerDsProvider).GetExternalLabels()
	}

	var curFingerprints []string
	for _, v := range queryRes {
		fingerprint := v.GetFingerprint()
		event := process.BuildEvent(rule, func() map[string]interface{} {
			metric := v.GetMetric()
			metric["rule_name"] = rule.RuleName
			metric["severity"] = rule.Severity
			metric["fingerprint"] = fingerprint
			metric["service"] = rule.JaegerConfig.Service
			metric["traceId"] = v.TraceId
			for ek, ev := range externalLabels {
				metric[ek] = ev
			}
			for ek, ev := range rule.ExternalLabels {
				metric[ek] = ev
			}
			return metric
		})
		event.DatasourceId = datasourceId
		event.Fingerprint = fingerprint
		event.SearchQL = rule.JaegerConfig.Tags
		event.Annotations = fmt.Sprintf("服务: %s 链路中存在异常, TraceId: %s", rule.JaegerConfig.Service, v.TraceId)

		curFingerprints = append(curFingerprints, event.Fingerprint)
		process.PushEventToFaultCenter(ctx, &event)
	}

	return curFingerprints
}
