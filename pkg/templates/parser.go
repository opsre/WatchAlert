package templates

import (
	"bytes"
	"context"
	"fmt"
	"text/template"
	"time"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"

	"github.com/zeromicro/go-zero/core/logc"
)

// ParserTemplate 处理告警推送的消息模版
func ParserTemplate(defineName string, alert models.AlertCurEvent, templateStr string) string {
	// 预处理时间格式化
	alert = prepareAlertData(alert)

	// 解析模板
	tmpl, err := template.New("tmpl").Parse(templateStr)
	if err != nil {
		logc.Error(context.Background(), "模板解析失败", "error", err, "template", templateStr)
		return ""
	}

	return renderNamedTemplate(tmpl, defineName, alert)
}

// prepareAlertData 预处理告警数据，格式化时间字段
func prepareAlertData(alert models.AlertCurEvent) models.AlertCurEvent {
	alert.FirstTriggerTimeFormat = time.Unix(alert.FirstTriggerTime, 0).Format(global.Layout)
	alert.RecoverTimeFormat = time.Unix(alert.RecoverTime, 0).Format(global.Layout)
	return alert
}

// renderNamedTemplate 渲染模板
func renderNamedTemplate(tmpl *template.Template, name string, alert models.AlertCurEvent) string {
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, alert); err != nil {
		// 如果命名模板不存在，尝试直接执行主模板
		if err := tmpl.Execute(&buf, alert); err != nil {
			logc.Error(context.Background(), fmt.Sprintf("%s模板执行失败", name), "error", err)
			return ""
		}
	}

	// 解析变量并返回
	data := tools.ConvertStructToMap(alert)
	return tools.ParserVariables(buf.String(), data)
}
