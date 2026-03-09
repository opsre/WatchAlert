package templates

import (
	"bytes"
	"context"
	"text/template"
	"time"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"

	"github.com/zeromicro/go-zero/core/logc"
)

// ParserTemplate 处理告警推送的消息模版
func ParserTemplate(defineName string, alert models.AlertCurEvent, templateStr string) string {
	// 1. 定义模板函数
	funcMap := template.FuncMap{
		// 时间戳转格式化字符串: {{ .FirstTriggerTime | formatTime }}
		"formatTime": func(timestamp int64) string {
			if timestamp == 0 {
				return "-"
			}
			return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
		},
		// 计算持续时间: {{ duration .FirstTriggerTime }}
		"duration": func(first int64) string {
			cur := time.Now().Unix()
			if first == 0 || cur == 0 || cur < first {
				return "0s"
			}
			d := time.Duration(cur-first) * time.Second
			return d.String()
		},
	}

	// 2. 解析模板并注入函数
	tmpl, err := template.New("tmpl").Funcs(funcMap).Parse(templateStr)
	if err != nil {
		logc.Errorf(context.Background(), "模板解析失败: %v, template: %s", err, templateStr)
		return ""
	}

	return renderNamedTemplate(tmpl, defineName, alert)
}

// renderNamedTemplate 渲染模板
func renderNamedTemplate(tmpl *template.Template, name string, alert models.AlertCurEvent) string {
	var buf bytes.Buffer
	// 尝试执行指定的 define 块，如果失败则执行整个模板
	if err := tmpl.ExecuteTemplate(&buf, name, alert); err != nil {
		if err := tmpl.Execute(&buf, alert); err != nil {
			logc.Errorf(context.Background(), "%s 模板执行失败: %v", name, err)
			return ""
		}
	}

	// 解析变量并返回
	data := tools.ConvertStructToMap(alert)
	return tools.ParserVariables(buf.String(), data)
}
