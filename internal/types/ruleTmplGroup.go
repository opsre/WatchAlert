package types

import "watchAlert/internal/models"

type RequestRuleTemplateGroupCreate struct {
	Name   string `json:"name"`
	Number int    `json:"number"`
	Type   string `json:"type"`
}

type RequestRuleTemplateGroupUpdate struct {
	Name   string `json:"name"`
	Number int    `json:"number"`
	Type   string `json:"type"`
}

type RequestRuleTemplateGroupQuery struct {
	Name  string `json:"name" form:"name"`
	Type  string `json:"type" form:"type"`
	Query string `json:"query" form:"query"`
	models.Page
}

type ResponseRuleTemplateGroupList struct {
	List []models.RuleTemplateGroup `json:"list"`
	models.Page
}
