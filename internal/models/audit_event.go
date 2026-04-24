package models

// AuditEventMap 审计日志接口与事件名称映射表（使用完整API路径作为key）
var AuditEventMap = map[string]string{
	// ========== 用户相关 ==========
	"/api/w8t/user/userUpdate":     "更新用户",
	"/api/w8t/user/userDelete":     "删除用户",
	"/api/w8t/user/userChangePass": "修改密码",

	// ========== 租户相关 ==========
	"/api/w8t/tenant/createTenant":         "创建租户",
	"/api/w8t/tenant/updateTenant":         "更新租户",
	"/api/w8t/tenant/deleteTenant":         "删除租户",
	"/api/w8t/tenant/addUsersToTenant":     "添加用户到租户",
	"/api/w8t/tenant/delUsersOfTenant":     "从租户删除用户",
	"/api/w8t/tenant/changeTenantUserRole": "更改租户用户角色",

	// ========== 告警规则相关 ==========
	"/api/w8t/rule/ruleCreate":       "创建告警规则",
	"/api/w8t/rule/ruleUpdate":       "更新告警规则",
	"/api/w8t/rule/ruleDelete":       "删除告警规则",
	"/api/w8t/rule/import":           "导入告警规则",
	"/api/w8t/rule/ruleChangeStatus": "修改告警规则状态",
	"/api/w8t/rule/change":           "修改告警规则",

	// ========== 规则组相关 ==========
	"/api/w8t/ruleGroup/ruleGroupCreate": "创建告警规则组",
	"/api/w8t/ruleGroup/ruleGroupUpdate": "更新告警规则组",
	"/api/w8t/ruleGroup/ruleGroupDelete": "删除告警规则组",

	// ========== 规则模版相关 ==========
	"/api/w8t/ruleTmpl/ruleTmplCreate": "创建规则模版",
	"/api/w8t/ruleTmpl/ruleTmplUpdate": "更新规则模版",
	"/api/w8t/ruleTmpl/ruleTmplDelete": "删除规则模版",

	// ========== 规则模版组相关 ==========
	"/api/w8t/ruleTmplGroup/ruleTmplGroupCreate": "创建规则模版组",
	"/api/w8t/ruleTmplGroup/ruleTmplGroupUpdate": "更新规则模版组",
	"/api/w8t/ruleTmplGroup/ruleTmplGroupDelete": "删除规则模版组",

	// ========== 静默规则相关 ==========
	"/api/w8t/silence/silenceCreate": "创建静默规则",
	"/api/w8t/silence/silenceUpdate": "更新静默规则",
	"/api/w8t/silence/silenceDelete": "删除静默规则",

	// ========== 通知对象相关 ==========
	"/api/w8t/notice/noticeCreate": "创建通知对象",
	"/api/w8t/notice/noticeUpdate": "更新通知对象",
	"/api/w8t/notice/noticeDelete": "删除通知对象",
	"/api/w8t/notice/noticeTest":   "测试通知",

	// ========== 通知模版相关 ==========
	"/api/w8t/noticeTemplate/noticeTemplateCreate": "创建通知模版",
	"/api/w8t/noticeTemplate/noticeTemplateUpdate": "更新通知模版",
	"/api/w8t/noticeTemplate/noticeTemplateDelete": "删除通知模版",

	// ========== 数据源相关 ==========
	"/api/w8t/datasource/dataSourceCreate":      "创建数据源",
	"/api/w8t/datasource/dataSourceUpdate":      "更新数据源",
	"/api/w8t/datasource/dataSourceDelete":      "删除数据源",
	"/api/w8t/datasource/dataSourcePing":        "测试数据源连接",
	"/api/w8t/datasource/searchViewLogsContent": "查询日志内容",

	// ========== 值班管理相关 ==========
	"/api/w8t/dutyManage/dutyManageCreate": "创建值班管理",
	"/api/w8t/dutyManage/dutyManageUpdate": "更新值班管理",
	"/api/w8t/dutyManage/dutyManageDelete": "删除值班管理",

	// ========== 值班日历相关 ==========
	"/api/w8t/calendar/calendarCreate": "创建值班表",
	"/api/w8t/calendar/calendarUpdate": "更新值班表",
	"/api/w8t/calendar/calendarDelete": "删除值班表",

	// ========== 仪表盘相关 ==========
	"/api/w8t/dashboard/createFolder": "创建仪表盘目录",
	"/api/w8t/dashboard/updateFolder": "更新仪表盘目录",
	"/api/w8t/dashboard/deleteFolder": "删除仪表盘目录",

	// ========== 拨测相关 ==========
	"/api/w8t/probing/createProbing": "创建拨测规则",
	"/api/w8t/probing/updateProbing": "更新拨测规则",
	"/api/w8t/probing/deleteProbing": "删除拨测规则",
	"/api/w8t/probing/onceProbing":   "手动执行拨测",
	"/api/w8t/probing/changeState":   "修改拨测状态",

	// ========== 故障中心相关 ==========
	"/api/w8t/faultCenter/faultCenterCreate": "创建故障中心",
	"/api/w8t/faultCenter/faultCenterUpdate": "更新故障中心",
	"/api/w8t/faultCenter/faultCenterDelete": "删除故障中心",
	"/api/w8t/faultCenter/faultCenterReset":  "修改故障中心基本信息",

	// ========== 用户角色相关 ==========
	"/api/w8t/role/roleCreate": "创建用户角色",
	"/api/w8t/role/roleUpdate": "更新用户角色",
	"/api/w8t/role/roleDelete": "删除用户角色",

	// ========== 订阅相关 ==========
	"/api/w8t/subscribe/createSubscribe": "创建告警订阅",
	"/api/w8t/subscribe/deleteSubscribe": "删除告警订阅",

	// ========== 系统设置相关 ==========
	"/api/w8t/setting/saveSystemSetting": "保存系统设置",
	"/api/w8t/setting/syncLdapUser":      "同步LDAP用户",

	// ========== 事件相关 ==========
	"/api/w8t/event/delete":        "删除告警事件",
	"/api/w8t/event/addComment":    "添加评论",
	"/api/w8t/event/deleteComment": "删除评论",
	"/api/w8t/event/process":       "处理告警事件",

	// ========== 记录规则相关 ==========
	"/api/w8t/recordingRule/recordingRuleCreate":       "创建记录规则",
	"/api/w8t/recordingRule/recordingRuleUpdate":       "更新记录规则",
	"/api/w8t/recordingRule/recordingRuleDelete":       "删除记录规则",
	"/api/w8t/recordingRule/recordingRuleChangeStatus": "修改记录规则状态",

	// ========== 记录规则组相关 ==========
	"/api/w8t/recordingRuleGroup/recordingRuleGroupCreate": "创建记录规则组",
	"/api/w8t/recordingRuleGroup/recordingRuleGroupUpdate": "更新记录规则组",
	"/api/w8t/recordingRuleGroup/recordingRuleGroupDelete": "删除记录规则组",

	// ========== API Key 相关 ==========
	"/api/w8t/apiKey/create": "创建API Key",
	"/api/w8t/apiKey/update": "更新API Key",
	"/api/w8t/apiKey/delete": "删除API Key",

	// ========== AI 相关 ==========
	"/api/w8t/ai/chat": "AI对话",
}
