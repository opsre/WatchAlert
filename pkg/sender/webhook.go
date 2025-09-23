package sender

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"watchAlert/pkg/tools"
)

type (
	// WebHookSender 自定义Hook发送策略
	WebHookSender struct{}
)

var WebhookTestContent = fmt.Sprintf(`{
  "text": "%s"
  }
}`, RobotTestContent)

func NewWebHookSender() SendInter {
	return &WebHookSender{}
}

func (w *WebHookSender) Send(params SendParams) error {
	return w.post(params.Hook, params.Content)
}

func (w *WebHookSender) Test(params SendParams) error {
	return w.post(params.Hook, WebhookTestContent)
}

func (w *WebHookSender) post(hook, content string) error {
	res, err := tools.Post(nil, hook, bytes.NewReader([]byte(content)), 10)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		bodyByte, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("读取 Body 失败, err: %s", err.Error())
		}
		return errors.New(string(bodyByte))
	}

	return nil
}
