package eval

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
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
	pools := ctx.Cache.ProviderPools()
	var (
		resQuery       []provider.Metrics
		externalLabels map[string]interface{}
		// 当前活跃告警的指纹列表
		curFingerprints []string
		// 按指纹分组存储事件，每个指纹只保留最高优先级的事件
		highestPriorityEvents = make(map[string]models.AlertCurEvent)
	)
	switch datasourceType {
	case provider.PrometheusDsProvider:
		cli, err := pools.GetClient(datasourceId)
		if err != nil {
			logc.Errorf(ctx.Ctx, err.Error())
			return nil
		}

		resQuery, err = cli.(provider.PrometheusProvider).Query(rule.PrometheusConfig.PromQL)
		if err != nil {
			logc.Error(ctx.Ctx, err.Error())
			return nil
		}

		externalLabels = cli.(provider.PrometheusProvider).GetExternalLabels()
	case provider.VictoriaMetricsDsProvider:
		cli, err := pools.GetClient(datasourceId)
		if err != nil {
			logc.Errorf(ctx.Ctx, err.Error())
			return nil
		}

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
		fingerprint := v.GetFingerprint()

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

			event := process.BuildEvent(rule)
			event.DatasourceId = datasourceId
			event.Fingerprint = fingerprint
			event.Metric = *v.GetMetric()
			event.Severity = ruleExpr.Severity
			event.Metric["severity"] = ruleExpr.Severity
			for ek, ev := range externalLabels {
				event.Metric[ek] = ev
			}
			event.Annotations = tools.ParserVariables(rule.PrometheusConfig.Annotations, event.Metric)
			event.SearchQL = rule.PrometheusConfig.PromQL

			if process.EvalCondition(option) {
				// 如果条件满足，检查是否已经有更高优先级的事件
				if _, exists := highestPriorityEvents[fingerprint]; !exists {
					// 如果该指纹还没有事件，添加当前事件
					highestPriorityEvents[fingerprint] = event
					curFingerprints = append(curFingerprints, fingerprint)
				}
				// 找到符合条件的规则后，跳过该指标的其他规则
				break
			} else if _, exist := fingerPrintMap[fingerprint]; exist {
				// 如果是 预告警 状态的事件，触发了恢复逻辑，但它并非是真正触发告警而恢复，所以只需要删除历史事件即可，无需继续处理恢复逻辑。
				if ctx.Cache.Event().GetEventStatusForFaultCenter(event.TenantId, event.FaultCenterId, fingerprint) == 0 {
					logc.Alert(ctx.Ctx, fmt.Sprintf("移除预告警恢复事件, Rule: %s, Fingerprint: %s", rule.RuleName, fingerprint))
					ctx.Cache.Event().RemoveEventFromFaultCenter(event.TenantId, event.FaultCenterId, fingerprint)
					continue
				}

				// 获取过恢复值则直接跳过
				if _, existRecoverValue := event.Metric["recover_value"]; existRecoverValue {
					continue
				}

				// 获取上一次告警值
				event.Metric["value"] = ctx.Cache.Event().GetLastFiringValueForFaultCenter(event.TenantId, event.FaultCenterId, event.Fingerprint)
				// 获取当前恢复值
				event.Metric["recover_value"] = v.GetValue()
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
		queryRes       []provider.Logs
		count          int
		evalOptions    models.EvalCondition
		externalLabels map[string]interface{}
	)

	pools := ctx.Cache.ProviderPools()
	switch datasourceType {
	case provider.LokiDsProviderName:
		cli, err := pools.GetClient(datasourceId)
		if err != nil {
			logc.Errorf(ctx.Ctx, err.Error())
			return []string{}
		}

		curAt := time.Now()
		startsAt := tools.ParserDuration(curAt, rule.LokiConfig.LogScope, "m")
		queryOptions := provider.LogQueryOptions{
			Loki: provider.Loki{
				Query: rule.LokiConfig.LogQL,
			},
			StartAt: startsAt.Unix(),
			EndAt:   curAt.Unix(),
		}
		queryRes, count, err = cli.(provider.LokiProvider).Query(queryOptions)
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
		cli, err := pools.GetClient(datasourceId)
		if err != nil {
			logc.Errorf(ctx.Ctx, err.Error())
			return []string{}
		}

		curAt := time.Now()
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
		queryRes, count, err = cli.(provider.AliCloudSlsDsProvider).Query(queryOptions)
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
		cli, err := pools.GetClient(datasourceId)
		if err != nil {
			logc.Errorf(ctx.Ctx, err.Error())
			return []string{}
		}

		curAt := time.Now()
		startsAt := tools.ParserDuration(curAt, int(rule.ElasticSearchConfig.Scope), "m")
		queryOptions := provider.LogQueryOptions{
			ElasticSearch: provider.Elasticsearch{
				Index:                rule.ElasticSearchConfig.Index,
				QueryFilter:          rule.ElasticSearchConfig.Filter,
				QueryFilterCondition: rule.ElasticSearchConfig.FilterCondition,
				QueryType:            rule.ElasticSearchConfig.EsQueryType,
				QueryWildcard:        rule.ElasticSearchConfig.QueryWildcard,
				RawJson:              rule.ElasticSearchConfig.RawJson,
			},
			StartAt: tools.FormatTimeToUTC(startsAt.Unix()),
			EndAt:   tools.FormatTimeToUTC(curAt.Unix()),
		}
		queryRes, count, err = cli.(provider.ElasticSearchDsProvider).Query(queryOptions)
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
		cli, err := pools.GetClient(datasourceId)
		if err != nil {
			logc.Errorf(ctx.Ctx, err.Error())
			return []string{}
		}

		curAt := time.Now()
		startsAt := tools.ParserDuration(curAt, rule.VictoriaLogsConfig.LogScope, "m")
		queryOptions := provider.LogQueryOptions{
			VictoriaLogs: provider.VictoriaLogs{
				Query: rule.VictoriaLogsConfig.LogQL,
				Limit: rule.VictoriaLogsConfig.Limit,
			},
			StartAt: int32(startsAt.Unix()),
			EndAt:   int32(curAt.Unix()),
		}
		queryRes, count, err = cli.(provider.VictoriaLogsProvider).Query(queryOptions)
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
	}

	if count <= 0 {
		return []string{}
	}

	var curFingerprints []string
	for _, v := range queryRes {
		event := func() *models.AlertCurEvent {
			event := process.BuildEvent(rule)
			event.DatasourceId = datasourceId
			event.Fingerprint = v.GetFingerprint()
			event.Metric = v.GetMetric()
			event.Metric["value"] = count
			for ek, ev := range externalLabels {
				event.Metric[ek] = ev
			}
			event.Annotations = fmt.Sprintf("共计日志 %d 条\n%s", count, tools.FormatJson(v.GetAnnotations()[0].(string)))

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
	}

	return curFingerprints
}

// Traces 包含 Jaeger 数据源
func traces(ctx *ctx.Context, datasourceId, datasourceType string, rule models.AlertRule) []string {
	var (
		queryRes       []provider.Traces
		externalLabels map[string]interface{}
	)

	pools := ctx.Cache.ProviderPools()
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
		event := process.BuildEvent(rule)
		event.DatasourceId = datasourceId
		event.Fingerprint = v.GetFingerprint()
		event.Metric = v.GetMetric()
		for ek, ev := range externalLabels {
			event.Metric[ek] = ev
		}
		event.SearchQL = rule.JaegerConfig.Tags
		event.Annotations = fmt.Sprintf("服务: %s 链路中存在异常状态码接口, TraceId: %s", rule.JaegerConfig.Service, v.TraceId)

		curFingerprints = append(curFingerprints, event.Fingerprint)
		process.PushEventToFaultCenter(ctx, &event)
	}

	return curFingerprints
}

