package types

import "fmt"

type RequestAiChatContent struct {
	// 规则名称，用来分析告警时，更明确当前是一个什么规则
	RuleName string `json:"ruleName" form:"ruleName"`
	RuleId   string `json:"RuleId" form:"ruleId"`
	SearchQL string `json:"searchQL" form:"searchQL"`
	// 用户内容
	Content string `json:"content" form:"content"`
	// 重新分析，不调用缓存
	Deep string `json:"deep" form:"deep"`
}

func (a RequestAiChatContent) ValidateParams() error {
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
