package eval

import (
	"fmt"
	"sort"
	"watchAlert/alert/process"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/provider"
	"watchAlert/pkg/tools"

	"github.com/zeromicro/go-zero/core/logc"
)

// maxSeriesPerQuery 单次评估最多处理的时间序列数量，超出部分会被截断以避免影响性能。
const maxSeriesPerQuery = 1000

// compiledRule 是解析过操作符与期望值之后的规则，避免在每个时间序列上重复解析表达式。
type compiledRule struct {
	models.Rules
	operator      string
	expectedValue float64
}

// metricEvalContext 收敛一次 Prometheus 规则评估过程中保持不变的公共依赖，
// 避免在各步骤函数之间传递大量重复参数。
type metricEvalContext struct {
	ctx            *ctx.Context
	rule           models.AlertRule
	datasourceId   string
	prometheus     provider.PrometheusProvider
	externalLabels map[string]interface{}
	cachedEvents   map[string]*models.AlertCurEvent
}

// metrics Prometheus 数据源评估入口。
func metrics(c *ctx.Context, datasourceId, _ string, rule models.AlertRule) []string {
	prometheus, err := getPrometheusProvider(c, datasourceId)
	if err != nil {
		logc.Errorf(c.Ctx, "获取Prometheus数据源失败, 规则ID: %s, 规则名称: %s, 数据源ID: %s, 错误: %v", rule.RuleId, rule.RuleName, datasourceId, err)
		return nil
	}

	series, err := queryMetrics(c, prometheus, rule)
	if err != nil {
		logc.Errorf(c.Ctx, "Prometheus查询失败, 规则ID: %s, 规则名称: %s, 数据源ID: %s, PromQL: %s, 错误: %v", rule.RuleId, rule.RuleName, datasourceId, rule.PrometheusConfig.PromQL, err)
		return nil
	}
	if len(series) == 0 {
		return nil
	}

	rules := compileRules(c, rule)
	if len(rules) == 0 {
		return nil
	}

	evalCtx := &metricEvalContext{
		ctx:            c,
		rule:           rule,
		datasourceId:   datasourceId,
		prometheus:     prometheus,
		externalLabels: prometheus.GetExternalLabels(),
		cachedEvents:   getAlertEvents(c, rule),
	}

	return evalCtx.evalSeries(series, rules)
}

// getPrometheusProvider 获取并断言指定数据源的 Prometheus 客户端。
func getPrometheusProvider(c *ctx.Context, datasourceId string) (prometheus provider.PrometheusProvider, err error) {
	cli, err := c.Redis.ProviderPools().GetClient(datasourceId)
	if err != nil {
		return prometheus, fmt.Errorf("获取数据源客户端失败: %w", err)
	}

	prometheus, ok := cli.(provider.PrometheusProvider)
	if !ok {
		return prometheus, fmt.Errorf("数据源客户端类型错误, 类型: %T", cli)
	}

	return prometheus, nil
}

// queryMetrics 执行 PromQL 查询，并对超出上限的结果做截断保护。
func queryMetrics(c *ctx.Context, prometheus provider.PrometheusProvider, rule models.AlertRule) ([]provider.Metrics, error) {
	series, err := prometheus.Query(rule.PrometheusConfig.PromQL)
	if err != nil {
		return nil, err
	}

	if len(series) > maxSeriesPerQuery {
		logc.Errorf(c.Ctx, "Prometheus查询结果过多，可能影响性能，仅提取前 %d 个数据点，规则ID: %s, 规则名称: %s, 结果数量: %d", maxSeriesPerQuery, rule.RuleId, rule.RuleName, len(series))
		series = series[:maxSeriesPerQuery]
	}

	return series, nil
}

// evalSeries 依次评估每个时间序列，返回本轮评估中处于触发/抑制状态的指纹列表。
func (e *metricEvalContext) evalSeries(series []provider.Metrics, rules []compiledRule) []string {
	curFingerprints := make([]string, 0, len(series))
	activeFingerprints := make(map[string]struct{}, len(series))

	for _, s := range series {
		fingerprints := e.evalOneMetric(s, rules, activeFingerprints)
		curFingerprints = append(curFingerprints, fingerprints...)
	}

	return curFingerprints
}

