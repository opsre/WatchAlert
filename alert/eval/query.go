package eval

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"watchAlert/alert/process"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/community/aws/cloudwatch"
	"watchAlert/pkg/community/aws/cloudwatch/types"
	"watchAlert/pkg/provider"
	"watchAlert/pkg/tools"

	"github.com/zeromicro/go-zero/core/logc"
	v1 "k8s.io/api/core/v1"
)

// Metrics Prometheus 数据源
func metrics(ctx *ctx.Context, datasourceId, datasourceType string, rule models.AlertRule) []string {
	pools := ctx.Redis.ProviderPools()
	var (
		resQuery       []provider.Metrics
		externalLabels map[string]interface{}
		// 当前活跃告警的指纹列表
		curFingerprints []string
		// 按指纹分组存储事件，相同规则只保留最高优先级的事件
		highestPriorityEvents = make(map[string]struct{})
	)

	cli, err := pools.GetClient(datasourceId)
	if err != nil {
		logc.Errorf(ctx.Ctx, "获取数据源客户端失败, 规则ID: %s, 规则名称: %s, 数据源ID: %s, 错误: %v", rule.RuleId, rule.RuleName, datasourceId, err)
		return nil
	}

	switch datasourceType {
	case provider.PrometheusDsProvider:
		resQuery, err = cli.(provider.PrometheusProvider).Query(rule.PrometheusConfig.PromQL)
		if err != nil {
			logc.Errorf(ctx.Ctx, "Prometheus查询失败, 规则ID: %s, 规则名称: %s, 数据源ID: %s, PromQL: %s, 错误: %v", rule.RuleId, rule.RuleName, datasourceId, rule.PrometheusConfig.PromQL, err)
			return nil
		}

		externalLabels = cli.(provider.PrometheusProvider).GetExternalLabels()
	default:
		logc.Errorf(ctx.Ctx, "不支持的指标类型, 规则ID: %s, 规则名称: %s, 数据源ID: %s, 类型: %s", rule.RuleId, rule.RuleName, datasourceId, datasourceType)
		return nil
	}

	if resQuery == nil {
		return nil
	}

	// 按优先级排序规则（P0 > P1 > P2）
	rules := sortRulesByPriority(rule.PrometheusConfig.Rules)

	for _, v := range resQuery {
		// 避免共享引用导致的指纹不一致问题
		metricLabels := make(map[string]interface{})
		for k, val := range v.GetMetric() {
			metricLabels[k] = val
		}

		// 使用独立的标签副本来生成指纹，避免修改原始数据
		fingerprintLabels := make(map[string]interface{})
		for k, val := range metricLabels {
			fingerprintLabels[k] = val
		}
		fingerprintLabels["rule_id"] = rule.RuleId
		fingerprintLabels["rule_name"] = rule.RuleName

		// 遍历按优先级排序后的规则
		for _, ruleExpr := range rules {
			fingerprintLabels["severity"] = ruleExpr.Severity
			operator, value, err := tools.ProcessRuleExpr(ruleExpr.Expr)
			if err != nil {
				logc.Errorf(ctx.Ctx, "处理规则表达式失败, 规则ID: %s, 规则名称: %s, 表达式: %s, 错误: %v", rule.RuleId, rule.RuleName, ruleExpr.Expr, err)
				continue
			}

			fingerprintMetric := provider.Metrics{
				Metric: fingerprintLabels,
			}
			fingerprint := fingerprintMetric.GetFingerprint()

			event := process.BuildEvent(rule, func() map[string]interface{} {
				newMetric := make(map[string]interface{})
				for k, val := range metricLabels {
					newMetric[k] = val
				}
				newMetric["rule_name"] = rule.RuleName
				newMetric["fingerprint"] = fingerprint
				newMetric["severity"] = ruleExpr.Severity
				newMetric["value"] = v.Value
				for ek, ev := range externalLabels {
					newMetric[ek] = ev
				}
				for ek, ev := range rule.ExternalLabels {
					newMetric[ek] = ev
				}

				// 获取初次触发值
				data, err := ctx.Redis.Alert().GetEventFromCache(rule.TenantId, rule.FaultCenterId, fingerprint)
				if err == nil && data.Labels["first_value"] != nil {
					newMetric["first_value"] = data.Labels["first_value"]
				} else {
					newMetric["first_value"] = v.Value
				}

				return newMetric
			})
			event.DatasourceId = datasourceId
			event.Fingerprint = fingerprint
			event.Severity = ruleExpr.Severity
			event.SearchQL = fmt.Sprintf("%s %s %v", rule.PrometheusConfig.PromQL, operator, value)
			event.ForDuration = rule.GetForDuration(ruleExpr.Severity)
			event.Annotations = tools.ParserVariables(rule.PrometheusConfig.Annotations, tools.ConvertStructToMap(event))
			event.Status = models.StatePreAlert

			// 告警评估
			if process.EvalCondition(models.EvalCondition{
				Operator:      operator,
				QueryValue:    v.Value,
				ExpectedValue: value,
			}) {
				if len(highestPriorityEvents) > 0 {
					// 如果有高优先级告警，则抑制掉低级告警
					event.LastSendTime = time.Now().Unix()
				}
				highestPriorityEvents[fingerprint] = struct{}{}
				event.Status = models.StatePreAlert
				process.PushEventToFaultCenter(ctx, &event)
				curFingerprints = append(curFingerprints, fingerprint)
			} else {
				// 更新恢复时最新值
				cache, err := ctx.Redis.Alert().GetEventFromCache(event.TenantId, event.FaultCenterId, event.Fingerprint)
				if err == nil {
					if !cache.IsRecovered && cache.Status != models.StateRecovered {
						event.Labels["value"] = v.GetValue()
						process.PushEventToFaultCenter(ctx, &event)
					}
				}
			}
		}
	}

	return curFingerprints
}

