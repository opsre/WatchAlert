package aliyun

import (
	"errors"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dyvmsapi "github.com/alibabacloud-go/dyvmsapi-intl-20211015/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/zeromicro/go-zero/core/logc"
	"go.uber.org/multierr"
	"watchAlert/internal/ctx"
)

type PhoneCall struct {
	Endpoint        string `json:"endpoint,omitempty"`
	AccessKeyId     string `json:"accessKeyId,omitempty"`
	AccessKeySecret string `json:"accessKeySecret,omitempty"`
	TtsCode         string `json:"ttsCode,omitempty"`
	Client          *dyvmsapi.Client
}

func (p *PhoneCall) CreateClient() error {
	config := &openapi.Config{
		AccessKeyId:     tea.String(p.AccessKeyId),
		AccessKeySecret: tea.String(p.AccessKeySecret),
		Endpoint:        tea.String(p.Endpoint),
	}
	client, err := dyvmsapi.NewClient(config)
	if err != nil {
		return err
	}
	p.Client = client
	return nil
}

func (p *PhoneCall) Call(message string, phoneNumbers []string) error {
	var resultError error
	for _, phoneNumber := range phoneNumbers {
		request := &dyvmsapi.VoiceSingleCallRequest{
			// 接收语音通知的手机号码
			CalledNumber:   tea.String(phoneNumber),
			CallerIdNumber: nil,
			// 语音播报次数
			PlayTimes: tea.Int64(2),
			TtsCode:   tea.String(p.TtsCode),
			TtsParam:  tea.String(message),
		}
		response, err := p.Client.VoiceSingleCall(request)
		if err != nil {
			logc.Errorf(ctx.Ctx, "呼叫失败，号码：%s，内容：%s\n", phoneNumber, message)
			resultError = multierr.Append(resultError, err)
			continue
		}
		if response.Body.Success != tea.Bool(true) {
			logc.Errorf(ctx.Ctx, "呼叫失败，号码：%s，内容：%s，原因：%s", phoneNumber, message, response.Body.String())
			resultError = multierr.Append(resultError, errors.New(*response.Body.AccessDeniedDetail))
			continue
		}
		logc.Info(ctx.Ctx, fmt.Sprintf("呼叫成功，号码：%s，内容：%s\n", phoneNumber, message))
	}
	return resultError
}
