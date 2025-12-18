package types

import "watchAlert/internal/models"

type RequestNoticeCreate struct {
	TenantId string         `json:"tenantId"`
	Name     string         `json:"name"`
	DutyId   *string        `json:"dutyId"`
	Routes   []models.Route `json:"routes" gorm:"column:routes;serializer:json"`
	UpdateBy string         `json:"updateBy"`
}

type RequestNoticeUpdate struct {
	TenantId string         `json:"tenantId"`
	Uuid     string         `json:"uuid"`
	Name     string         `json:"name"`
	DutyId   *string        `json:"dutyId"`
	Routes   []models.Route `json:"routes" gorm:"column:routes;serializer:json"`
	UpdateBy string         `json:"updateBy"`
}

func (requestNoticeUpdate *RequestNoticeUpdate) GetDutyId() *string {
	if requestNoticeUpdate.DutyId == nil {
		return new(string)
	}
	return requestNoticeUpdate.DutyId
}

type RequestNoticeQuery struct {
	TenantId     string `json:"tenantId" form:"tenantId"`
	EventId      string `json:"eventId" form:"eventId"`
	Uuid         string `json:"uuid" form:"uuid"`
	Name         string `json:"name" form:"name"`
	NoticeTmplId string `json:"noticeTmplId" form:"noticeTmplId"`
	Status       string `json:"status" form:"status"`
	Severity     string `json:"severity" form:"severity"`
	Query        string `json:"query" form:"query"`
	models.Page
}

type RequestNoticeTest struct {
	NoticeType string       `json:"noticeType"`
	Hook       string       `json:"hook"`
	Sign       string       `json:"sign"`
	Email      models.Email `json:"email"`
}
