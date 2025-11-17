package types

import "watchAlert/internal/models"

type RequestNoticeCreate struct {
	TenantId     string         `json:"tenantId"`
	Name         string         `json:"name"`
	DutyId       *string        `json:"dutyId"`
	NoticeType   string         `json:"noticeType"`
	NoticeTmplId string         `json:"noticeTmplId"`
	DefaultHook  string         `json:"hook" gorm:"column:hook"`
	DefaultSign  string         `json:"sign" gorm:"column:sign"`
	Routes       []models.Route `json:"routes" gorm:"column:routes;serializer:json"`
	Email        models.Email   `json:"email" gorm:"email;serializer:json"`
	PhoneNumber  []string       `json:"phoneNumber" gorm:"phoneNumber;serializer:json"`
	UpdateBy     string         `json:"updateBy"`
}

type RequestNoticeUpdate struct {
	TenantId     string         `json:"tenantId"`
	Uuid         string         `json:"uuid"`
	Name         string         `json:"name"`
	DutyId       *string        `json:"dutyId"`
	NoticeType   string         `json:"noticeType"`
	NoticeTmplId string         `json:"noticeTmplId"`
	DefaultHook  string         `json:"hook" gorm:"column:hook"`
	DefaultSign  string         `json:"sign" gorm:"column:sign"`
	Routes       []models.Route `json:"routes" gorm:"column:routes;serializer:json"`
	Email        models.Email   `json:"email" gorm:"email;serializer:json"`
	PhoneNumber  []string       `json:"phoneNumber" gorm:"phoneNumber;serializer:json"`
	UpdateBy     string         `json:"updateBy"`
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
	NoticeType  string         `json:"noticeType"`
	DefaultHook string         `json:"hook"`
	DefaultSign string         `json:"sign"`
	Routes      []models.Route `json:"routes"`
	Email       models.Email   `json:"email"`
}
