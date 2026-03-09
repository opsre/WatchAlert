package sender

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
	"watchAlert/internal/ctx"

	"github.com/jordan-wright/email"
)

// EmailSender 邮件发送策略
type EmailSender struct {
	ServerAddr string
	Port       int
	Email      *email.Email
	Auth       smtp.Auth
}

func NewEmailSender() (SendInter, error) {
	setting, err := ctx.DB.Setting().Get()
	if err != nil {
		return nil, errors.New("获取 系统配置/邮箱配置 失败: " + err.Error())
	}

	e := email.NewEmail()
	auth := smtp.PlainAuth("", setting.EmailConfig.Email, setting.EmailConfig.Token, setting.EmailConfig.ServerAddress)
	e.From = fmt.Sprintf("WatchAlert<%s>", setting.EmailConfig.Email)

	return &EmailSender{
		ServerAddr: setting.EmailConfig.ServerAddress,
		Port:       setting.EmailConfig.Port,
		Email:      e,
		Auth:       auth,
	}, nil
}

func (e *EmailSender) Send(params SendParams) error {
	if params.IsRecovered {
		params.Email.Subject = params.Email.Subject + "「已恢复」"
	} else {
		params.Email.Subject = params.Email.Subject + "「报警中」"
	}
	err := e.post(params.Email.To, params.Email.CC, params.Email.Subject, []byte(params.Content))
	if err != nil {
		return fmt.Errorf("%s, %s", err.Error(), "Content: "+params.Content)
	}

	return nil
}

func (e *EmailSender) Test(params SendParams) error {
	return e.post(params.Email.To, params.Email.CC, "WatchAlert 消息测试", []byte(RobotTestContent))
}

func (e *EmailSender) post(to, cc []string, subject string, msg []byte) error {
	e.Email.To = to
	e.Email.Cc = cc
	e.Email.HTML = msg
	e.Email.Subject = subject

	addr := fmt.Sprintf("%s:%d", e.ServerAddr, e.Port)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         e.ServerAddr,
	}

	// 如果端口是 465，使用标准的 SSL/TLS 加密
	if e.Port == 465 {
		return e.Email.SendWithTLS(addr, e.Auth, tlsConfig)
	}

	return e.Email.Send(addr, e.Auth)
}
