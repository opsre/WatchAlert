package models

const (
	FiringAlertCachePrefix  = "firing-alert-"
	PendingAlertCachePrefix = "pending-alert-"
)

type AlertCurEvent struct {
	TenantId               string                 `json:"tenantId"`
	RuleId                 string                 `json:"rule_id"`
	RuleName               string                 `json:"rule_name"`
	DatasourceType         string                 `json:"datasource_type"`
	DatasourceId           string                 `json:"datasource_id" gorm:"datasource_id"`
	Fingerprint            string                 `json:"fingerprint"`
	Severity               string                 `json:"severity"`
	Metric                 map[string]interface{} `json:"metric" gorm:"metric;serializer:json"`
	Labels                 map[string]string      `json:"labels" gorm:"labels;serializer:json"`
	EvalInterval           int64                  `json:"eval_interval"`
	ForDuration            int64                  `json:"for_duration"`
	NoticeId               string                 `json:"notice_id" gorm:"notice_id"` // 默认通知对象ID
	NoticeGroup            NoticeGroup            `json:"noticeGroup" gorm:"noticeGroup;serializer:json"`
	Annotations            string                 `json:"annotations" gorm:"-"`
	IsRecovered            bool                   `json:"is_recovered" gorm:"-"`
	FirstTriggerTime       int64                  `json:"first_trigger_time"` // 第一次触发时间
	FirstTriggerTimeFormat string                 `json:"first_trigger_time_format" gorm:"-"`
	RepeatNoticeInterval   int64                  `json:"repeat_notice_interval"`  // 重复通知间隔时间
	LastEvalTime           int64                  `json:"last_eval_time" gorm:"-"` // 上一次评估时间
	LastSendTime           int64                  `json:"last_send_time" gorm:"-"` // 上一次发送时间
	RecoverTime            int64                  `json:"recover_time" gorm:"-"`   // 恢复时间
	RecoverTimeFormat      string                 `json:"recover_time_format" gorm:"-"`
	DutyUser               string                 `json:"duty_user" gorm:"-"`
	EffectiveTime          EffectiveTime          `json:"effectiveTime" gorm:"effectiveTime;serializer:json"`
	RecoverNotify          *bool                  `json:"recoverNotify"`
	AlarmAggregation       *bool                  `json:"alarmAggregation"`

	ResponseTime  string `json:"response_time" gorm:"-"`
	TimeRemaining int64  `json:"time_remaining" gorm:"-"`
}

type AlertCurEventQuery struct {
	TenantId       string `json:"tenantId" form:"tenantId"`
	RuleId         string `json:"ruleId" form:"ruleId"`
	RuleName       string `json:"ruleName" form:"ruleName"`
	DatasourceType string `json:"datasourceType" form:"datasourceType"`
	DatasourceId   string `json:"datasourceId" form:"datasourceId"`
	Fingerprint    string `json:"fingerprint" form:"fingerprint"`
	Query          string `json:"query" form:"query"`
	Scope          int64  `json:"scope" form:"scope"`
	Severity       string `json:"severity" form:"severity"`
	Page
}

type CurEventResponse struct {
	List []AlertCurEvent `json:"list"`
	Page
}

func (ace *AlertCurEvent) GetFiringAlertCacheKey() string {
	return ace.TenantId + ":" + FiringAlertCachePrefix + ace.AlertCacheTailKey()
}

func (ace *AlertCurEvent) GetPendingAlertCacheKey() string {
	return ace.TenantId + ":" + PendingAlertCachePrefix + ace.AlertCacheTailKey()
}

func (ace *AlertCurEvent) AlertCacheTailKey() string {
	return ace.RuleId + "-" + ace.DatasourceId + "-" + ace.Fingerprint
}

func (ace *AlertCurEvent) GetRecoverNotify() *bool {
	if ace.RecoverNotify == nil {
		isOk := false
		return &isOk
	}
	return ace.RecoverNotify
}

func (ace *AlertCurEvent) GetAlarmAggregation() *bool {
	if ace.AlarmAggregation == nil {
		isOk := false
		return &isOk
	}
	return ace.AlarmAggregation
}
