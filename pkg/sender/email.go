package sender

import (
	"errors"
	"fmt"
	"watchAlert/pkg/client"
	"watchAlert/pkg/ctx"
)

// EmailSender 邮件发送策略
type EmailSender struct{}

func NewEmailSender() SendInter {
	return &EmailSender{}
}

func (e *EmailSender) Send(params SendParams) error {
	setting, err := ctx.DB.Setting().Get()
	if err != nil {
		return errors.New("获取系统配置失败: " + err.Error())
	}
	eCli := client.NewEmailClient(setting.EmailConfig.ServerAddress, setting.EmailConfig.Email, setting.EmailConfig.Token, setting.EmailConfig.Port)
	if params.IsRecovered {
		params.Email.Subject = params.Email.Subject + "「已恢复」"
	} else {
		params.Email.Subject = params.Email.Subject + "「报警中」"
	}
	err = eCli.Send(params.Email.To, params.Email.CC, params.Email.Subject, []byte(params.Content))
	if err != nil {
		return fmt.Errorf("%s, %s", err.Error(), "Content: "+params.Content)
	}

	return nil
}
