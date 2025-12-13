package provider

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	GetHTTPMethod  = "GET"
	PostHTTPMethod = "POST"
)

type HTTPer struct{}

// NewMetricsAwareHTTPer 创建支持指标的HTTP探测器
func NewMetricsAwareHTTPer() MetricsAwareProbe {
	return HTTPer{}
}

// HTTPResult 包含HTTP探测的核心结果
type HTTPResult struct {
	Address    string        `json:"address"`
	StatusCode int           `json:"status_code"`
	Latency    time.Duration `json:"latency"`
	IsTimeout  bool          `json:"is_timeout"`
	Error      string        `json:"error,omitempty"`
}

// PilotWithMetrics 执行HTTP探测并直接返回指标
func (h HTTPer) PilotWithMetrics(option EndpointOption, ruleInfo ProbeRuleInfo) []ProbeMetric {
	timestamp := time.Now().Unix()

	// 执行HTTP探测
	httpResult := h.executeHTTPProbe(option)

	// 创建基础标签
	baseLabels := map[string]any{
		"tenant_id": ruleInfo.TenantID,
		"rule_id":   ruleInfo.RuleID,
		"rule_name": ruleInfo.RuleName,
		"rule_type": ruleInfo.RuleType,
		"endpoint":  ruleInfo.Endpoint,
	}

	// 添加HTTP特定标签
	httpLabels := make(map[string]any)
	for k, v := range baseLabels {
		httpLabels[k] = v
	}

	// 创建指标列表
	metrics := []ProbeMetric{
		{
			Name:      "probe_http_response_time_ms",
			Help:      "HTTP response time in milliseconds",
			Type:      "gauge",
			Labels:    copyLabelsMap(httpLabels),
			Value:     float64(httpResult.Latency.Milliseconds()),
			Timestamp: timestamp,
		},
		{
			Name:      "probe_http_status_code",
			Help:      "HTTP response status code",
			Type:      "gauge",
			Labels:    copyLabelsMap(baseLabels),
			Value:     float64(httpResult.StatusCode),
			Timestamp: timestamp,
		},
		{
			Name:      "probe_http_success",
			Help:      "HTTP probe success (1 for reachable, 0 for unreachable)",
			Type:      "gauge",
			Labels:    copyLabelsMap(baseLabels),
			Value:     getHTTPSuccessValueFromResult(httpResult),
			Timestamp: timestamp,
		},
	}

	return metrics
}

// executeHTTPProbe 执行HTTP探测并收集核心指标数据
func (h HTTPer) executeHTTPProbe(option EndpointOption) HTTPResult {
	result := HTTPResult{
		Address: option.Endpoint,
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: time.Duration(option.Timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Proxy: http.ProxyFromEnvironment,
		},
	}

	// 创建请求
	var req *http.Request
	var err error

	switch strings.ToUpper(option.HTTP.Method) {
	case GetHTTPMethod:
		req, err = http.NewRequest(http.MethodGet, option.Endpoint, nil)
	case PostHTTPMethod:
		body := bytes.NewReader([]byte(option.HTTP.Body))
		req, err = http.NewRequest(http.MethodPost, option.Endpoint, body)
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
		}
	default:
		result.Error = fmt.Sprintf("unsupported HTTP method: %s", option.HTTP.Method)
		return result
	}

	if err != nil {
		result.Error = fmt.Sprintf("failed to create request: %v", err)
		return result
	}

	// 设置自定义头部
	for k, v := range option.HTTP.Header {
		req.Header.Set(k, v)
	}

	// 设置默认User-Agent
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "WatchAlert-Probe/1.0")
	}

	// 执行HTTP请求
	startTime := time.Now()
	resp, err := client.Do(req)
	result.Latency = time.Since(startTime)

	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	// 收集响应信息
	result.StatusCode = resp.StatusCode

	return result
}

// getHTTPSuccessValueFromResult 仅在网络请求失败（不通）时返回 0，
// 只要收到 HTTP 响应（无论状态码），即视为 success=1。
func getHTTPSuccessValueFromResult(result HTTPResult) float64 {
	if result.Error != "" {
		return 0.0 // 请求未完成，网络不通
	}
	return 1.0 // 成功收到 HTTP 响应，视为“通”
}
