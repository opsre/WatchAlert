package templates

import (
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
)

type Template struct {
	CardContentMsg string
}

// NewTemplate 创建模板
func NewTemplate(ctx *ctx.Context, alert models.AlertCurEvent, route models.Route) Template {
	noticeTmpl := ctx.DB.NoticeTmpl().Get(route.NoticeTmplId)
	switch route.NoticeType {
	case "FeiShu":
		return Template{CardContentMsg: feishuTemplate(alert, noticeTmpl)}
	case "DingDing":
		return Template{CardContentMsg: dingdingTemplate(alert, noticeTmpl)}
	case "Email":
		return Template{CardContentMsg: emailTemplate(alert, noticeTmpl)}
	case "WeChat":
		return Template{CardContentMsg: wechatTemplate(alert, noticeTmpl)}
	case "PhoneCall":
		return Template{CardContentMsg: phoneCallTemplate(alert, noticeTmpl)}
	case "Slack":
		return Template{CardContentMsg: slackTemplate(alert, noticeTmpl)}
	}

	return Template{}
}
