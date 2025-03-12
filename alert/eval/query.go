package eval

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"strings"
	"time"
	"watchAlert/alert/process"
	models "watchAlert/internal/models"
	"watchAlert/pkg/community/aws/cloudwatch"
	"watchAlert/pkg/community/aws/cloudwatch/types"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/provider"
	"watchAlert/pkg/tools"
)

// Metrics 包含 Prometheus、VictoriaMetrics 数据源
func metrics(ctx *ctx.Context, datasourceId, datasourceType string, rule models.AlertRule) ([]string, map[string]interface{}) {
	pools := ctx.Redis.ProviderPools()
	var (
		resQuery       []provider.Metrics
		externalLabels map[string]interface{}
	)
	fingerPrintMetrics := make(map[string]interface{})
	switch datasourceType {
	case provider.PrometheusDsProvider:
		cli, err := pools.GetClient(datasourceId)
		if err != nil {
			logc.Errorf(ctx.Ctx, err.Error())
			return nil, nil
		}

		resQuery, err = cli.(provider.PrometheusProvider).Query(rule.PrometheusConfig.PromQL)
		if err != nil {
			logc.Error(ctx.Ctx, err.Error())
			return nil, nil
		}

		externalLabels = cli.(provider.PrometheusProvider).GetExternalLabels()
	case provider.VictoriaMetricsDsProvider:
		cli, err := pools.GetClient(datasourceId)
		if err != nil {
			logc.Errorf(ctx.Ctx, err.Error())
			return nil, nil
		}

		resQuery, err = cli.(provider.VictoriaMetricsProvider).Query(rule.PrometheusConfig.PromQL)
		if err != nil {
			logc.Error(ctx.Ctx, err.Error())
			return nil, nil
		}

		externalLabels = cli.(provider.VictoriaMetricsProvider).GetExternalLabels()
	default:
		logc.Errorf(ctx.Ctx, fmt.Sprintf("Unsupported metrics type, type: %s", datasourceType))
		return nil, nil
	}

	if resQuery == nil {
		return nil, nil
	}

	// 获取已缓存事件指纹
	fingerPrintMap := process.GetFingerPrint(ctx, rule.TenantId, rule.FaultCenterId, rule.RuleId)

	var curFingerprints []string
	for _, v := range resQuery {
		for _, ruleExpr := range rule.PrometheusConfig.Rules {
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
			event.Fingerprint = v.GetFingerprint()
			event.Metric = v.GetMetric()
			event.Metric["severity"] = ruleExpr.Severity
			for ek, ev := range externalLabels {
				event.Metric[ek] = ev
			}
			event.Severity = ruleExpr.Severity
			event.Annotations = tools.ParserVariables(rule.PrometheusConfig.Annotations, event.Metric)
			event.SearchQL = rule.PrometheusConfig.PromQL

			if process.EvalCondition(option) {
				// 如果告警条件满足需要将告警事件指纹加入，不满足时应当直接跳过
				curFingerprints = append(curFingerprints, event.Fingerprint)

				process.PushEventToFaultCenter(ctx, &event)
			} else {
				// 仅更新已经触发事件的指纹对应指标
				if _, exist := fingerPrintMap[event.Fingerprint]; exist {
					fingerPrintMetrics[event.Fingerprint] = event.Metric

					process.PushEventToFaultCenter(ctx, &event)
				}
			}
		}
	}

	return curFingerprints, fingerPrintMetrics
}

// Logs 包含 AliSLS、Loki、ElasticSearch 数据源
func logs(ctx *ctx.Context, datasourceId, datasourceType string, rule models.AlertRule) []string {
	var (
		queryRes       []provider.Logs
		count          int
		evalOptions    models.EvalCondition
		externalLabels map[string]interface{}
	)

	pools := ctx.Redis.ProviderPools()
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
