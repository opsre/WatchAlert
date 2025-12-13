package probe

import (
	"fmt"
	"strings"
	"watchAlert/internal/models"
)

// ValidateProbeRule 验证拨测规则
func ValidateProbeRule(rule models.ProbeRule) error {
	if rule.RuleId == "" {
		return fmt.Errorf("rule ID cannot be empty")
	}

	if rule.RuleName == "" {
		return fmt.Errorf("rule name cannot be empty")
	}

	if rule.ProbingEndpointConfig.Endpoint == "" {
		return fmt.Errorf("endpoint cannot be empty")
	}

	if rule.ProbingEndpointConfig.Strategy.EvalInterval <= 0 {
		return fmt.Errorf("eval interval must be greater than 0")
	}

	return nil
}

// ValidatePrometheusConfig 验证Prometheus配置
func ValidateWriteConfig(config MetricsWriterConfig) error {
	if config.Endpoint == "" {
		return fmt.Errorf("Write endpoint不能为空")
	}

	// 验证URL格式
	if !strings.HasPrefix(config.Endpoint, "http://") && !strings.HasPrefix(config.Endpoint, "https://") {
		return fmt.Errorf("Write endpoint必须以http://或https://开头")
	}

	return nil
}

// FormatProbeMetricName 格式化探测指标名称
func FormatProbeMetricName(name string) string {
	// 确保指标名称符合Prometheus命名规范
	return name
}

// MergeProbeLabels 合并探测标签
func MergeProbeLabels(base, additional map[string]any) map[string]any {
	result := make(map[string]any)

	// 复制基础标签
	for k, v := range base {
		result[k] = v
	}

	// 添加额外标签
	for k, v := range additional {
		result[k] = v
	}

	return result
}

// CopyLabels 复制标签映射
func CopyLabels(labels map[string]any) map[string]any {
	copied := make(map[string]any)
	for k, v := range labels {
		copied[k] = v
	}
	return copied
}

// BoolToFloat 将布尔值转换为浮点数
func BoolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}
