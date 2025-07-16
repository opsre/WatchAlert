package templates

import (
	"fmt"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

func slackTemplate(alert models.AlertCurEvent, noticeTmpl models.NoticeTemplateExample) string {
	if alert.DutyUser != "暂无" {
		alert.DutyUser = fmt.Sprintf("%s", alert.DutyUser)
	}

	t := models.SlackMsgTemplate{
		Text: ParserTemplate("Event", alert, noticeTmpl.Template),
	}

	cardContentString := tools.JsonMarshal(t)
	return cardContentString
}
