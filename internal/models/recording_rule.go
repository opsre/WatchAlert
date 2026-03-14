package models

import (
	"fmt"
	"regexp"
)

// RecordingRule 记录规则模型
type RecordingRule struct {
	TenantId       string            `json:"tenantId" gorm:"column:tenant_id;index"`
	RuleId         string            `json:"ruleId" gorm:"column:rule_id;primaryKey"`
	DatasourceType string            `json:"datasourceType" gorm:"column:datasource_type"`
	DatasourceId   string            `json:"datasourceId" gorm:"column:datasource_id;serializer:json"`
	MetricName     string            `json:"metricName" gorm:"column:metric_name"`
	PromQL         string            `json:"promQL" gorm:"column:prom_ql"`
	Labels         map[string]string `json:"labels" gorm:"column:labels;serializer:json"`
	EvalInterval   int64             `json:"evalInterval" gorm:"column:eval_interval"`
	Enabled        *bool             `json:"enabled" gorm:"column:enabled"`
	CreateAt       int64             `json:"createAt" gorm:"column:create_at"`
	UpdateAt       int64             `json:"updateAt" gorm:"column:update_at"`
	CreateBy       string            `json:"createBy" gorm:"column:create_by"`
	UpdateBy       string            `json:"updateBy" gorm:"column:update_by"`
	RuleGroupId    int64             `json:"ruleGroupId" gorm:"column:rule_group_id;index"`
}

func (RecordingRule) TableName() string {
	return "w8t_recording_rules"
}

// GetEnabled 获取启用状态
func (r *RecordingRule) GetEnabled() *bool {
	if r.Enabled == nil {
		isOk := false
		return &isOk
	}
	return r.Enabled
}

// Validate 验证记录规则
func (r *RecordingRule) Validate() error {
	// 必填字段验证
	if r.MetricName == "" {
		return fmt.Errorf("指标名称不能为空")
	}
	if r.PromQL == "" {
		return fmt.Errorf("PromQL查询语句不能为空")
	}
	if len(r.DatasourceId) == 0 {
		return fmt.Errorf("数据源ID不能为空")
	}

	// 执行频率验证 (30秒 - 86400秒)
	if r.EvalInterval < 30 || r.EvalInterval > 86400 {
		return fmt.Errorf("执行频率必须在30秒到86400秒之间")
	}

	// 指标名称格式验证
	metricNamePattern := regexp.MustCompile(`^[a-zA-Z_:][a-zA-Z0-9_:]*$`)
	if !metricNamePattern.MatchString(r.MetricName) {
		return fmt.Errorf("指标名称格式不符合规范，必须匹配 [a-zA-Z_:][a-zA-Z0-9_:]*")
	}

	// 标签名称格式验证
	labelNamePattern := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	for labelName := range r.Labels {
		if !labelNamePattern.MatchString(labelName) {
			return fmt.Errorf("标签名称 '%s' 格式不符合规范，必须匹配 [a-zA-Z_][a-zA-Z0-9_]*", labelName)
		}
	}

	return nil
}
