package types

import "watchAlert/internal/models"

type RequestRecordingRuleGroupCreate struct {
	TenantId string `json:"tenantId"`
	Name     string `json:"name"`
}

type RequestRecordingRuleGroupUpdate struct {
	TenantId string `json:"tenantId"`
	ID       int64  `json:"id"`
	Name     string `json:"name"`
}

type RequestRecordingRuleGroupQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	ID       int64  `json:"id" form:"id"`
	Name     string `json:"name" form:"name"`
	Query    string `json:"query" form:"query"`
	models.Page
}

type ResponseRecordingRuleGroupList struct {
	List []models.RecordingRuleGroup `json:"list"`
	models.Page
}
