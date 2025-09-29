package templates

import (
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

func slackTemplate(alert models.AlertCurEvent, noticeTmpl models.NoticeTemplateExample) string {
	t := models.SlackMsgTemplate{
		Text: ParserTemplate("Event", alert, noticeTmpl.Template),
	}

	return tools.JsonMarshalToString(t)
}