func cloudWatch(ctx *ctx.Context, datasourceId string, rule models.AlertRule) []string {
	var externalLabels map[string]interface{}
	pools := ctx.Cache.ProviderPools()
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

		event := process.BuildEvent(rule)
		event.DatasourceId = datasourceId
		event.Fingerprint = query.GetFingerprint()
		event.Metric = query.GetMetrics()
		for ek, ev := range externalLabels {
			event.Metric[ek] = ev
		}
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

	pools := ctx.Cache.ProviderPools()
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

	if len(k8sEvent.Items) < rule.KubernetesConfig.Value {
		return []string{}
	}

	var eventMapping = make(map[string][]string)
	var curFingerprints []string
	for _, item := range process.FilterKubeEvent(k8sEvent, rule.KubernetesConfig.Filter).Items {
		// 同一个资源可能有多条不同的事件信息
		eventMapping[item.InvolvedObject.Name] = append(eventMapping[item.InvolvedObject.Name], "\n"+strings.ReplaceAll(item.Message, "\"", "'"))
		k8sItem := process.KubernetesAlertEvent(ctx, item)
		event := process.BuildEvent(rule)
		event.DatasourceId = datasourceId
		event.Fingerprint = k8sItem.GetFingerprint()
		event.Metric = k8sItem.GetMetrics()
		for ek, ev := range externalLabels {
			event.Metric[ek] = ev
		}
		event.SearchQL = rule.KubernetesConfig.Resource
		event.Annotations = fmt.Sprintf("- 环境: %s\n- 命名空间: %s\n- 资源类型: %s\n- 资源名称: %s\n- 事件类型: %s\n- 事件详情: %s\n",
			datasourceObj.Name, item.Namespace, item.InvolvedObject.Kind,
			item.InvolvedObject.Name, item.Reason, eventMapping[item.InvolvedObject.Name],
		)

		curFingerprints = append(curFingerprints, event.Fingerprint)
		process.PushEventToFaultCenter(ctx, &event)
	}

	return curFingerprints
}
