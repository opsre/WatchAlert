package models

type AiContentRecord struct {
	RuleId string `json:"RuleId" form:"ruleId"`
	// Ai 分析后的内容
	Content string `json:"content" form:"content"`
}

func (a AiContentRecord) TableName() string {
	return "w8t_ai_content_record"
}
