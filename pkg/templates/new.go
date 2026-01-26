package templates

import (
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
)

type Template struct {
	CardContentMsg string
}

// NewTemplate 创建模板
func NewTemplate(ctx *ctx.Context, alert models.AlertCurEvent, route models.Route) (Template, error) {
	noticeTmpl, err := ctx.DB.NoticeTmpl().Get(route.NoticeTmplId)
	if err != nil {
		return Template{}, err
	}
	switch route.NoticeType {
	case "FeiShu":
		return Template{CardContentMsg: feishuTemplate(alert, noticeTmpl)}, nil
	case "DingDing":
		return Template{CardContentMsg: dingdingTemplate(alert, noticeTmpl)}, nil
	case "Email":
		return Template{CardContentMsg: emailTemplate(alert, noticeTmpl)}, nil
	case "WeChat":
		return Template{CardContentMsg: wechatTemplate(alert, noticeTmpl)}, nil
	case "PhoneCall":
		return Template{CardContentMsg: phoneCallTemplate(alert, noticeTmpl)}, nil
	case "Slack":
		return Template{CardContentMsg: slackTemplate(alert, noticeTmpl)}, nil
	}

	return Template{}, nil
}
