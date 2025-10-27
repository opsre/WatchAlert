package models

type Tenant struct {
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

func (t *Tenant) GetRemoveProtection() *bool {
	if t.RemoveProtection == nil {
		isOk := false
		return &isOk
	}
	return t.RemoveProtection
}

type TenantLinkedUsers struct {
	ID    string       `json:"id"`
	Users []TenantUser `json:"users" gorm:"users;serializer:json"`
}

type TenantUser struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	UserRole string `json:"userRole"`
}
