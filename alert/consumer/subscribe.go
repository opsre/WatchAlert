package consumer

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"strings"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/sender"
	"watchAlert/pkg/templates"
	"watchAlert/pkg/tools"
)

type toUser struct {
	Email            string
	NoticeSubject    string
	NoticeTemplateId string
}

// 向已订阅的用户中发送告警消息
func processSubscribe(ctx *ctx.Context, alert *models.AlertCurEvent) error {
	var toUsers []toUser

	// 获取所有用户订阅列表
	subscribes, err := getSubscribes(alert)
	if err != nil {
		return err
	}

	for _, subscribe := range subscribes {
		// Severity检查
		severitySet := make(map[string]struct{})
		for _, sev := range subscribe.SRuleSeverity {
			severitySet[sev] = struct{}{}
		}
		if _, exists := severitySet[alert.Severity]; !exists {
			continue
		}

		// 过滤器检查
		if len(subscribe.SFilter) > 0 {
			allMatched := true
			for _, f := range subscribe.SFilter {
				if !strings.Contains(tools.JsonMarshalToString(alert.Labels), f) && !strings.Contains(alert.Annotations, f) {
					allMatched = false
					break
				}
			}
			if !allMatched {
				continue
			}
		}

		toUsers = append(toUsers, toUser{
			Email:            subscribe.SUserEmail,
			NoticeSubject:    subscribe.SNoticeSubject,
			NoticeTemplateId: subscribe.SNoticeTemplateId,
		})
	}

	return sendToSubscribeUser(ctx, *alert, toUsers)
}

func getSubscribes(alert *models.AlertCurEvent) ([]models.AlertSubscribe, error) {
	list, err := ctx.DB.Subscribe().List(alert.TenantId, alert.RuleId, "")
	if err != nil {
		return nil, fmt.Errorf("获取订阅用户失败, err: %s", err.Error())
	}

	return list, nil
}

func sendToSubscribeUser(ctx *ctx.Context, alert models.AlertCurEvent, toUsers []toUser) error {
	if len(toUsers) <= 0 {
		return nil
	}

	var sem = make(chan struct{}, 10)
	for _, user := range toUsers {
		u := user
		// 插入信号量，超过 10 则阻塞协程启动
		sem <- struct{}{}
		go func(u toUser, sem chan struct{}) {
			defer func() {
				// 释放信号量
				<-sem
			}()
			emailTemp := templates.NewTemplate(ctx, alert, models.AlertNotice{NoticeType: "Email", NoticeTmplId: u.NoticeTemplateId})
			err := sender.NewEmailSender().Send(sender.SendParams{
				IsRecovered: alert.IsRecovered,
				Email: models.Email{
					Subject: u.NoticeSubject,
					To:      []string{u.Email},
					CC:      nil,
				},
				Content: emailTemp.CardContentMsg,
			})
			if err != nil {
				logc.Errorf(ctx.Ctx, fmt.Sprintf("Email: %s, 邮件发送失败, err: %s", u.Email, err.Error()))
			}
		}(u, sem)
	}

	return nil
}
