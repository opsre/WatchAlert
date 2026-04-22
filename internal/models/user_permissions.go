package models

type UserPermissions struct {
	Key string `json:"key"`
	API string `json:"api"`
}

func PermissionsInfo() map[string]UserPermissions {
	return map[string]UserPermissions{
		"ruleSearch": {
			Key: "预览告警规则详情",
			API: "/api/w8t/rule/ruleSearch",
		},
		"calendarCreate": {
			Key: "创建值班表",
			API: "/api/w8t/calendar/calendarCreate",
		},
		"calendarSearch": {
			Key: "预览值班表",
			API: "/api/w8t/calendar/calendarSearch",
		},
		"calendarUpdate": {
			Key: "更新值班表",
			API: "/api/w8t/calendar/calendarUpdate",
		},
		"dataSourceCreate": {
			Key: "创建数据源",
			API: "/api/w8t/datasource/dataSourceCreate",
		},
		"dataSourceDelete": {
			Key: "删除数据源",
			API: "/api/w8t/datasource/dataSourceDelete",
		},
		"dataSourceGet": {
			Key: "预览数据源详情",
			API: "/api/w8t/datasource/dataSourceGet",
		},
		"dataSourceList": {
			Key: "查看数据源列表",
			API: "/api/w8t/datasource/dataSourceList",
		},
		"dataSourceUpdate": {
			Key: "更新数据源",
			API: "/api/w8t/datasource/dataSourceUpdate",
		},
		"dutyManageCreate": {
			Key: "创建值班管理",
			API: "/api/w8t/dutyManage/dutyManageCreate",
		},
		"dutyManageDelete": {
			Key: "删除值班管理",
			API: "/api/w8t/dutyManage/dutyManageDelete",
		},
		"dutyManageList": {
			Key: "查看值班管理列表",
			API: "/api/w8t/dutyManage/dutyManageList",
		},
		"dutyManageUpdate": {
			Key: "更新值班管理",
			API: "/api/w8t/dutyManage/dutyManageUpdate",
		},
		"noticeCreate": {
			Key: "创建通知对象",
			API: "/api/w8t/notice/noticeCreate",
		},
		"noticeDelete": {
			Key: "删除通知对象",
			API: "/api/w8t/notice/noticeDelete",
		},
		"noticeList": {
			Key: "查看通知对象列表",
			API: "/api/w8t/notice/noticeList",
		},
		"noticeTemplateCreate": {
			Key: "创建通知模版",
			API: "/api/w8t/noticeTemplate/noticeTemplateCreate",
		},
		"noticeTemplateDelete": {
			Key: "删除通知模版",
			API: "/api/w8t/noticeTemplate/noticeTemplateDelete",
		},
		"noticeTemplateList": {
			Key: "查看通知模版列表",
			API: "/api/w8t/noticeTemplate/noticeTemplateList",
		},
		"noticeTemplateGet": {
			Key: "预览通知模版",
			API: "/api/w8t/noticeTemplate/noticeTemplateGet",
		},
		"noticeTemplateUpdate": {
			Key: "更新通知模版",
			API: "/api/w8t/noticeTemplate/noticeTemplateUpdate",
		},
		"noticeUpdate": {
			Key: "更新通知对象",
			API: "/api/w8t/notice/noticeUpdate",
		},
		"ruleCreate": {
			Key: "创建告警规则",
			API: "/api/w8t/rule/ruleCreate",
		},
		"ruleDelete": {
			Key: "删除告警规则",
			API: "/api/w8t/rule/ruleDelete",
		},
		"ruleGroupCreate": {
			Key: "创建告警规则组",
			API: "/api/w8t/ruleGroup/ruleGroupCreate",
		},
		"ruleGroupDelete": {
			Key: "删除告警规则组",
			API: "/api/w8t/ruleGroup/ruleGroupDelete",
		},
		"ruleGroupList": {
			Key: "查看告警规则组列表",
			API: "/api/w8t/ruleGroup/ruleGroupList",
		},
		"ruleGroupUpdate": {
			Key: "更新告警规则组",
			API: "/api/w8t/ruleGroup/ruleGroupUpdate",
		},
		"ruleList": {
			Key: "查看告警规则列表",
			API: "/api/w8t/rule/ruleList",
		},
		"ruleTmplCreate": {
			Key: "创建规则模版",
			API: "/api/w8t/ruleTmpl/ruleTmplCreate",
		},
		"ruleTmplUpdate": {
			Key: "更新规则模版",
			API: "/api/w8t/ruleTmpl/ruleTmplUpdate",
		},
		"ruleTmplDelete": {
			Key: "删除规则模版",
			API: "/api/w8t/ruleTmpl/ruleTmplDelete",
		},
		"ruleTmplGroupCreate": {
			Key: "创建规则模版组",
			API: "/api/w8t/ruleTmplGroup/ruleTmplGroupCreate",
		},
		"ruleTmplGroupUpdate": {
			Key: "更新规则模版组",
			API: "/api/w8t/ruleTmplGroup/ruleTmplGroupUpdate",
		},
		"ruleTmplGroupDelete": {
			Key: "删除规则模版组",
			API: "/api/w8t/ruleTmplGroup/ruleTmplGroupDelete",
		},
		"ruleTmplGroupList": {
			Key: "查看规则模版组列表",
			API: "/api/w8t/ruleTmplGroup/ruleTmplGroupList",
		},
		"ruleTmplList": {
			Key: "查看规则模版列表",
			API: "/api/w8t/ruleTmpl/ruleTmplList",
		},
		"ruleUpdate": {
			Key: "更新告警规则",
			API: "/api/w8t/rule/ruleUpdate",
		},
		"silenceCreate": {
			Key: "创建静默规则",
			API: "/api/w8t/silence/silenceCreate",
		},
		"silenceDelete": {
			Key: "删除静默规则",
			API: "/api/w8t/silence/silenceDelete",
		},
		"silenceList": {
			Key: "查看静默规则列表",
			API: "/api/w8t/silence/silenceList",
		},
		"silenceUpdate": {
			Key: "更新静默规则",
			API: "/api/w8t/silence/silenceUpdate",
		},
		"createProbing": {
			Key: "创建拨测规则",
			API: "/api/w8t/probing/createProbing",
		},
		"updateProbing": {
			Key: "更新拨测规则",
			API: "/api/w8t/probing/updateProbing",
		},
		"deleteProbing": {
			Key: "删除拨测规则",
			API: "/api/w8t/probing/deleteProbing",
		},
		"listProbing": {
			Key: "查看拨测规则列表",
			API: "/api/w8t/probing/listProbing",
		},
		"searchProbing": {
			Key: "预览拨测规则详情",
			API: "/api/w8t/probing/searchProbing",
		},
		"listFolder": {
			Key: "查看仪表盘目录列表",
			API: "/api/w8t/dashboard/listFolder",
		},
		"getFolder": {
			Key: "预览仪表盘目录详情",
			API: "/api/w8t/dashboard/getFolder",
		},
		"createFolder": {
			Key: "创建仪表盘目录",
			API: "/api/w8t/dashboard/createFolder",
		},
		"updateFolder": {
			Key: "更新仪表盘目录",
			API: "/api/w8t/dashboard/updateFolder",
		},
		"deleteFolder": {
			Key: "删除仪表盘目录",
			API: "/api/w8t/dashboard/deleteFolder",
		},
		"listGrafanaDashboards": {
			Key: "查看仪表盘列表",
			API: "/api/w8t/dashboard/listGrafanaDashboards",
		},
		"noticeRecordList": {
			Key: "查看通知记录列表",
			API: "/api/w8t/notice/noticeRecordList",
		},
		"faultCenterList": {
			Key: "查看故障中心列表",
			API: "/api/w8t/faultCenter/faultCenterList",
		},
		"faultCenterSearch": {
			Key: "预览故障中心",
			API: "/api/w8t/faultCenter/faultCenterSearch",
		},
		"faultCenterCreate": {
			Key: "创建故障中心",
			API: "/api/w8t/faultCenter/faultCenterCreate",
		},
		"faultCenterUpdate": {
			Key: "更新故障中心",
			API: "/api/w8t/faultCenter/faultCenterUpdate",
		},
		"faultCenterDelete": {
			Key: "删除故障中心",
			API: "/api/w8t/faultCenter/faultCenterDelete",
		},
		"faultCenterReset": {
			Key: "修改故障中心基本信息",
			API: "/api/w8t/faultCenter/faultCenterReset",
		},
		"listComments": {
			Key: "查看评论列表",
			API: "/api/w8t/event/listComments",
		},
		"addComment": {
			Key: "添加评论",
			API: "/api/w8t/event/addComment",
		},
		"deleteComment": {
			Key: "删除评论",
			API: "/api/w8t/event/deleteComment",
		},
		"recordingRuleCreate": {
			Key: "创建记录规则",
			API: "/api/w8t/recordingRule/recordingRuleCreate",
		},
		"recordingRuleUpdate": {
			Key: "更新记录规则",
			API: "/api/w8t/recordingRule/recordingRuleUpdate",
		},
		"recordingRuleDelete": {
			Key: "删除记录规则",
			API: "/api/w8t/recordingRule/recordingRuleDelete",
		},
		"userList": {
			Key: "查看用户列表",
			API: "/api/w8t/user/userList",
		},
	}
}
