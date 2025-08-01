package models

type AlertSubscribe struct {
	SId               string   `json:"sId"`                                                // 订阅规则ID
	STenantId         string   `json:"sTenantId"`                                          // 订阅规则的租户
	SUserId           string   `json:"sUserId"`                                            // 订阅的用户 ID
	SUserEmail        string   `json:"sUserEmail"`                                         // 订阅的用户邮箱
	SRuleId           string   `json:"sRuleId"`                                            // 订阅的规则 ID
	SRuleName         string   `json:"sRuleName"`                                          // 订阅的规则名称
	SRuleType         string   `json:"sRuleType"`                                          // 订阅的规则类型
	SRuleSeverity     []string `json:"sRuleSeverity" gorm:"sRuleSeverity;serializer:json"` // 订阅的告警等级
	SNoticeSubject    string   `json:"sNoticeSubject"`                                     // 发布订阅消息的 Title
	SNoticeTemplateId string   `json:"sNoticeTemplateId"`                                  // 发送订阅消息的通知模版 ID
	SFilter           []string `json:"sFilter" gorm:"sFilter;serializer:json"`             // 过滤
	SCreateAt         int64    `json:"sCreateAt"`
}
