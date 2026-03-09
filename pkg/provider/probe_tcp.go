package provider

import (
	"net"
	"time"
)

type Tcper struct{}

// NewMetricsAwareTcper 创建支持指标的TCP探测器
func NewMetricsAwareTcper() MetricsAwareProbe {
	return Tcper{}
}

// PilotWithMetrics 执行TCP探测并直接返回指标
func (p Tcper) PilotWithMetrics(option EndpointOption, ruleInfo ProbeRuleInfo) []ProbeMetric {
	timestamp := time.Now().Unix()
	startTime := time.Now()

	// 尝试拨测指定地址和端口
	conn, err := net.DialTimeout("tcp", option.Endpoint, time.Duration(option.Timeout)*time.Second)
	responseTime := time.Since(startTime)

	// 创建基础标签
	baseLabels := map[string]any{
		"tenant_id":  ruleInfo.TenantID,
		"probe_id":   ruleInfo.RuleID,
		"probe_name": ruleInfo.RuleName,
		"probe_type": ruleInfo.RuleType,
		"endpoint":   ruleInfo.Endpoint,
	}

	// 确定成功状态
	isSuccessful := err == nil
	if isSuccessful && conn != nil {
		conn.Close()
	}

	// 创建TCP指标
	metrics := []ProbeMetric{
		{
			Name:      "probe_tcp_success",
			Help:      "TCP probe success (1 for success, 0 for failure)",
			Type:      "gauge",
			Labels:    copyLabelsMap(baseLabels),
			Value:     BoolToFloat(isSuccessful),
			Timestamp: timestamp,
		},
		{
			Name:      "probe_tcp_response_time_ms",
			Help:      "TCP connection response time in milliseconds",
			Type:      "gauge",
			Labels:    copyLabelsMap(baseLabels),
			Value:     float64(responseTime.Milliseconds()),
			Timestamp: timestamp,
		},
	}

	return metrics
}
