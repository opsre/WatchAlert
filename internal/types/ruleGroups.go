package types

import "watchAlert/internal/models"

type RequestRuleGroupCreate struct {
	TenantId string `json:"tenantId"`
	ID       string `json:"id"`
	Name     string `json:"name"`
}

type RequestRuleGroupUpdate struct {
	TenantId string `json:"tenantId"`
	ID       string `json:"id"`
	Name     string `json:"name"`
}

type RequestRuleGroupQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	ID       string `json:"id" form:"id"`
	Name     string `json:"name" form:"name"`
	Query    string `json:"query" form:"query"`
	models.Page
}

type ResponseRuleGroupList struct {
	List []models.RuleGroups `json:"list"`
	models.Page
}
