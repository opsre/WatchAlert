package sender

import (
	"bytes"
	"errors"
	"fmt"
	"watchAlert/pkg/tools"
)

type (
	// FeiShuSender 飞书发送策略
	FeiShuSender struct{}

	FeiShuResponse struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
)

func NewFeiShuSender() SendInter {
	return &FeiShuSender{}
}

func (f *FeiShuSender) Send(params SendParams) error {
	cardContentByte := bytes.NewReader([]byte(params.Content))
	res, err := tools.Post(nil, params.Hook, cardContentByte, 10)
	if err != nil {
		return err
	}

	var response FeiShuResponse
	if err := tools.ParseReaderBody(res.Body, &response); err != nil {
		return errors.New(fmt.Sprintf("Error unmarshalling Feishu response: %s", err.Error()))
	}
	if response.Code != 0 {
		return errors.New(response.Msg)
	}

	return nil
}
