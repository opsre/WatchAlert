package models

type AlertNotice struct {
	TenantId     string   `json:"tenantId"`
	Uuid         string   `json:"uuid"`
	Name         string   `json:"name"`
	DutyId       *string  `json:"dutyId"`
	NoticeType   string   `json:"noticeType"`
	NoticeTmplId string   `json:"noticeTmplId"`
	DefaultHook  string   `json:"hook" gorm:"column:hook"`
	DefaultSign  string   `json:"sign" gorm:"column:sign"`
	Routes       []Route  `json:"routes" gorm:"column:routes;serializer:json"`
	Email        Email    `json:"email" gorm:"email;serializer:json"`
	PhoneNumber  []string `json:"phoneNumber" gorm:"phoneNumber;serializer:json"`
	UpdateAt     int64    `json:"updateAt"`
	UpdateBy     string   `json:"updateBy"`
}

func (alertNotice *AlertNotice) GetDutyId() *string {
	if alertNotice.DutyId == nil {
		return new(string)
	}
	return alertNotice.DutyId
}

type Route struct {
	// 告警等级
	Severity string `json:"severity"`
	// WebHook
	Hook string `json:"hook"`
	// 签名
	Sign string `json:"sign"`
	// 收件人
	To []string `json:"to" gorm:"column:to;serializer:json"`
	// 抄送人
	CC []string `json:"cc" gorm:"column:cc;serializer:json"`
}

type Email struct {
	Subject string   `json:"subject"`
	To      []string `json:"to" gorm:"column:to;serializer:json"`
	CC      []string `json:"cc" gorm:"column:cc;serializer:json"`
}

type NoticeRecord struct {
	EventId  string `json:"eventId"`  // 事件ID
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

type ResponseNoticeRecords struct {
	List []NoticeRecord `json:"list"`
	Page
}
