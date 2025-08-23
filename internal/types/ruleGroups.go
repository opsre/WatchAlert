package types

import "watchAlert/internal/models"

type RequestRuleGroupCreate struct {
	TenantId    string `json:"tenantId"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Number      int    `json:"number"`
	Description string `json:"description"`
}

type RequestRuleGroupUpdate struct {
	TenantId    string `json:"tenantId"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Number      int    `json:"number"`
	Description string `json:"description"`
}

type RequestRuleGroupQuery struct {
	TenantId    string `json:"tenantId" form:"tenantId"`
	ID          string `json:"id" form:"id"`
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
	Query       string `json:"query" form:"query"`
	models.Page
}

type ResponseRuleGroupList struct {
	List []models.RuleGroups `json:"list"`
	models.Page
}
