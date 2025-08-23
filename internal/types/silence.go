package types

import "watchAlert/internal/models"

// RequestSilenceCreate 请求创建静默规则
type RequestSilenceCreate struct {
	TenantId      string                `json:"tenantId"`
	Name          string                `json:"name"`
	Labels        []models.SilenceLabel `json:"labels" gorm:"labels;serializer:json"`
	StartsAt      int64                 `json:"startsAt"`
	UpdateBy      string                `json:"updateBy"`
	EndsAt        int64                 `json:"endsAt"`
	UpdateAt      int64                 `json:"updateAt"`
	FaultCenterId string                `json:"faultCenterId"`
	Comment       string                `json:"comment"`
	Status        int                   `json:"status"` // 0 未生效, 1 进行中, 2 已失效
}

// RequestSilenceUpdate 请求更新静默规则
type RequestSilenceUpdate struct {
	TenantId      string                `json:"tenantId"`
	Name          string                `json:"name"`
	ID            string                `json:"id"`
	Labels        []models.SilenceLabel `json:"labels" gorm:"labels;serializer:json"`
	StartsAt      int64                 `json:"startsAt"`
	UpdateBy      string                `json:"updateBy"`
	EndsAt        int64                 `json:"endsAt"`
	UpdateAt      int64                 `json:"updateAt"`
	FaultCenterId string                `json:"faultCenterId"`
	Comment       string                `json:"comment"`
	Status        int                   `json:"status"` // 0 未生效, 1 进行中, 2 已失效
}

// RequestSilenceQuery 请求查询静默规则
type RequestSilenceQuery struct {
	TenantId      string `json:"tenantId" form:"tenantId"`
	ID            string `json:"id" form:"id"`
	Query         string `json:"query" form:"query"`
	FaultCenterId string `json:"faultCenterId" form:"faultCenterId"`
	Status        int    `json:"status" form:"status"`
	models.Page
}

// ResponseSilenceList 返回静默规则列表
type ResponseSilenceList struct {
	List []models.AlertSilences `json:"list"`
	models.Page
}
