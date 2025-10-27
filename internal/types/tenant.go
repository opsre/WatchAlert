package types

import "watchAlert/internal/models"

// RequestTenantCreate 请求创建租户
type RequestTenantCreate struct {
	Name             string `json:"name"`
	Manager          string `json:"manager"`
	Description      string `json:"description"`
	UserNumber       int64  `json:"userNumber"`
	RuleNumber       int64  `json:"ruleNumber"`
	DutyNumber       int64  `json:"dutyNumber"`
	NoticeNumber     int64  `json:"noticeNumber"`
	RemoveProtection *bool  `json:"removeProtection" gorm:"type:BOOL"`
	UserId           string `json:"userId" gorm:"-"`
	UpdateAt         int64  `json:"updateAt"`
}

func (requestTenantCreate *RequestTenantCreate) GetRemoveProtection() *bool {
	if requestTenantCreate.RemoveProtection == nil {
		isOk := false
		return &isOk
	}
	return requestTenantCreate.RemoveProtection
}

// RequestTenantUpdate 请求更新租户
type RequestTenantUpdate struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Manager          string `json:"manager"`
	Description      string `json:"description"`
	UserNumber       int64  `json:"userNumber"`
	RuleNumber       int64  `json:"ruleNumber"`
	DutyNumber       int64  `json:"dutyNumber"`
	NoticeNumber     int64  `json:"noticeNumber"`
	RemoveProtection *bool  `json:"removeProtection" gorm:"type:BOOL"`
	UserId           string `json:"userId" gorm:"-"`
	UpdateAt         int64  `json:"updateAt"`
}

func (requestTenantUpdate *RequestTenantUpdate) GetRemoveProtection() *bool {
	if requestTenantUpdate.RemoveProtection == nil {
		isOk := false
		return &isOk
	}
	return requestTenantUpdate.RemoveProtection
}

// RequestTenantQuery 请求搜索租户
type RequestTenantQuery struct {
	ID     string `json:"id" form:"id"`
	Name   string `json:"name" form:"name"`
	UserID string `json:"userId" form:"userId"`
}

// RequestTenantChangeUserRole 请求修改用户角色
type RequestTenantChangeUserRole struct {
	ID       string `json:"id"`
	UserID   string `json:"userId"`
	UserRole string `json:"userRole" `
}

// RequestTenantAddUsers 请求向租户内添加用户
type RequestTenantAddUsers struct {
	ID       string              `json:"id"`
	UserRole string              `json:"userRole" gorm:"-"` // 用于新增成员时统一的用户角色
	Users    []models.TenantUser `json:"users" gorm:"users;serializer:json"`
}
