use watchalert;
INSERT ignore INTO watchalert.notice_template_examples (id,name,description,template,enable_fei_shu_json_card,template_firing,template_recover,notice_type) VALUES
	 ('nt-cqh3uppd6gvj2ctaqd60','飞书通知模版','发送飞书的普通消息模版','{{- define "Title" -}}
{{- if not .IsRecovered -}}
【报警中】- WatchAlert 业务系统 🔥
{{- else -}}
【已恢复】- WatchAlert 业务系统 ✨
{{- end -}}
{{- end }}

{{- define "TitleColor" -}}
{{- if not .IsRecovered -}}
red
{{- else -}}
green
{{- end -}}
{{- end }}

{{ define "Event" -}}
{{- if not .IsRecovered -}}
**🤖 报警类型:** ${rule_name}
**🫧 报警指纹:** ${fingerprint}
**📌 报警等级:** ${severity}
**🖥 报警主机:** ${labels.instance}
**🕘 开始时间:** ${first_trigger_time_format}
**👤 值班人员:** ${duty_user}
**📝 报警事件:** ${annotations}
[查看事件](http://localhost:3000/faultCenter/detail/${faultCenterId}?tab=1&query=${rule_name})
{{- else -}}
**🤖 报警类型:** ${rule_name}
**🫧 报警指纹:** ${fingerprint}
**📌 报警等级:** ${severity}
**🖥 报警主机:** ${labels.instance}
**🕘 开始时间:** ${first_trigger_time_format}
**🕘 恢复时间:** ${recover_time_format}
**👤 值班人员:** ${duty_user}
**📝 报警事件:** ${annotations}
[查看事件](http://localhost:3000/faultCenter/detail/${faultCenterId}?tab=1&query=${rule_name})
{{- end -}}
{{ end }}

{{- define "Footer" -}}
🧑‍💻 WatchAlert - 运维团队
{{- end }}',0,'','','FeiShu'),
	 ('nt-cqh4361d6gvj80netqk0','飞书卡片通知模版','发送飞书的高级消息卡片模版','',1,'{
  "elements": [
    {
      "tag": "column_set",
      "flexMode": "none",
      "background_style": "default",
      "text": {
        "content": "",
        "tag": ""
      },
      "actions": null,
      "columns": [],
      "elements": null
    },
    {
      "tag": "column_set",
      "flexMode": "none",
      "background_style": "default",
      "text": {
        "content": "",
        "tag": ""
      },
      "actions": null,
      "columns": [
        {
          "tag": "column",
          "width": "weighted",
          "weight": 1,
          "vertical_align": "top",
          "elements": [
            {
              "tag": "div",
              "text": {
                "content": "**🫧 报警指纹：**\n${fingerprint}",
                "tag": "lark_md"
              }
            }
          ]
        },
        {
          "tag": "column",
          "width": "weighted",
          "weight": 1,
          "vertical_align": "top",
          "elements": [
            {
              "tag": "div",
              "text": {
                "content": "**🕘 开始时间：**\n${first_trigger_time_format}",
                "tag": "lark_md"
              }
            }
          ]
        }
      ],
      "elements": null
    },
    {
      "tag": "column_set",
      "flexMode": "none",
      "background_style": "default",
      "text": {
        "content": "",
        "tag": ""
      },
      "actions": null,
      "columns": [
        {
          "tag": "column",
          "width": "weighted",
          "weight": 1,
          "vertical_align": "top",
          "elements": [
            {
              "tag": "div",
              "text": {
                "content": "**🖥 报警主机：**\n${labels.instance}",
                "tag": "lark_md"
              }
            }
          ]
        }
      ],
      "elements": null
    },
    {
      "tag": "hr"
    },
    {
      "tag": "column_set",
      "flexMode": "none",
      "background_style": "default",
      "text": {
        "content": "",
        "tag": ""
      },
      "actions": null,
      "columns": [
        {
          "tag": "column",
          "width": "weighted",
          "weight": 1,
          "vertical_align": "top",
          "elements": [
            {
              "tag": "div",
              "text": {
                "content": "**⛩️ 故障中心：**\n${faultCenter.name}",
                "tag": "lark_md"
              }
            }
          ]
        },
        {
          "tag": "column",
          "width": "weighted",
          "weight": 1,
          "vertical_align": "top",
          "elements": [
            {
              "tag": "div",
              "text": {
                "content": "**👤 值班人员：**\n${duty_user}",
                "tag": "lark_md"
              }
            }
          ]
        }
      ],
      "elements": null
    },
    {
      "tag": "column_set",
      "flexMode": "none",
      "background_style": "default",
      "text": {
        "content": "",
        "tag": ""
      },
      "actions": null,
      "columns": [],
      "elements": null
    },
    {
      "tag": "hr"
    },
    {
      "tag": "column_set",
      "flexMode": "none",
      "background_style": "default",
      "text": {
        "content": "",
        "tag": ""
      },
      "actions": null,
      "columns": [
        {
          "tag": "column",
          "width": "weighted",
          "weight": 1,
          "vertical_align": "top",
          "elements": [
            {
              "tag": "div",
              "text": {
                "content": "**📝 报警事件：**\n${annotations}",
                "tag": "lark_md"
              }
            }
          ]
        }
      ],
      "elements": null
    },
    {
      "tag": "hr",
      "flexMode": "",
      "background_style": "",
      "text": {
        "content": "",
        "tag": ""
      },
      "actions": null,
      "columns": null,
      "elements": null
    },
    {
      "tag": "note",
      "flexMode": "",
      "background_style": "",
      "text": {
        "content": "",
        "tag": ""
      },
      "actions": null,
      "columns": null,
      "elements": [
        {
          "tag": "plain_text",
          "content": "🧑‍💻 WatchAlert - 运维团队"
        }
      ]
    }
  ],
  "header": {
    "template": "red",
    "title": {
      "content": "【 ${severity} 报警中】- ${rule_name} 🔥",
      "tag": "plain_text"
    }
  },
  "tag": ""
}','{
  "elements": [
    {
      "tag": "column_set",
      "flexMode": "none",
      "background_style": "default",
      "text": {
        "content": "",
        "tag": ""
      },
      "actions": null,
      "columns": [],
      "elements": null
    },
    {
      "tag": "column_set",
      "flexMode": "none",
      "background_style": "default",
      "text": {
        "content": "",
        "tag": ""
      },
      "actions": null,
      "columns": [
        {
          "tag": "column",
          "width": "weighted",
          "weight": 1,
          "vertical_align": "top",
          "elements": [
            {
              "tag": "div",
              "text": {
                "content": "**🫧 报警指纹：**\n${fingerprint}",
                "tag": "lark_md"
              }
            }
          ]
        },
        {
          "tag": "column",
          "width": "weighted",
          "weight": 1,
          "vertical_align": "top",
          "elements": [
            {
              "tag": "div",
              "text": {
                "content": "**🕘 开始时间：**\n${first_trigger_time_format}",
                "tag": "lark_md"
              }
            }
          ]
        }
      ],
      "elements": null
    },
    {
      "tag": "column_set",
      "flexMode": "none",
      "background_style": "default",
      "text": {
        "content": "",
        "tag": ""
      },
      "actions": null,
      "columns": [
        {
          "tag": "column",
          "width": "weighted",
          "weight": 1,
          "vertical_align": "top",
          "elements": [
            {
              "tag": "div",
              "text": {
                "content": "**🖥 报警主机：**\n${labels.instance}",
                "tag": "lark_md"
              }
            }
          ]
        },
        {
          "tag": "column",
          "width": "weighted",
          "weight": 1,
          "vertical_align": "top",
          "elements": [
            {
              "tag": "div",
              "text": {
                "content": "**🕘 恢复时间：**\n${recover_time_format}",
                "tag": "lark_md"
              }
            }
          ]
        }
      ],
      "elements": null
    },
    {
      "tag": "hr"
    },
    {
      "tag": "column_set",
      "flexMode": "none",
      "background_style": "default",
      "text": {
        "content": "",
        "tag": ""
      },
      "actions": null,
      "columns": [
        {
          "tag": "column",
          "width": "weighted",
          "weight": 1,
          "vertical_align": "top",
          "elements": [
            {
              "tag": "div",
              "text": {
                "content": "**⛩️ 故障中心：**\n${faultCenter.name}",
                "tag": "lark_md"
              }
            }
          ]
        },
        {
          "tag": "column",
          "width": "weighted",
          "weight": 1,
          "vertical_align": "top",
          "elements": [
            {
              "tag": "div",
              "text": {
                "content": "**👤 值班人员：**\n${duty_user}",
                "tag": "lark_md"
              }
            }
          ]
        }
      ],
      "elements": null
    },
    {
      "tag": "hr",
      "flexMode": "",
      "background_style": "",
      "text": {
        "content": "",
        "tag": ""
      },
      "actions": null,
      "columns": null,
      "elements": null
    },
    {
      "tag": "column_set",
      "flexMode": "none",
      "background_style": "default",
      "text": {
        "content": "",
        "tag": ""
      },
      "actions": null,
      "columns": [
        {
          "tag": "column",
          "width": "weighted",
          "weight": 1,
          "vertical_align": "top",
          "elements": [
            {
              "tag": "div",
              "text": {
                "content": "**📝 报警事件：**\n${annotations}",
                "tag": "lark_md"
              }
            }
          ]
        }
      ],
      "elements": null
    },
    {
      "tag": "hr",
      "flexMode": "",
      "background_style": "",
      "text": {
        "content": "",
        "tag": ""
      },
      "actions": null,
      "columns": null,
      "elements": null
    },
    {
      "tag": "note",
      "flexMode": "",
      "background_style": "",
      "text": {
        "content": "",
        "tag": ""
      },
      "actions": null,
      "columns": null,
      "elements": [
        {
          "tag": "plain_text",
          "content": "🧑‍💻 WatchAlert - 运维团队"
        }
      ]
    }
  ],
  "header": {
    "template": "green",
    "title": {
      "content": "【 ${severity} 已恢复】- ${rule_name} ✨",
      "tag": "plain_text"
    }
  },
  "tag": ""
}','FeiShu'),
	 ('nt-cqh4599d6gvj80netql0','邮件通知模版','发送邮件的普通消息模版','{{ define "Event" -}}
{{- if not .IsRecovered -}}
<p>==========<strong>告警通知</strong>==========</p>
<strong>🤖 报警类型:</strong> ${rule_name}<br>
<strong>🫧 报警指纹:</strong> ${fingerprint}<br>
<strong>📌 报警等级:</strong> ${severity}<br>
<strong>🖥 报警主机:</strong> ${labels.node_name}<br>
<strong>🧚 容器名称:</strong> ${labels.pod}<br>
<strong>☘️ 业务环境:</strong> ${labels.namespace}<br>
<strong>🕘 开始时间:</strong> ${first_trigger_time_format}<br>
<strong>👤 值班人员:</strong> ${duty_user}<br>
<strong>📝 报警事件:</strong> ${annotations}<br>
[查看事件](http://localhost:3000/faultCenter/detail/${faultCenterId}?tab=1&query=${rule_name})
{{- else -}}
<p>==========<strong>恢复通知</strong>==========</p>
<strong>🤖 报警类型:</strong> ${rule_name}<br>
<strong>🫧 报警指纹:</strong> ${fingerprint}<br>
<strong>📌 报警等级:</strong> ${severity}<br>
<strong>🖥 报警主机:</strong> ${labels.node_name}<br>
<strong>🧚 容器名称:</strong> ${labels.pod}<br>
<strong>☘️ 业务环境:</strong> ${labels.namespace}<br>
<strong>🕘 开始时间:</strong> ${first_trigger_time_format}<br>
<strong>🕘 恢复时间:</strong> ${recover_time_format}<br>
<strong>👤 值班人员:</strong> ${duty_user}<br>
<strong>📝 报警事件:</strong> ${annotations}<br>
[查看事件](http://localhost:3000/faultCenter/detail/${faultCenterId}?tab=1&query=${rule_name})
{{- end -}}
{{ end }}',0,'','','Email'),
	 ('nt-crscirlvi7nhfu2tpf00','钉钉通知模版','发送钉钉的普通消息模版','{{- define "Title" -}}
{{- if not .IsRecovered -}}
【报警中】- WatchAlert 业务系统 🔥
{{- else -}}
【已恢复】- WatchAlert 业务系统 ✨
{{- end -}}
{{- end }}

{{- define "TitleColor" -}}
{{- if not .IsRecovered -}}
red
{{- else -}}
green
{{- end -}}
{{- end }}

{{ define "Event" -}}
{{- if not .IsRecovered -}}
&nbsp;**🔔 报警类型:** ${rule_name}<br>
**🔐 报警指纹:** ${fingerprint}<br>
**🚨 报警等级:** ${severity}<br>
**🖥 报警主机:** ${labels.instance}<br>
**🕘 开始时间:** ${first_trigger_time_format}<br>
**🧑‍🔧 值班人员:** ${duty_user}<br>
**📝 报警事件:** ${annotations}<br>
**👀 查看事件:** [点击跳转](http://localhost:3000/faultCenter/detail/${faultCenterId}?tab=1&query=${rule_name})<br>
{{- else -}}
&nbsp;**🔔 报警类型:** ${rule_name}<br>
**🔐 报警指纹:** ${fingerprint}<br>
**🚨 报警等级:** ${severity}<br>
**🖥 报警主机:** ${labels.instance}<br>
**🕘 开始时间:** ${first_trigger_time_format}<br>
**🕘 恢复时间:** ${recover_time_format}<br>
**🧑‍🔧 值班人员:** ${duty_user}<br>
**📝 报警事件:** ${annotations}<br>
**👀 查看事件:** [点击跳转](http://localhost:3000/faultCenter/detail/${faultCenterId}?tab=1&query=${rule_name})<br>
{{- end -}}
{{ end }}

{{- define "Footer" -}}
🧑‍💻 WatchAlert - 运维团队
{{- end }}',0,'','','DingDing'),
	 ('nt-cte1re5vi7ngs77mh190','企微通知模版','发送企业微信的普通消息模版','{{- define "Title" -}}
{{- if not .IsRecovered -}}
【报警中】- WatchAlert 业务系统 🔥
{{- else -}}
【已恢复】- WatchAlert 业务系统 ✨
{{- end -}}
{{- end }}

{{- define "TitleColor" -}}
{{- if not .IsRecovered -}}
red
{{- else -}}
green
{{- end -}}
{{- end }}

{{ define "Event" -}}
{{- if not .IsRecovered -}}
>**🤖 报警类型:** ${rule_name}
>**🫧 报警指纹:** ${fingerprint}
>**📌 报警等级:** ${severity}
>**🖥 报警主机:** ${labels.instance}
>**🕘 开始时间:** ${first_trigger_time_format}
>**👤 值班人员:** ${duty_user}
>**📝 报警事件:** ${annotations}
[查看事件](http://localhost:3000/faultCenter/detail/${faultCenterId}?tab=1&query=${rule_name})
{{- else -}}
>**🤖 报警类型:** ${rule_name}
>**🫧 报警指纹:** ${fingerprint}
>**📌 报警等级:** ${severity}
>**🖥 报警主机:** ${labels.instance}
>**🕘 开始时间:** ${first_trigger_time_format}
>**🕘 恢复时间:** ${recover_time_format}
>**👤 值班人员:** ${duty_user}
>**📝 报警事件:** ${annotations}
[查看事件](http://localhost:3000/faultCenter/detail/${faultCenterId}?tab=1&query=${rule_name})
{{- end -}}
{{ end }}

{{- define "Footer" -}}
🧑‍💻 WatchAlert - 运维团队
{{- end }}',0,'','','WeChat'),
    ('nt-crscirlvi7nhbb2tpf01','日志通知模版','日志类告警通知模版','{{- define "Title" -}}
{{- if not .IsRecovered -}}
【报警中】- WatchAlert 业务系统 🔥
{{- else -}}
【已恢复】- WatchAlert 业务系统 ✨
{{- end -}}
{{- end }}

{{- define "TitleColor" -}}
{{- if not .IsRecovered -}}
red
{{- else -}}
green
{{- end -}}
{{- end }}

{{ define "Event" -}}
{{- if not .IsRecovered -}}
**🤖 报警类型:** ${rule_name}</br>
**📌 报警等级:** ${severity}</br>
**🕘 开始时间:** ${first_trigger_time_format}</br>
**👤 值班人员:** ${duty_user}</br>
**📝 服务名称:** ${log.app}</br>
**📝 TraceId:** ${log.trace_id}</br>
**📝 日志内容:** ${log.message}</br>
{{- else -}}
**🤖 报警类型:** ${rule_name}</br>
**📌 报警等级:** ${severity}</br>
**🕘 开始时间:** ${first_trigger_time_format}</br>
**🕘 恢复时间:** ${recover_time_format}</br>
**👤 值班人员:** ${duty_user}</br>
{{- end -}}
{{ end }}

{{- define "Footer" -}}
🧑‍💻 WatchAlert - 运维团队
{{- end }}',0,'','','FeiShu');
