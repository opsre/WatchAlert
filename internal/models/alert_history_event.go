package models

type AlertHisEvent struct {
	TenantId         string                 `json:"tenantId"`
	EventId          string                 `json:"eventId"`
	DatasourceId     string                 `json:"datasource_id" gorm:"datasource_id"`
	DatasourceType   string                 `json:"datasource_type"`
	Fingerprint      string                 `json:"fingerprint"`
	RuleId           string                 `json:"rule_id"`
	RuleName         string                 `json:"rule_name"`
	Severity         string                 `json:"severity"`
	Labels           map[string]interface{} `json:"labels" gorm:"labels;serializer:json"`
	EvalInterval     int64                  `json:"eval_interval"`
	Annotations      string                 `json:"annotations"`
	FirstTriggerTime int64                  `json:"first_trigger_time"` // 第一次触发时间
	LastEvalTime     int64                  `json:"last_eval_time"`     // 最近评估时间
	LastSendTime     int64                  `json:"last_send_time"`     // 最近发送时间
	RecoverTime      int64                  `json:"recover_time"`       // 恢复时间
	FaultCenterId    string                 `json:"faultCenterId"`
	ConfirmState     ConfirmState           `json:"confirmState" gorm:"metric;serializer:json"`
	AlarmDuration    int64                  `json:"alarmDuration"` // 告警持续时长
	SearchQL         string                 `json:"searchQL"`
}
