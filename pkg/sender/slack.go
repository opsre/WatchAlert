package sender

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

type (
	// SlackSender Slack 发送策略
	SlackSender struct{}
)

func NewSlackSender() SendInter {
	return &SlackSender{}
}

func (f *SlackSender) Send(params SendParams) error {
	msg := params.GetSendMsg()
	return f.post(params.Hook, tools.JsonMarshalToString(msg))
}

func (f *SlackSender) Test(params SendParams) error {
	msg := models.SlackMsgTemplate{
		Text: RobotTestContent,
	}
	return f.post(params.Hook, tools.JsonMarshalToString(msg))
}

func (f *SlackSender) post(hook, content string) error {
	res, err := tools.Post(nil, hook, bytes.NewReader([]byte(content)), 10)
	if err != nil {
		return err
	}

	bodyByte, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("Error unmarshalling Slack response: %s", err.Error()))
	}

	if string(bodyByte) != "ok" {
		return errors.New(string(bodyByte))
	}

	return nil
}
