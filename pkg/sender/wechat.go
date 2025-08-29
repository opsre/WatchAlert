package sender

import (
	"bytes"
	"errors"
	"fmt"
	"watchAlert/pkg/tools"
)

type (
	// WeChatSender 企业微信发送策略
	WeChatSender struct{}

	WeChatResponse struct {
		Code int    `json:"errcode"`
		Msg  string `json:"errmsg"`
	}
)

func NewWeChatSender() SendInter {
	return &WeChatSender{}
}

func (w *WeChatSender) Send(params SendParams) error {
	cardContentByte := bytes.NewReader([]byte(params.Content))
	res, err := tools.Post(nil, params.Hook, cardContentByte, 10)
	if err != nil {
		return err
	}

	var response WeChatResponse
	if err := tools.ParseReaderBody(res.Body, &response); err != nil {
		return errors.New(fmt.Sprintf("Error unmarshalling Feishu response: %s", err.Error()))
	}
	if response.Code != 0 {
		return errors.New(response.Msg)
	}

	return nil
}
