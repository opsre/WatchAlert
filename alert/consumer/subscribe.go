package consumer

import (
	"fmt"
	"strings"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
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
func processSubscribe(ctx *ctx.Context, alert *models.AlertCurEvent, notice models.AlertNotice) error {
	list, err := ctx.DB.Subscribe().List(models.AlertSubscribeQuery{
		STenantId: alert.TenantId,
		Query:     alert.RuleId,
	})
	if err != nil {
		return fmt.Errorf("获取订阅用户失败, err: %s", err.Error())
	}

	notice.NoticeType = "Email"
	var toUsers []toUser
	for _, s := range list {
		var foundSeverity, foundFilter bool
		for _, severity := range s.SRuleSeverity {
			if severity == alert.Severity {
				foundSeverity = true
				break
			}
		}

		if foundSeverity {
			if len(s.SFilter) > 0 {
				for _, f := range s.SFilter {
					if strings.Contains(tools.JsonMarshal(alert.Metric), f) || strings.Contains(alert.Annotations, f) {
						foundFilter = true
						break
					}
				}
			} else {
				foundFilter = true
			}

			if foundFilter {
				toUsers = append(toUsers, toUser{
					Email:            s.SUserEmail,
					NoticeSubject:    s.SNoticeSubject,
					NoticeTemplateId: s.SNoticeTemplateId,
				})
			}
		}
	}

	if len(toUsers) > 0 {
		for _, u := range toUsers {
			notice.NoticeTmplId = u.NoticeTemplateId
			emailTemp := templates.NewTemplate(ctx, *alert, notice)

			err = sender.NewEmailSender().Send(sender.SendParams{
				IsRecovered: alert.IsRecovered,
				Email: models.Email{
					Subject: u.NoticeSubject,
					To:      []string{u.Email},
					CC:      nil,
				},
				Content: emailTemp.CardContentMsg,
			})
			if err != nil {
				return fmt.Errorf("邮件发送失败, err: %s", err.Error())
			}
		}
	}

	return nil
}
