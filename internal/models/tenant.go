package models

type Tenant struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	CreateAt         int64  `json:"createAt"`
	CreateBy         string `json:"createBy"`
	Manager          string `json:"manager"`
	Description      string `json:"description"`
	UserNumber       int64  `json:"userNumber"`
	RuleNumber       int64  `json:"ruleNumber"`
	DutyNumber       int64  `json:"dutyNumber"`
	NoticeNumber     int64  `json:"noticeNumber"`
	RemoveProtection *bool  `json:"removeProtection" gorm:"type:BOOL"`
	UserId           string `json:"userId" gorm:"-"`
}

func (t *Tenant) GetRemoveProtection() *bool {
	if t.RemoveProtection == nil {
		isOk := false
		return &isOk
	}
	return t.RemoveProtection
}

type TenantQuery struct {
	ID     string `json:"id" form:"id"`
	Name   string `json:"name" form:"name"`
	UserID string `json:"userId" form:"userId"`
}

type TenantLinkedUsers struct {
	ID       string       `json:"id"`
	UserRole string       `json:"userRole" gorm:"-"` // 用于新增成员时统一的用户角色
	Users    []TenantUser `json:"users" gorm:"users;serializer:json"`
}

type TenantUser struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	UserRole string `json:"userRole"`
}

type GetTenantLinkedUserInfo struct {
	ID     string `json:"id" form:"id"`
	UserID string `json:"userId" form:"userId"`
}

type ChangeTenantUserRole struct {
	ID       string `json:"id"`
	UserID   string `json:"userId"`
	UserRole string `json:"userRole" `
}
