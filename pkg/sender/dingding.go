package sender

import (
	"bytes"
	"errors"
	"fmt"
	"watchAlert/pkg/tools"
)

type (
	// DingDingSender 钉钉发送策略
	DingDingSender struct{}

	DingResponse struct {
		Code int    `json:"errcode"`
		Msg  string `json:"errmsg"`
	}
)

var DingTestContent = fmt.Sprintf(`{
	"msgtype": "text",
    "text": {
        "content": "%s"
    }
}`, RobotTestContent)

func NewDingSender() SendInter {
	return &DingDingSender{}
}

func (d *DingDingSender) Send(params SendParams) error {
	return d.post(params.Hook, params.Content)
}

func (d *DingDingSender) Test(params SendParams) error {
	return d.post(params.Hook, DingTestContent)
}

func (d *DingDingSender) post(hook, content string) error {
	cardContentByte := bytes.NewReader([]byte(content))
	res, err := tools.Post(nil, hook, cardContentByte, 10)
	if err != nil {
		return err
	}

	var response DingResponse
	if err := tools.ParseReaderBody(res.Body, &response); err != nil {
		return errors.New(fmt.Sprintf("Error unmarshalling Dingding response: %s", err.Error()))
	}
	if response.Code != 0 {
		return errors.New(response.Msg)
	}

	return nil
}