// sortRulesByPriority 按优先级排序规则
func sortRulesByPriority(rules []models.Rules) []models.Rules {
	sortedRules := make([]models.Rules, len(rules))
	copy(sortedRules, rules)

	sort.Slice(sortedRules, func(i, j int) bool {
		return getPriorityValue(sortedRules[i].Severity) > getPriorityValue(sortedRules[j].Severity)
	})

	return sortedRules
}

// getPriorityValue 获取优先级的数值表示，用于排序
// p0 优先级最高
// p1 次之
// p2 最低
// 其他情况排在后面
func getPriorityValue(severity string) int {
	switch severity {
	case "P0":
		return 3
	case "P1":
		return 2
	case "P2":
		return 1
	default:
		return 0
	}
}

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
		operator, value, err := tools.ProcessRuleExpr(rule.LogEvalCondition)
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
		operator, value, err := tools.ProcessRuleExpr(rule.LogEvalCondition)
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
		operator, value, err := tools.ProcessRuleExpr(rule.LogEvalCondition)
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
		operator, value, err := tools.ProcessRuleExpr(rule.LogEvalCondition)
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
		operator, value, err := tools.ProcessRuleExpr(rule.LogEvalCondition)
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

func cloudWatch(ctx *ctx.Context, datasourceId, datasourceType string, rule models.AlertRule) []string {
	var externalLabels map[string]interface{}
	pools := ctx.Redis.ProviderPools()
	cfg, err := pools.GetClient(datasourceId)
	if err != nil {
		logc.Errorf(ctx.Ctx, "获取CloudWatch数据源客户端失败, 规则ID: %s, 规则名称: %s, 数据源ID: %s, 错误: %v", rule.RuleId, rule.RuleName, datasourceId, err)
		return []string{}
	}

	externalLabels = cfg.(provider.AwsConfig).GetExternalLabels()

	cli := cfg.(provider.AwsConfig).CloudWatchCli()
	curAt := time.Now().UTC()
	startsAt := tools.ParserDuration(curAt, rule.CloudWatchConfig.Period, "m")

	var curFingerprints []string
	for _, endpoint := range rule.CloudWatchConfig.Endpoints {
		query := types.CloudWatchQuery{
			Endpoint:   endpoint,
			Dimension:  rule.CloudWatchConfig.Dimension,
			Period:     int32(rule.CloudWatchConfig.Period * 60),
			Namespace:  rule.CloudWatchConfig.Namespace,
			MetricName: rule.CloudWatchConfig.MetricName,
			Statistic:  rule.CloudWatchConfig.Statistic,
			Form:       startsAt,
			To:         curAt,
		}
		_, values := cloudwatch.MetricDataQuery(cli, query)
		if len(values) == 0 {
			return []string{}
		}

		event := process.BuildEvent(rule, func() map[string]interface{} {
			metric := query.GetMetrics()
			metric["severity"] = rule.Severity
			for ek, ev := range externalLabels {
				metric[ek] = ev
			}
			for ek, ev := range rule.ExternalLabels {
				metric[ek] = ev
			}
			metric["rule_name"] = rule.RuleName
			return metric
		})
		event.DatasourceId = datasourceId
		event.Fingerprint = query.GetFingerprint()
		event.Annotations = fmt.Sprintf("%s %s %s %s %d", query.Namespace, query.MetricName, query.Statistic, rule.CloudWatchConfig.Expr, rule.CloudWatchConfig.Threshold)

		options := models.EvalCondition{
			Operator:      rule.CloudWatchConfig.Expr,
			QueryValue:    values[0],
			ExpectedValue: float64(rule.CloudWatchConfig.Threshold),
		}

		curFingerprints = append(curFingerprints, event.Fingerprint)
		if process.EvalCondition(options) {
			process.PushEventToFaultCenter(ctx, &event)
		}
	}

	return curFingerprints
}

