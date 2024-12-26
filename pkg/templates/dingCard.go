package templates

import (
	"fmt"
	models2 "watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

func dingdingTemplate(alert models2.AlertCurEvent, noticeTmpl models2.NoticeTemplateExample) string {

	Title := ParserTemplate("Title", alert, noticeTmpl.Template)
	Footer := ParserTemplate("Footer", alert, noticeTmpl.Template)

	userId := alert.DutyUser

	if alert.DutyUser != "暂无" {
		alert.DutyUser = fmt.Sprintf("@%s", alert.DutyUser)
	}

	t := models2.DingMsg{
		Msgtype: "markdown",
		Markdown: models2.Markdown{
			Title: Title,
			Text: "**" + Title + "**" +
				"\n" + "\n" +
				ParserTemplate("Event", alert, noticeTmpl.Template) +
				"\n" +
				Footer,
		},
		At: models2.At{
			AtUserIds: []string{userId},
			AtMobiles: []string{userId},
			IsAtAll:   false,
		},
	}

	cardContentString := tools.JsonMarshal(t)

	return cardContentString

}
