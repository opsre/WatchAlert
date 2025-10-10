package templates

import (
	"strings"
	models2 "watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

func dingdingTemplate(alert models2.AlertCurEvent, noticeTmpl models2.NoticeTemplateExample) string {
	Title := ParserTemplate("Title", alert, noticeTmpl.Template)
	Footer := ParserTemplate("Footer", alert, noticeTmpl.Template)

	dutyUser := alert.DutyUser
	var dutyUsers []string
	for _, user := range strings.Split(dutyUser, " ") {
		u := strings.Trim(user, "@")
		dutyUsers = append(dutyUsers, u)
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
			AtUserIds: dutyUsers,
			AtMobiles: dutyUsers,
			IsAtAll:   false,
		},
	}

	if strings.Trim(alert.DutyUser, " ") == "all" {
		t.At = models2.At{
			AtUserIds: []string{},
			AtMobiles: []string{},
			IsAtAll:   true,
		}
	}

	return tools.JsonMarshalToString(t)
}
