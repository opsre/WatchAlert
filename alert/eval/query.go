package eval

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	v1 "k8s.io/api/core/v1"
	"sort"
	"strings"
	"time"
	"watchAlert/alert/process"
	"watchAlert/internal/models"
	"watchAlert/pkg/community/aws/cloudwatch"
	"watchAlert/pkg/community/aws/cloudwatch/types"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/provider"
	"watchAlert/pkg/tools"
)

// Metrics 包含 Prometheus、VictoriaMetrics 数据源
func metrics(ctx *ctx.Context, datasourceId, datasourceType string, rule models.AlertRule) []string {
	pools := ctx.Redis.ProviderPools()
	var (
		resQuery       []provider.Metrics
		externalLabels map[string]interface{}
		// 当前活跃告警的指纹列表
		curFingerprints []string
		// 按指纹分组存储事件，每个指纹只保留最高优先级的事件
		highestPriorityEvents = make(map[string]models.AlertCurEvent)
	)

	cli, err := pools.GetClient(datasourceId)
	if err != nil {
		logc.Errorf(ctx.Ctx, err.Error())
		return nil
	}

	switch datasourceType {
	case provider.PrometheusDsProvider:
		resQuery, err = cli.(provider.PrometheusProvider).Query(rule.PrometheusConfig.PromQL)
		if err != nil {
			logc.Error(ctx.Ctx, err.Error())
			return nil
		}

		externalLabels = cli.(provider.PrometheusProvider).GetExternalLabels()
	case provider.VictoriaMetricsDsProvider:
		resQuery, err = cli.(provider.VictoriaMetricsProvider).Query(rule.PrometheusConfig.PromQL)
		if err != nil {
			logc.Error(ctx.Ctx, err.Error())
			return nil
		}

		externalLabels = cli.(provider.VictoriaMetricsProvider).GetExternalLabels()
	default:
		logc.Errorf(ctx.Ctx, fmt.Sprintf("Unsupported metrics type, type: %s", datasourceType))
		return nil
	}

	if resQuery == nil {
		return nil
	}

	// 获取已缓存事件指纹
	fingerPrintMap := process.GetFingerPrint(ctx, rule.TenantId, rule.FaultCenterId, rule.RuleId)

	// 按优先级排序规则（P0 > P1 > P2）
	rules := sortRulesByPriority(rule.PrometheusConfig.Rules)

	for _, v := range resQuery {
		// 遍历按优先级排序后的规则
		for _, ruleExpr := range rules {
			operator, value, err := tools.ProcessRuleExpr(ruleExpr.Expr)
			if err != nil {
				logc.Errorf(ctx.Ctx, err.Error())
				continue
			}

			option := models.EvalCondition{
				Operator:      operator,
				QueryValue:    v.Value,
				ExpectedValue: value,
			}

			fingerprint := v.GetFingerprint()
			event := process.BuildEvent(rule, func() map[string]interface{} {
				metric := v.GetMetric()
				metric["rule_name"] = rule.RuleName
				metric["fingerprint"] = fingerprint
				metric["severity"] = ruleExpr.Severity
				metric["value"] = v.Value
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
			event.Severity = ruleExpr.Severity
			event.SearchQL = rule.PrometheusConfig.PromQL
			event.Annotations = tools.ParserVariables(rule.PrometheusConfig.Annotations, tools.ConvertEventToMap(event))

			// 告警评估
			if process.EvalCondition(option) {
				// 如果条件满足，检查是否已经有更高优先级的事件
				if _, exists := highestPriorityEvents[fingerprint]; !exists {
					// 如果该指纹还没有事件，添加当前事件
					event.Status = models.StatePreAlert
					highestPriorityEvents[fingerprint] = event
					curFingerprints = append(curFingerprints, fingerprint)
					if _, e := event.Labels["recover_value"]; e {
						delete(event.Labels, "recover_value")
					}
				}
				// 找到符合条件的规则后，跳过该指标的其他规则
				break
			} else if _, exist := fingerPrintMap[fingerprint]; exist {
				// 获取过恢复值则直接跳过
				if _, existRecoverValue := event.Labels["recover_value"]; existRecoverValue {
					continue
				}

				// 获取上一次告警值
				event.Labels["value"] = ctx.Redis.Alert().GetLastFiringValue(event.TenantId, event.FaultCenterId, event.Fingerprint)
				// 获取当前恢复值
				event.Labels["recover_value"] = v.GetValue()
				process.PushEventToFaultCenter(ctx, &event)
			}
		}
	}

	// 推送最高优先级的事件
	for _, event := range highestPriorityEvents {
		process.PushEventToFaultCenter(ctx, &event)
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
		logc.Errorf(ctx.Ctx, err.Error())
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
			logc.Error(ctx.Ctx, err.Error())
			return []string{}
		}

		externalLabels = cli.(provider.LokiProvider).GetExternalLabels()
		operator, value, err := tools.ProcessRuleExpr(rule.LogEvalCondition)
		if err != nil {
			logc.Errorf(ctx.Ctx, err.Error())
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
			logc.Error(ctx.Ctx, err.Error())
			return []string{}
		}

		externalLabels = cli.(provider.AliCloudSlsDsProvider).GetExternalLabels()
		operator, value, err := tools.ProcessRuleExpr(rule.LogEvalCondition)
		if err != nil {
			logc.Errorf(ctx.Ctx, err.Error())
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
			logc.Error(ctx.Ctx, err.Error())
			return []string{}
		}

		externalLabels = cli.(provider.ElasticSearchDsProvider).GetExternalLabels()
		operator, value, err := tools.ProcessRuleExpr(rule.LogEvalCondition)
		if err != nil {
			logc.Errorf(ctx.Ctx, err.Error())
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
			logc.Error(ctx.Ctx, err.Error())
			return []string{}
		}

		externalLabels = cli.(provider.VictoriaLogsProvider).GetExternalLabels()
		operator, value, err := tools.ProcessRuleExpr(rule.LogEvalCondition)
		if err != nil {
			logc.Errorf(ctx.Ctx, err.Error())
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
			logc.Error(ctx.Ctx, err.Error())
			return []string{}
		}

		externalLabels = cli.(provider.ClickHouseProvider).GetExternalLabels()
		operator, value, err := tools.ProcessRuleExpr(rule.LogEvalCondition)
		if err != nil {
			logc.Errorf(ctx.Ctx, err.Error())
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
				event.SearchQL = tools.JsonMarshal(rule.ElasticSearchConfig.Filter)
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
			logc.Errorf(ctx.Ctx, err.Error())
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
			logc.Error(ctx.Ctx, err.Error())
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

func cloudWatch(ctx *ctx.Context, datasourceId string, rule models.AlertRule) []string {
	var externalLabels map[string]interface{}
	pools := ctx.Redis.ProviderPools()
	cfg, err := pools.GetClient(datasourceId)
	if err != nil {
		logc.Errorf(ctx.Ctx, err.Error())
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

func kubernetesEvent(ctx *ctx.Context, datasourceId string, rule models.AlertRule) []string {
	var externalLabels map[string]interface{}
	datasourceObj, err := ctx.DB.Datasource().GetInstance(datasourceId)
	if err != nil {
		logc.Error(ctx.Ctx, err.Error())
		return []string{}
	}

	pools := ctx.Redis.ProviderPools()
	cli, err := pools.GetClient(datasourceId)
	if err != nil {
		logc.Errorf(ctx.Ctx, err.Error())
		return []string{}
	}

	k8sEvent, err := cli.(provider.KubernetesClient).GetWarningEvent(rule.KubernetesConfig.Reason, rule.KubernetesConfig.Scope)
	if err != nil {
		logc.Error(ctx.Ctx, err.Error())
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
