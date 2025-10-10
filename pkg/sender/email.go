package sender

import (
	"errors"
	"fmt"
	"watchAlert/internal/ctx"
	"watchAlert/pkg/client"
)

// EmailSender 邮件发送策略
type EmailSender struct{}

func NewEmailSender() SendInter {
	return &EmailSender{}
}

func (e *EmailSender) Send(params SendParams) error {
	setting, err := ctx.DB.Setting().Get()
	if err != nil {
		return errors.New("获取 系统配置/邮箱配置 失败: " + err.Error())
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

func (e *EmailSender) Test(params SendParams) error {
	setting, err := ctx.DB.Setting().Get()
	if err != nil {
		return errors.New("获取 系统配置/邮箱配置 失败: " + err.Error())
	}

	eCli := client.NewEmailClient(setting.EmailConfig.ServerAddress, setting.EmailConfig.Email, setting.EmailConfig.Token, setting.EmailConfig.Port)
	return eCli.Send(params.Email.To, params.Email.CC, "WatchAlert 消息测试", []byte(RobotTestContent))
}
