package types

import "watchAlert/internal/models"

type RequestAuditLogQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	Query    string `json:"query" form:"query"`
	Scope    string `json:"scope" form:"scope"`
	models.Page
}

type ResponseAuditLog struct {
	List []models.AuditLog `json:"list"`
	models.Page
}
