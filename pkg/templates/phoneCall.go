package templates

import "watchAlert/internal/models"

func phoneCallTemplate(alert models.AlertCurEvent, noticeTmpl models.NoticeTemplateExample) string {
	return ParserTemplate("Event", alert, noticeTmpl.Template)
}
