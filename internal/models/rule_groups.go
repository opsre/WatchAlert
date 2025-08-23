package models

type RuleGroups struct {
	TenantId    string `json:"tenantId"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Number      int    `json:"number"`
	Description string `json:"description"`
}
