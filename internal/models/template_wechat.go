package models

type WeChatMsgTemplate struct {
	MsgType  string         `json:"msgtype"`
	MarkDown WeChatMarkDown `json:"markdown"`
}

type WeChatMarkDown struct {
	Content string `json:"content"`
}
