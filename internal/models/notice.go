package models

type AlertNotice struct {
	TenantId string  `json:"tenantId"`
	Uuid     string  `json:"uuid"`
	Name     string  `json:"name"`
	DutyId   *string `json:"dutyId"`
	Routes   []Route `json:"routes" gorm:"column:routes;serializer:json"`
	UpdateAt int64   `json:"updateAt"`
	UpdateBy string  `json:"updateBy"`
}

func (alertNotice *AlertNotice) GetDutyId() *string {
	if alertNotice.DutyId == nil {
		return new(string)
	}
	return alertNotice.DutyId
}

type Route struct {
	// 通知类型
	NoticeType string `json:"noticeType"`
	// 通知模版 ID
	NoticeTmplId string `json:"noticeTmplId"`
	// 告警等级
	Severitys []string `json:"severitys"`
	// WebHook
	Hook string `json:"hook"`
	// 签名
	Sign string `json:"sign"`
	// 邮件主题
	Subject string `json:"subject"`
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