// evalOneMetric 评估单个时间序列：命中的规则会推送/更新告警事件，未命中的规则用于更新恢复前的最后观测值。
func (e *metricEvalContext) evalOneMetric(s provider.Metrics, rules []compiledRule, activeFingerprints map[string]struct{}) []string {
	value := s.Value
	metricLabels := cloneLabels(s.GetMetric())
	matched, unmatched := classifyRules(rules, value)

	var newFingerprints []string
	if len(matched) > 0 {
		newFingerprints = e.triggerMatchedRules(matched, metricLabels, value, activeFingerprints)
	}

	// 更新未命中规则对应事件的恢复前最后观测值。
	e.updateRecoveryValues(metricLabels, unmatched, activeFingerprints, value)

	return newFingerprints
}

// triggerMatchedRules 处理命中的规则：生成/更新告警事件，并计算抑制状态后推送到故障中心。
func (e *metricEvalContext) triggerMatchedRules(matched []compiledRule, metricLabels map[string]interface{}, value float64, activeFingerprints map[string]struct{}) []string {
	newFingerprints := make([]string, 0, len(matched))
	highestPriority := getPriorityValue(matched[0].Severity)

	for _, rule := range matched {
		fingerprint := e.generateFingerprint(metricLabels, rule.Severity)
		event := e.buildAlertEvent(fingerprint, metricLabels, value, rule)
		desiredStatus := metricEventStatus(e.cachedEvents[fingerprint], rule, highestPriority)

		processCallbackPromQL(e.ctx, e.prometheus, e.rule.PrometheusConfig.CallbakPromQLs, event.Labels)
		if desiredStatus == "" {
			process.PushEventToFaultCenter(e.ctx, &event)
		} else {
			process.PushEventToFaultCenterWithStatus(e.ctx, &event, desiredStatus)
		}
		e.cachedEvents[fingerprint] = &event

		if _, exists := activeFingerprints[fingerprint]; !exists {
			activeFingerprints[fingerprint] = struct{}{}
			newFingerprints = append(newFingerprints, fingerprint)
		}
	}

	return newFingerprints
}

// buildAlertEvent 基于命中的规则构建一条告警事件。
func (e *metricEvalContext) buildAlertEvent(fingerprint string, metricLabels map[string]interface{}, value float64, rule compiledRule) models.AlertCurEvent {
	event := process.BuildEvent(e.rule, func() map[string]interface{} {
		labels := cloneLabels(metricLabels)
		labels["rule_name"] = e.rule.RuleName
		labels["fingerprint"] = fingerprint
		labels["severity"] = rule.Severity
		labels["value"] = value
		for key, labelValue := range e.externalLabels {
			labels[key] = labelValue
		}
		for key, labelValue := range e.rule.ExternalLabels {
			labels[key] = labelValue
		}
		if cachedEvent, exists := e.cachedEvents[fingerprint]; exists && cachedEvent.Labels["first_value"] != nil {
			labels["first_value"] = cachedEvent.Labels["first_value"]
		} else {
			labels["first_value"] = value
		}
		return labels
	})
	event.DatasourceId = e.datasourceId
	event.Fingerprint = fingerprint
	event.Severity = rule.Severity
	event.SearchQL = fmt.Sprintf("%s %s %v", e.rule.PrometheusConfig.PromQL, rule.operator, rule.expectedValue)
	event.ForDuration = e.rule.GetForDuration(rule.Severity)
	event.Annotations = tools.ParserVariables(e.rule.PrometheusConfig.Annotations, tools.ConvertStructToMap(event))
	event.Status = models.StatePreAlert
	return event
}

// generateFingerprint 基于规则 ID、规则名称、severity 以及指标标签生成事件指纹。
func (e *metricEvalContext) generateFingerprint(metricLabels map[string]interface{}, severity string) string {
	fingerprintLabels := cloneLabels(metricLabels)
	fingerprintLabels["rule_id"] = e.rule.RuleId
	fingerprintLabels["rule_name"] = e.rule.RuleName
	fingerprintLabels["severity"] = severity
	return provider.Metrics{Labels: fingerprintLabels}.GetFingerprint()
}

// updateRecoveryValues 为未命中规则对应的、尚未恢复的历史事件更新最后观测值。
func (e *metricEvalContext) updateRecoveryValues(metricLabels map[string]interface{}, unmatched []compiledRule, activeFingerprints map[string]struct{}, value float64) {
	updatedFingerprints := make(map[string]struct{}, len(unmatched))
	for _, rule := range unmatched {
		fingerprint := e.generateFingerprint(metricLabels, rule.Severity)
		if _, active := activeFingerprints[fingerprint]; active {
			continue
		}
		if _, updated := updatedFingerprints[fingerprint]; updated {
			continue
		}
		updatedFingerprints[fingerprint] = struct{}{}

		event, exists := e.cachedEvents[fingerprint]
		if !exists || event.IsRecovered || event.Status == models.StateRecovered {
			continue
		}
		if event.Labels == nil {
			event.Labels = make(map[string]interface{})
		}
		event.Labels["value"] = value
		process.PushEventToFaultCenter(e.ctx, event)
	}
}

