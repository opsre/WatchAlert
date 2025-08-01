package types

import "watchAlert/internal/models"

type RequestDutyCalendarCreate struct {
	TenantId   string              `json:"tenantId"`
	DutyId     string              `json:"dutyId"`
	DutyPeriod int                 `json:"dutyPeriod"`
	Month      string              `json:"month"`
	UserGroup  [][]models.DutyUser `json:"userGroup"`
	DateType   string              `json:"dateType"`
	Status     string              `json:"status" `
}

type RequestDutyCalendarUpdate struct {
	TenantId string            `json:"tenantId"`
	DutyId   string            `json:"dutyId"`
	Time     string            `json:"time"`
	Users    []models.DutyUser `json:"users"`
	Status   string            `json:"status" `
}

type RequestDutyCalendarQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	DutyId   string `json:"dutyId" form:"dutyId"`
	Time     string `json:"time" form:"time"`
}