func kubernetesEvent(ctx *ctx.Context, datasourceId, datasourceType string, rule models.AlertRule) []string {
	var externalLabels map[string]interface{}
	datasourceObj, err := ctx.DB.Datasource().GetInstance(datasourceId)
	if err != nil {
		logc.Errorf(ctx.Ctx, "获取数据源实例失败, 规则ID: %s, 规则名称: %s, 数据源ID: %s, 错误: %v", rule.RuleId, rule.RuleName, datasourceId, err)
		return []string{}
	}

	pools := ctx.Redis.ProviderPools()
	cli, err := pools.GetClient(datasourceId)
	if err != nil {
		logc.Errorf(ctx.Ctx, "获取Kubernetes数据源客户端失败, 规则ID: %s, 规则名称: %s, 数据源ID: %s, 错误: %v", rule.RuleId, rule.RuleName, datasourceId, err)
		return []string{}
	}

	k8sEvent, err := cli.(provider.KubernetesClient).GetWarningEvent(rule.KubernetesConfig.Reason, rule.KubernetesConfig.Scope)
	if err != nil {
		logc.Errorf(ctx.Ctx, "获取Kubernetes警告事件失败, 规则ID: %s, 规则名称: %s, 数据源ID: %s, 资源: %s, 错误: %v", rule.RuleId, rule.RuleName, datasourceId, rule.KubernetesConfig.Resource, err)
		return []string{}
	}

	externalLabels = cli.(provider.KubernetesClient).GetExternalLabels()

	if len(k8sEvent.Items) == 0 {
		return []string{}
	}

	// 分组：key = resourceName + eventReason
	groupedEvents := make(map[string][]v1.Event)

	// 过滤并分组事件
	for _, item := range process.FilterKubeEvent(k8sEvent, rule.KubernetesConfig.Filter).Items {
		key := fmt.Sprintf("%s/%s", item.InvolvedObject.Name, item.Reason)
		groupedEvents[key] = append(groupedEvents[key], item)
	}

	var curFingerprints []string

	for _, items := range groupedEvents {
		// 不满足阈值，跳过
		if len(items) < rule.KubernetesConfig.Value {
			continue
		}

		// 取第一个作为代表生成告警
		item := items[0]

		// 构造告警内容
		k8sItem := process.KubernetesAlertEvent(ctx, item)
		fingerprint := k8sItem.GetFingerprint()
		event := process.BuildEvent(rule, func() map[string]interface{} {
			metric := k8sItem.GetMetrics()
			metric["rule_name"] = rule.RuleName
			metric["severity"] = rule.Severity
			metric["fingerprint"] = fingerprint
			metric["value"] = len(items)
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
		event.SearchQL = rule.KubernetesConfig.Resource

		// 拼接注释信息
		var msgList []string
		for _, e := range items {
			msg := strings.ReplaceAll(e.Message, "\"", "'")
			msgList = append(msgList, msg)
		}
		event.Annotations = fmt.Sprintf(
			"- 数据源: %s\n- 命名空间: %s\n- 资源类型: %s\n- 资源名称: %s\n- 事件类型: %s\n- 事件详情:\n%s",
			datasourceObj.Name,
			item.Namespace,
			item.InvolvedObject.Kind,
			item.InvolvedObject.Name,
			item.Reason,
			strings.Join(msgList, "\n"),
		)

		curFingerprints = append(curFingerprints, event.Fingerprint)
		process.PushEventToFaultCenter(ctx, &event)
	}

	return curFingerprints
}