// getAlertEvents 获取当前规则关联的所有缓存告警事件。
func getAlertEvents(c *ctx.Context, rule models.AlertRule) map[string]*models.AlertCurEvent {
	events, err := c.Redis.Alert().GetAllEvents(models.BuildAlertEventCacheKey(rule.TenantId, rule.FaultCenterId))
	if err != nil {
		logc.Errorf(c.Ctx, "获取指标告警缓存失败, 规则ID: %s, 规则名称: %s, 错误: %v", rule.RuleId, rule.RuleName, err)
		return make(map[string]*models.AlertCurEvent)
	}

	ruleEvents := make(map[string]*models.AlertCurEvent)
	for fingerprint, event := range events {
		if event.RuleId == rule.RuleId {
			ruleEvents[fingerprint] = event
		}
	}
	return ruleEvents
}

// compileRules 对规则排序并解析表达式一次，避免在每个时间序列上重复解析。
func compileRules(c *ctx.Context, rule models.AlertRule) []compiledRule {
	sortedRules := sortRulesByPriority(rule.PrometheusConfig.Rules)
	compiledRules := make([]compiledRule, 0, len(sortedRules))
	for _, r := range sortedRules {
		operator, expectedValue, err := process.ProcessRuleExpr(r.Expr)
		if err != nil {
			logc.Errorf(c.Ctx, "处理规则表达式失败, 规则ID: %s, 规则名称: %s, 表达式: %s, 错误: %v", rule.RuleId, rule.RuleName, r.Expr, err)
			continue
		}
		compiledRules = append(compiledRules, compiledRule{
			Rules:         r,
			operator:      operator,
			expectedValue: expectedValue,
		})
	}
	return compiledRules
}

// classifyRules 保持命中规则的优先级顺序，并返回未命中规则供恢复逻辑使用。
func classifyRules(rules []compiledRule, queryValue float64) (matched, unmatched []compiledRule) {
	matched = make([]compiledRule, 0, len(rules))
	unmatched = make([]compiledRule, 0, len(rules))
	for _, rule := range rules {
		if process.EvalCondition(models.EvalCondition{
			Operator:      rule.operator,
			QueryValue:    queryValue,
			ExpectedValue: rule.expectedValue,
		}) {
			matched = append(matched, rule)
			continue
		}
		unmatched = append(unmatched, rule)
	}
	return matched, unmatched
}

// cloneLabels 返回标签 map 的浅拷贝，避免调用方相互污染。
func cloneLabels(labels map[string]interface{}) map[string]interface{} {
	clone := make(map[string]interface{}, len(labels))
	for key, value := range labels {
		clone[key] = value
	}
	return clone
}

// processCallbackPromQL 对渲染后相同的回调查询语句只执行一次查询。
func processCallbackPromQL(c *ctx.Context, prometheus provider.PrometheusProvider, callbacks []models.CallbakPromQLs, labels map[string]interface{}) {
	results := make(map[string][]provider.Metrics, len(callbacks))
	for _, callback := range callbacks {
		query := tools.ParserVariables(callback.Value, map[string]interface{}{"labels": labels})
		result, queried := results[query]
		if !queried {
			var err error
			result, err = prometheus.Query(query)
			if err != nil {
				logc.Errorf(c.Ctx, "query callback promql error: %v, callback_key: %s, callback_promql: %s", err, callback.Key, callback.Value)
			}
			results[query] = result
		}
		if len(result) > 0 {
			labels[callback.Key] = result[0].GetValue()
		}
	}
}

// metricEventStatus 返回某个指标标签序列所需的显式状态转换。
// 优先级更低的命中规则会被最高优先级的命中规则抑制；
// 反之，一旦某条曾被抑制的规则成为最高命中规则，则恢复为告警状态。
func metricEventStatus(cachedEvent *models.AlertCurEvent, rule compiledRule, highestPriority int) models.AlertStatus {
	if getPriorityValue(rule.Severity) < highestPriority {
		return models.StateSuppression
	}
	if cachedEvent != nil && cachedEvent.Status == models.StateSuppression {
		return models.StateAlerting
	}
	return ""
}

// sortRulesByPriority 按优先级排序规则
func sortRulesByPriority(rules []models.Rules) []models.Rules {
	sortedRules := make([]models.Rules, len(rules))
	copy(sortedRules, rules)

	sort.SliceStable(sortedRules, func(i, j int) bool {
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
