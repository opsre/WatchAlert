package types

import "watchAlert/internal/models"

type RequestDutyManagementCreate struct {
	TenantId    string            `json:"tenantId"`
	Name        string            `json:"name"`
	Manager     models.DutyUser   `json:"manager" gorm:"manager;serializer:json"`
	Description string            `json:"description"`
	CurDutyUser []models.DutyUser `json:"curDutyUser" gorm:"curDutyUser;serializer:json"`
	UpdateBy    string            `json:"updateBy"`
	UpdateAt    int64             `json:"updateAt"`
}

type RequestDutyManagementUpdate struct {
	TenantId    string            `json:"tenantId"`
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Manager     models.DutyUser   `json:"manager" gorm:"manager;serializer:json"`
	Description string            `json:"description"`
	CurDutyUser []models.DutyUser `json:"curDutyUser" gorm:"curDutyUser;serializer:json"`
	UpdateBy    string            `json:"updateBy"`
	UpdateAt    int64             `json:"updateAt"`
}

type RequestDutyManagementQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	ID       string `json:"id" form:"id"`
	Name     string `json:"name" form:"name"`
}
