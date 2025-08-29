package client

import (
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
	"strconv"
)

type EmailClient struct {
	ServerAddr string
	Port       int
	Email      *email.Email
	Auth       smtp.Auth
}

func NewEmailClient(serverAddr, username, password string, port int) EmailClient {
	e := email.NewEmail()
	auth := smtp.PlainAuth("", username, password, serverAddr)
	e.From = fmt.Sprintf("WatchAlert<%s>", username)

	return EmailClient{
		Email:      e,
		Auth:       auth,
		ServerAddr: serverAddr,
		Port:       port,
	}
}

func (a EmailClient) Send(to, cc []string, subject string, msg []byte) error {
	a.Email.To = to
	a.Email.Cc = cc
	a.Email.HTML = msg
	a.Email.Subject = subject
	port := strconv.FormatInt(int64(a.Port), 10)
	//err := a.Email.SendWithTLS(a.ServerAddr+":"+port, a.Auth, &tls.Config{ServerName: a.ServerAddr})
	err := a.Email.Send(a.ServerAddr+":"+port, a.Auth)
	if err != nil {
		return err
	}

	return nil
}
