package models

// Comment 告警事件评论表
type Comment struct {
	// 自增 ID
	ID uint `json:"id" gorm:"primaryKey"`
	// 租户 ID
	TenantId string `json:"tenantId"`
	// 评论 ID
	CommentId string `json:"commentId"`
	// 事件指纹
	Fingerprint string `json:"fingerprint"`
	// 用户名
	Username string `json:"username"`
	// 用户 ID
	UserId string `json:"userId"`
	// 时间
	Time int64 `json:"time"`
	// 内容
	Content string `json:"content"`
}

type RequestAddEventComment struct {
	// 租户
	TenantId string `json:"tenantId"`
	// 故障中心
	FaultCenterId string `json:"faultCenterId"`
	// 告警指纹
	Fingerprint string `json:"fingerprint"`
	// 用户名
	Username string `json:"username"`
	// 用户 ID
	UserId string `json:"userId"`
	// 评论内容
	Content string `json:"content"`
}

type RequestDeleteEventComment struct {
	// 租户
	TenantId string `json:"tenantId"`
	// 评论 ID
	CommentId string `json:"commentId"`
}

type RequestListEventComments struct {
	// 租户
	TenantId string `json:"tenantId" form:"tenantId"`
	// 告警指纹
	Fingerprint string `json:"fingerprint" form:"fingerprint"`
}
