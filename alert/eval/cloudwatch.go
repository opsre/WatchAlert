package eval

import (
	"fmt"
	"time"
	"watchAlert/alert/process"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/aws/cloudwatch"
	"watchAlert/pkg/provider"
	"watchAlert/pkg/tools"

	"github.com/zeromicro/go-zero/core/logc"
)

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
		query := cloudwatch.CloudWatchQuery{
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
