package types

import "watchAlert/internal/models"

// RequestProcessAlertEvent 请求处理告警事件
type RequestProcessAlertEvent struct {
	TenantId      string   `json:"tenantId"`
	FaultCenterId string   `json:"faultCenterId"`
	Fingerprints  []string `json:"fingerprints"`
	Time          int64    `json:"time"`
	Username      string   `json:"username"`
}

// RequestAlertCurEventQuery 请求活跃告警事件
type RequestAlertCurEventQuery struct {
	TenantId       string `json:"tenantId" form:"tenantId"`
	RuleId         string `json:"ruleId" form:"ruleId"`
	RuleName       string `json:"ruleName" form:"ruleName"`
	DatasourceType string `json:"datasourceType" form:"datasourceType"`
	DatasourceId   string `json:"datasourceId" form:"datasourceId"`
	Fingerprint    string `json:"fingerprint" form:"fingerprint"`
	Query          string `json:"query" form:"query"`
	Scope          int64  `json:"scope" form:"scope"`
	Severity       string `json:"severity" form:"severity"`
	FaultCenterId  string `json:"faultCenterId" form:"faultCenterId"`
	Status         string `json:"status" form:"status"`
	SortOrder      string `json:"sortOrder" form:"sortOrder"`
	models.Page
}

// ResponseAlertCurEventList 返回活跃告警列表
type ResponseAlertCurEventList struct {
	List []models.AlertCurEvent `json:"list"`
	models.Page
}

// RequestAlertHisEventQuery 请求查询历史事件
type RequestAlertHisEventQuery struct {
	TenantId       string `json:"tenantId" form:"tenantId"`
	DatasourceId   string `json:"datasourceId" form:"datasourceId"`
	DatasourceType string `json:"datasourceType" form:"datasourceType"`
	Fingerprint    string `json:"fingerprint" form:"fingerprint"`
	Severity       string `json:"severity" form:"severity"`
	RuleId         string `json:"ruleId" form:"ruleId"`
	RuleName       string `json:"ruleName" form:"ruleName"`
	StartAt        int64  `json:"startAt" form:"startAt"`
	EndAt          int64  `json:"endAt" form:"endAt"`
	Query          string `json:"query" form:"query"`
	FaultCenterId  string `json:"faultCenterId" form:"faultCenterId"`
	SortOrder      string `json:"sortOrder" form:"sortOrder"`
	models.Page
}

// ResponseHistoryEventList 返回历史事件列表
type ResponseHistoryEventList struct {
	List []models.AlertHisEvent `json:"list"`
	models.Page
}

// RequestAddEventComment 添加评论
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

// RequestDeleteEventComment 删除评论
type RequestDeleteEventComment struct {
	// 租户
	TenantId string `json:"tenantId"`
	// 评论 ID
	CommentId string `json:"commentId"`
}

// RequestListEventComments 获取评论
type RequestListEventComments struct {
	// 租户
	TenantId string `json:"tenantId" form:"tenantId"`
	// 告警指纹
	Fingerprint string `json:"fingerprint" form:"fingerprint"`
}
