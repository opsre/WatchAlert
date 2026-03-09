package models

type NoticeTemplateExample struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	NoticeType           string `json:"noticeType"`
	Description          string `json:"description"`
	Template             string `json:"template"`
	TemplateFiring       string `json:"templateFiring"`
	TemplateRecover      string `json:"templateRecover"`
	EnableFeiShuJsonCard *bool  `json:"enableFeiShuJsonCard"`
	UpdateAt             int64  `json:"updateAt"`
	UpdateBy             string `json:"updateBy"`
}
