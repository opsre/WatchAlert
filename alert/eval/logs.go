package eval

import (
	"time"
	"watchAlert/alert/process"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/provider"
	"watchAlert/pkg/tools"

	"github.com/zeromicro/go-zero/core/logc"
)

// Logs 包含 AliSLS、Loki、ElasticSearch 数据源
func logs(ctx *ctx.Context, datasourceId, datasourceType string, rule models.AlertRule) []string {
	var (
		// 日志信息
		log provider.Logs
		// 日志总数
		count int
		// 评估
		evalOptions models.EvalCondition
		// 额外的标签
		externalLabels map[string]interface{}
		// 当前时间
		curAt = time.Now()
	)

	pools := ctx.Redis.ProviderPools()
	cli, err := pools.GetClient(datasourceId)
	if err != nil {
		logc.Errorf(ctx.Ctx, "获取数据源客户端失败, 规则ID: %s, 规则名称: %s, 数据源ID: %s, 错误: %v", rule.RuleId, rule.RuleName, datasourceId, err)
		return []string{}
	}

	switch datasourceType {
	case provider.LokiDsProviderName:
		startsAt := tools.ParserDuration(curAt, rule.LokiConfig.LogScope, "m")
		queryOptions := provider.LogQueryOptions{
			Loki: provider.Loki{
				Query: rule.LokiConfig.LogQL,
			},
			StartAt: startsAt.Unix(),
			EndAt:   curAt.Unix(),
		}
		log, count, err = cli.(provider.LokiProvider).Query(queryOptions)
		if err != nil {
			logc.Errorf(ctx.Ctx, "Loki查询失败, 规则ID: %s, 规则名称: %s, 数据源ID: %s, LogQL: %s, 错误: %v", rule.RuleId, rule.RuleName, datasourceId, rule.LokiConfig.LogQL, err)
			return []string{}
		}

		externalLabels = cli.(provider.LokiProvider).GetExternalLabels()
		operator, value, err := process.ProcessRuleExpr(rule.LogEvalCondition)
		if err != nil {
			logc.Errorf(ctx.Ctx, "处理日志规则表达式失败, 规则ID: %s, 规则名称: %s, 表达式: %s, 错误: %v", rule.RuleId, rule.RuleName, rule.LogEvalCondition, err)
			return []string{}
		}

		evalOptions = models.EvalCondition{
			Operator:      operator,
			QueryValue:    float64(count),
			ExpectedValue: value,
		}
	case provider.AliCloudSLSDsProviderName:
		startsAt := tools.ParserDuration(curAt, rule.AliCloudSLSConfig.LogScope, "m")
		queryOptions := provider.LogQueryOptions{
			AliCloudSLS: provider.AliCloudSLS{
				Query:    rule.AliCloudSLSConfig.LogQL,
				Project:  rule.AliCloudSLSConfig.Project,
				LogStore: rule.AliCloudSLSConfig.Logstore,
			},
			StartAt: int32(startsAt.Unix()),
			EndAt:   int32(curAt.Unix()),
		}
		log, count, err = cli.(provider.AliCloudSlsDsProvider).Query(queryOptions)
		if err != nil {
			logc.Errorf(ctx.Ctx, "AliCloudSLS查询失败, 规则ID: %s, 规则名称: %s, 数据源ID: %s, LogQL: %s, 错误: %v", rule.RuleId, rule.RuleName, datasourceId, rule.AliCloudSLSConfig.LogQL, err)
			return []string{}
		}

		externalLabels = cli.(provider.AliCloudSlsDsProvider).GetExternalLabels()
		operator, value, err := process.ProcessRuleExpr(rule.LogEvalCondition)
		if err != nil {
			logc.Errorf(ctx.Ctx, "处理日志规则表达式失败, 规则ID: %s, 规则名称: %s, 表达式: %s, 错误: %v", rule.RuleId, rule.RuleName, rule.LogEvalCondition, err)
			return []string{}
		}

		evalOptions = models.EvalCondition{
			Operator:      operator,
			QueryValue:    float64(count),
			ExpectedValue: value,
		}
	case provider.ElasticSearchDsProviderName:
		queryOptions := provider.LogQueryOptions{
			ElasticSearch: provider.Elasticsearch{
				Index:                rule.ElasticSearchConfig.Index,
				QueryFilter:          rule.ElasticSearchConfig.Filter,
				QueryFilterCondition: rule.ElasticSearchConfig.FilterCondition,
				QueryType:            rule.ElasticSearchConfig.EsQueryType,
				QueryWildcard:        rule.ElasticSearchConfig.QueryWildcard,
				RawJson:              rule.ElasticSearchConfig.RawJson,
			},
		}
		log, count, err = cli.(provider.ElasticSearchDsProvider).Query(queryOptions)
		if err != nil {
			logc.Errorf(ctx.Ctx, "ElasticSearch查询失败, 规则ID: %s, 规则名称: %s, 数据源ID: %s, 索引: %s, 错误: %v", rule.RuleId, rule.RuleName, datasourceId, rule.ElasticSearchConfig.Index, err)
			return []string{}
		}

		externalLabels = cli.(provider.ElasticSearchDsProvider).GetExternalLabels()
		operator, value, err := process.ProcessRuleExpr(rule.LogEvalCondition)
		if err != nil {
			logc.Errorf(ctx.Ctx, "处理日志规则表达式失败, 规则ID: %s, 规则名称: %s, 表达式: %s, 错误: %v", rule.RuleId, rule.RuleName, rule.LogEvalCondition, err)
			return []string{}
		}

		evalOptions = models.EvalCondition{
			Operator:      operator,
			QueryValue:    float64(count),
			ExpectedValue: value,
		}
	case provider.VictoriaLogsDsProviderName:
		startsAt := tools.ParserDuration(curAt, rule.VictoriaLogsConfig.LogScope, "m")
		queryOptions := provider.LogQueryOptions{
			VictoriaLogs: provider.VictoriaLogs{
				Query: rule.VictoriaLogsConfig.LogQL,
				Limit: rule.VictoriaLogsConfig.Limit,
			},
			StartAt: int32(startsAt.Unix()),
			EndAt:   int32(curAt.Unix()),
		}
		log, count, err = cli.(provider.VictoriaLogsProvider).Query(queryOptions)
		if err != nil {
			logc.Errorf(ctx.Ctx, "VictoriaLogs查询失败, 规则ID: %s, 规则名称: %s, 数据源ID: %s, LogQL: %s, 错误: %v", rule.RuleId, rule.RuleName, datasourceId, rule.VictoriaLogsConfig.LogQL, err)
			return []string{}
		}

		externalLabels = cli.(provider.VictoriaLogsProvider).GetExternalLabels()
		operator, value, err := process.ProcessRuleExpr(rule.LogEvalCondition)
		if err != nil {
			logc.Errorf(ctx.Ctx, "处理日志规则表达式失败, 规则ID: %s, 规则名称: %s, 表达式: %s, 错误: %v", rule.RuleId, rule.RuleName, rule.LogEvalCondition, err)
			return []string{}
		}

		evalOptions = models.EvalCondition{
			Operator:      operator,
			QueryValue:    float64(count),
			ExpectedValue: value,
		}
	case provider.ClickHouseDsProviderName:
		queryOptions := provider.LogQueryOptions{
			ClickHouse: provider.ClickHouse{
				Query: rule.ClickHouseConfig.LogQL,
			},
		}
		log, count, err = cli.(provider.ClickHouseProvider).Query(queryOptions)
		if err != nil {
			logc.Errorf(ctx.Ctx, "ClickHouse查询失败, 规则ID: %s, 规则名称: %s, 数据源ID: %s, LogQL: %s, 错误: %v", rule.RuleId, rule.RuleName, datasourceId, rule.ClickHouseConfig.LogQL, err)
			return []string{}
		}

		externalLabels = cli.(provider.ClickHouseProvider).GetExternalLabels()
		operator, value, err := process.ProcessRuleExpr(rule.LogEvalCondition)
		if err != nil {
			logc.Errorf(ctx.Ctx, "处理日志规则表达式失败, 规则ID: %s, 规则名称: %s, 表达式: %s, 错误: %v", rule.RuleId, rule.RuleName, rule.LogEvalCondition, err)
			return []string{}
		}

		evalOptions = models.EvalCondition{
			Operator:      operator,
			QueryValue:    float64(count),
			ExpectedValue: value,
		}
	}

	if count <= 0 {
		return []string{}
	}

	// 唯一指纹基于 RuleId
	fingerprint := log.GenerateFingerprint(rule.RuleId)
	var curFingerprints []string
	event := func() *models.AlertCurEvent {
		event := process.BuildEvent(rule, func() map[string]interface{} {
			labels := map[string]interface{}{
				"value":       count,
				"severity":    rule.Severity,
				"fingerprint": fingerprint,
				"rule_name":   rule.RuleName,
			}
			for ek, ev := range externalLabels {
				labels[ek] = ev
			}
			for ek, ev := range rule.ExternalLabels {
				labels[ek] = ev
			}
			for logKey, logValue := range log.GetAnnotations() {
				labels[logKey] = logValue
			}
			return labels
		})
		event.DatasourceId = datasourceId
		event.Fingerprint = fingerprint

		switch datasourceType {
		case provider.LokiDsProviderName:
			event.SearchQL = rule.LokiConfig.LogQL
		case provider.AliCloudSLSDsProviderName:
			event.SearchQL = rule.AliCloudSLSConfig.LogQL
		case provider.ElasticSearchDsProviderName:
			if rule.ElasticSearchConfig.RawJson != "" {
				event.SearchQL = rule.ElasticSearchConfig.RawJson
			} else {
				event.SearchQL = tools.JsonMarshalToString(rule.ElasticSearchConfig.Filter)
			}
		case provider.VictoriaLogsDsProviderName:
			event.SearchQL = rule.VictoriaLogsConfig.LogQL
		}

		curFingerprints = append(curFingerprints, event.Fingerprint)

		return &event
	}

	// 评估告警条件
	if process.EvalCondition(evalOptions) {
		process.PushEventToFaultCenter(ctx, event())
	}

	return curFingerprints
}
