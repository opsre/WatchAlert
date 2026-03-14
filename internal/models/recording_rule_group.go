package models

type RecordingRuleGroup struct {
	ID       int64  `json:"id" gorm:"autoIncrement"`
	TenantId string `json:"tenantId" gorm:"index"`
	Name     string `json:"name" gorm:"column:name"`
}

func (RecordingRuleGroup) TableName() string {
	return "w8t_recording_rule_groups"
}
