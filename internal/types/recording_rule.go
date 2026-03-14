package types

import "watchAlert/internal/models"

type RequestRecordingRuleCreate struct {
	TenantId       string            `json:"tenantId"`
	RuleId         string            `json:"ruleId"`
	DatasourceType string            `json:"datasourceType"`
	DatasourceId   string            `json:"datasourceId"`
	MetricName     string            `json:"metricName"`
	PromQL         string            `json:"promQL"`
	Labels         map[string]string `json:"labels"`
	EvalInterval   int64             `json:"evalInterval"`
	UpdateBy       string            `json:"updateBy"`
	Enabled        *bool             `json:"enabled"`
	RuleGroupId    int64             `json:"ruleGroupId"`
}

func (requestRecordingRuleCreate *RequestRecordingRuleCreate) GetEnabled() *bool {
	if requestRecordingRuleCreate.Enabled == nil {
		isOk := false
		return &isOk
	}
	return requestRecordingRuleCreate.Enabled
}

type RequestRecordingRuleUpdate struct {
	TenantId       string            `json:"tenantId"`
	RuleId         string            `json:"ruleId"`
	DatasourceType string            `json:"datasourceType"`
	DatasourceId   string            `json:"datasourceId"`
	MetricName     string            `json:"metricName"`
	PromQL         string            `json:"promQL"`
	Labels         map[string]string `json:"labels"`
	EvalInterval   int64             `json:"evalInterval"`
	UpdateBy       string            `json:"updateBy"`
	Enabled        *bool             `json:"enabled"`
	RuleGroupId    int64             `json:"ruleGroupId"`
}

func (requestRecordingRuleUpdate *RequestRecordingRuleUpdate) GetEnabled() *bool {
	if requestRecordingRuleUpdate.Enabled == nil {
		isOk := false
		return &isOk
	}
	return requestRecordingRuleUpdate.Enabled
}

type RequestRecordingRuleQuery struct {
	TenantId       string `json:"tenantId" form:"tenantId"`
	RuleId         string `json:"ruleId" form:"ruleId"`
	DatasourceType string `json:"datasourceType" form:"datasourceType"`
	Enabled        string `json:"enabled" form:"enabled"`
	Query          string `json:"query" form:"query"`
	Status         string `json:"status" form:"status"`
	RuleGroupId    int64  `json:"ruleGroupId" form:"ruleGroupId"`
	models.Page
}

type ResponseRecordingRuleList struct {
	List []models.RecordingRule `json:"list"`
	models.Page
}

type RequestRecordingRuleChangeStatus struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	RuleId   string `json:"ruleId" form:"ruleId"`
	Enabled  *bool  `json:"enabled" form:"enabled"`
}

func (r *RequestRecordingRuleChangeStatus) GetEnabled() *bool {
	if r.Enabled == nil {
		isOk := false
		return &isOk
	}
	return r.Enabled
}
