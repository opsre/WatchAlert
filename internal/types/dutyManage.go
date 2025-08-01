package types

import "watchAlert/internal/models"

type RequestDutyManagementCreate struct {
	TenantId    string            `json:"tenantId"`
	Name        string            `json:"name"`
	Manager     models.DutyUser   `json:"manager" gorm:"manager;serializer:json"`
	Description string            `json:"description"`
	CurDutyUser []models.DutyUser `json:"curDutyUser" gorm:"curDutyUser;serializer:json"`
	CreateBy    string            `json:"create_by"`
	CreateAt    int64             `json:"create_at"`
}

type RequestDutyManagementUpdate struct {
	TenantId    string            `json:"tenantId"`
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Manager     models.DutyUser   `json:"manager" gorm:"manager;serializer:json"`
	Description string            `json:"description"`
	CurDutyUser []models.DutyUser `json:"curDutyUser" gorm:"curDutyUser;serializer:json"`
	CreateBy    string            `json:"create_by"`
	CreateAt    int64             `json:"create_at"`
}

type RequestDutyManagementQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	ID       string `json:"id" form:"id"`
	Name     string `json:"name" form:"name"`
}
