package sender

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	msgStr, _ := json.Marshal(msg)
	msgByte := bytes.NewReader(msgStr)
	res, err := tools.Post(nil, params.Hook, msgByte, 10)
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
