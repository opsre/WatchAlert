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

// traces evaluates Jaeger and SkyWalking trace datasources.
func traces(ctx *ctx.Context, datasourceId, datasourceType string, rule models.AlertRule) []string {
	var (
		queryRes       []provider.Traces
		externalLabels map[string]interface{}
		service        string
		tags           string
		scope          int
	)

	switch datasourceType {
	case provider.JaegerDsProviderName:
		service = rule.JaegerConfig.Service
		tags = rule.JaegerConfig.Tags
		scope = rule.JaegerConfig.Scope
	case provider.SkyWalkingDsProviderName:
		service = rule.SkyWalkingConfig.Service
		tags = rule.SkyWalkingConfig.Tags
		scope = rule.SkyWalkingConfig.Scope
	default:
		return []string{}
	}

	curAt := time.Now()
	startsAt := tools.ParserDuration(curAt, scope, "m")
	cli, err := ctx.Redis.ProviderPools().GetClient(datasourceId)
	if err != nil {
		logc.Errorf(ctx.Ctx, "获取链路数据源客户端失败, 规则ID: %s, 规则名称: %s, 数据源ID: %s, 错误: %v", rule.RuleId, rule.RuleName, datasourceId, err)
		return []string{}
	}

	traceClient, ok := cli.(provider.TracesFactoryProvider)
	if !ok {
		logc.Errorf(ctx.Ctx, "链路数据源客户端类型错误, 规则ID: %s, 数据源ID: %s", rule.RuleId, datasourceId)
		return []string{}
	}

	queryRes, err = traceClient.Query(provider.TraceQueryOptions{
		Tags:    tags,
		Service: service,
		StartAt: startsAt.UnixMicro(),
		EndAt:   curAt.UnixMicro(),
	})
	if err != nil {
		logc.Errorf(ctx.Ctx, "%s查询失败, 规则ID: %s, 规则名称: %s, 数据源ID: %s, 服务: %s, 错误: %v", datasourceType, rule.RuleId, rule.RuleName, datasourceId, service, err)
		return []string{}
	}
	externalLabels = traceClient.GetExternalLabels()

	var curFingerprints []string
	for _, v := range queryRes {
		fingerprint := v.GetFingerprint()
		event := process.BuildEvent(rule, func() map[string]interface{} {
			metric := v.GetMetric()
			metric["rule_name"] = rule.RuleName
			metric["severity"] = rule.Severity
			metric["fingerprint"] = fingerprint
			metric["service"] = service
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
		event.SearchQL = tags
		event.Annotations = fmt.Sprintf("服务: %s 链路中存在异常, TraceId: %s", service, v.TraceId)

		curFingerprints = append(curFingerprints, event.Fingerprint)
		process.PushEventToFaultCenter(ctx, &event)
	}

	return curFingerprints
}
