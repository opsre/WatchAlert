package provider

import (
	"time"
	"watchAlert/pkg/tools"
)

type Ssler struct{}

// NewMetricsAwareSSLer 创建支持指标的SSL探测器
func NewMetricsAwareSSLer() MetricsAwareProbe {
	return Ssler{}
}

// PilotWithMetrics 执行SSL探测并直接返回指标
func (p Ssler) PilotWithMetrics(option EndpointOption, ruleInfo ProbeRuleInfo) []Metrics {
	timestamp := time.Now().Unix()
	startTime := time.Now()

	// 创建基础标签
	baseLabels := map[string]any{
		"tenant_id":  ruleInfo.TenantID,
		"probe_id":   ruleInfo.RuleID,
		"probe_name": ruleInfo.RuleName,
		"probe_type": ruleInfo.RuleType,
		"endpoint":   ruleInfo.Endpoint,
	}

	// 发起 HTTPS 请求
	resp, err := tools.Get(nil, "https://"+option.Endpoint, option.Timeout)
	if err != nil {
		// 返回失败指标
		return p.createFailureMetrics(baseLabels, timestamp, time.Since(startTime))
	}
	defer resp.Body.Close()

	// 证书为空, 跳过检测
	if resp.TLS == nil || len(resp.TLS.PeerCertificates) == 0 {
		// 返回失败指标
		return p.createFailureMetrics(baseLabels, timestamp, time.Since(startTime))
	}

	// 获取证书信息
	cert := resp.TLS.PeerCertificates[0]
	notAfter := cert.NotAfter // 证书过期时间
	currentTime := time.Now()

	// 计算剩余有效期（单位：天）
	timeRemaining := int64(notAfter.Sub(currentTime).Hours() / 24)
	responseTime := time.Since(startTime)

	// 创建SSL指标
	metrics := []Metrics{
		{
			Name:   "probe_ssl_certificate_expiry_days",
			Help:   "SSL certificate expiry time in days",
			Labels: copyLabelsMap(baseLabels),
			Value:  float64(timeRemaining),
		},
		{
			Name:   "probe_ssl_response_time_ms",
			Help:   "SSL handshake response time in milliseconds",
			Labels: copyLabelsMap(baseLabels),
			Value:  float64(responseTime.Milliseconds()),
		},
		{
			Name:   "probe_ssl_certificate_valid",
			Help:   "SSL certificate validity (1 for valid, 0 for invalid/expired)",
			Labels: copyLabelsMap(baseLabels),
			Value:  boolToFloat(timeRemaining > 0),
		},
	}

	return metrics
}

// createFailureMetrics 创建失败时的指标
func (p Ssler) createFailureMetrics(baseLabels map[string]any, timestamp int64, responseTime time.Duration) []Metrics {
	return []Metrics{
		{
			Name:   "probe_ssl_certificate_expiry_days",
			Help:   "SSL certificate expiry time in days",
			Labels: copyLabelsMap(baseLabels),
			Value:  -1, // 表示无法获取
		},
		{
			Name:   "probe_ssl_response_time_ms",
			Help:   "SSL handshake response time in milliseconds",
			Labels: copyLabelsMap(baseLabels),
			Value:  float64(responseTime.Milliseconds()),
		},
		{
			Name:   "probe_ssl_certificate_valid",
			Help:   "SSL certificate validity (1 for valid, 0 for invalid/expired)",
			Labels: copyLabelsMap(baseLabels),
			Value:  0.0, // 失败
		},
	}
}
