package models

import (
	"gorm.io/gorm"
)

type AlertNotice struct {
	TenantId     string   `json:"tenantId"`
	Uuid         string   `json:"uuid"`
	Name         string   `json:"name"`
	DutyId       string   `json:"dutyId"`
	NoticeType   string   `json:"noticeType"`
	NoticeTmplId string   `json:"noticeTmplId"`
	Hook         string   `json:"hook"`
	Email        Email    `json:"email" gorm:"email;serializer:json"`
	// 签名
	Sign         string   `json:"sign" gorm:"sign"`
	PhoneNumber  []string `json:"phoneNumber" gorm:"phoneNumber;serializer:json"`
}

type Email struct {
	Subject string   `json:"subject"`
	To      []string `json:"to" gorm:"to;serializer:json"`
	CC      []string `json:"cc" gorm:"cc;serializer:json"`
}

type AlertRecord struct {
	gorm.Model
	AlertName   string `json:"alertName"`
	Description string `json:"description"`
	Metric      string `json:"metric"`
	Severity    string `json:"severity"`
	Status      string `json:"status"`
}

type NoticeTemplateExample struct {
	Id                   string `json:"id"`
	Name                 string `json:"name"`
	NoticeType           string `json:"noticeType"`
	Description          string `json:"description"`
	Template             string `json:"template"`
	TemplateFiring       string `json:"templateFiring"`
	TemplateRecover      string `json:"templateRecover"`
	EnableFeiShuJsonCard *bool  `json:"enableFeiShuJsonCard"`
}

type NoticeQuery struct {
	TenantId     string `json:"tenantId" form:"tenantId"`
	Uuid         string `json:"uuid" form:"uuid"`
	Name         string `json:"name" form:"name"`
	NoticeTmplId string `json:"noticeTmplId" form:"noticeTmplId"`
	Status       string `json:"status" form:"status"`
	Severity     string `json:"severity" form:"severity"`
	Query        string `json:"query" form:"query"`
}

type NoticeTemplateExampleQuery struct {
	Id         string `json:"id" form:"id"`
	Name       string `json:"name" form:"name"`
	NoticeType string `json:"noticeType" form:"noticeType"`
	Query      string `json:"query" form:"query"`
}

type NoticeRecord struct {
	Date     string `json:"date"`     // 记录日期
	CreateAt int64  `json:"createAt"` // 记录时间
	TenantId string `json:"tenantId"` // 租户
	RuleName string `json:"ruleName"` // 规则名称
	NType    string `json:"nType"`    // 通知类型
	NObj     string `json:"nObj"`     // 通知对象
	Severity string `json:"severity"` // 告警等级
	Status   int    `json:"status"`   // 通知状态 0 成功 1 失败
	AlarmMsg string `json:"alarmMsg"` // 告警信息
	ErrMsg   string `json:"errMsg"`   // 错误信息
}

type CountRecord struct {
	Date     string `json:"date"`     // 记录日期
	TenantId string `json:"tenantId"` // 租户
	Severity string `json:"severity"` // 告警等级
}
