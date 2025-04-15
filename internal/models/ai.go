package models

import "fmt"

type AiParams struct {
	// 规则名称，用来分析告警时，更明确当前是一个什么规则
	RuleName string `json:"ruleName" form:"ruleName"`
	RuleId   string `json:"RuleId" form:"ruleId"`
	SearchQL string `json:"searchQL" form:"searchQL"`
	// 用户内容
	Content string `json:"content" form:"content"`
	// 重新分析，不调用缓存
	Deep string `json:"deep" form:"deep"`
}

func (a AiParams) ValidateParams() error {
	if a.Content == "" {
		return fmt.Errorf("告警事件详情不可为空")
	}
	if a.RuleName == "" {
		return fmt.Errorf("规则名称不可为空")
	}
	if a.RuleId == "" {
		return fmt.Errorf("规则 ID 不可为空")
	}

	return nil
}

type AiContentRecord struct {
	RuleId string `json:"RuleId" form:"ruleId"`
	// Ai 分析后的内容
	Content string `json:"content" form:"content"`
}

func (a AiContentRecord) TableName() string {
	return "w8t_ai_content_record"
}
