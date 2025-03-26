package models

type UserPermissions struct {
	Key string `json:"key"`
	API string `json:"api"`
}

func PermissionsInfo() map[string]UserPermissions {
	return map[string]UserPermissions{
		"ruleSearch": {
			Key: "搜索告警规则",
			API: "/api/w8t/rule/ruleSearch",
		},
		"calendarCreate": {
			Key: "发布日历表",
			API: "/api/w8t/calendar/calendarCreate",
		},
		"calendarSearch": {
			Key: "搜索日历表",
			API: "/api/w8t/calendar/calendarSearch",
		},
		"getCalendarUsers": {
			Key: "获取值班表用户列表",
			API: "/api/w8t/calendar/getCalendarUsers",
		},
		"calendarUpdate": {
			Key: "更新日历表",
			API: "/api/w8t/calendar/calendarUpdate",
		},
		"createDashboard": {
			Key: "创建仪表盘",
			API: "/api/w8t/dashboard/createDashboard",
		},
		"createTenant": {
			Key: "创建租户",
			API: "/api/w8t/tenant/createTenant",
		},
		"curEvent": {
			Key: "查看当前告警事件",
			API: "/api/w8t/event/curEvent",
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
			Key: "获取数据源",
			API: "/api/w8t/datasource/dataSourceGet",
		},
		"dataSourceList": {
			Key: "查看数据源",
			API: "/api/w8t/datasource/dataSourceList",
		},
		"dataSourceSearch": {
			Key: "搜索数据源",
			API: "/api/w8t/datasource/dataSourceSearch",
		},
		"dataSourceUpdate": {
			Key: "更新数据源",
			API: "/api/w8t/datasource/dataSourceUpdate",
		},
		"deleteDashboard": {
			Key: "删除仪表盘",
			API: "/api/w8t/dashboard/deleteDashboard",
		},
		"deleteTenant": {
			Key: "删除租户",
			API: "/api/w8t/tenant/deleteTenant",
		},
		"dutyManageCreate": {
			Key: "创建值班表",
			API: "/api/w8t/dutyManage/dutyManageCreate",
		},
		"dutyManageDelete": {
			Key: "删除值班表",
			API: "/api/w8t/dutyManage/dutyManageDelete",
		},
		"dutyManageList": {
			Key: "查看值班表",
			API: "/api/w8t/dutyManage/dutyManageList",
		},
		"dutyManageSearch": {
			Key: "搜索值班表",
			API: "/api/w8t/dutyManage/dutyManageSearch",
		},
		"dutyManageUpdate": {
			Key: "更新值班表",
			API: "/api/w8t/dutyManage/dutyManageUpdate",
		},
		"getDashboard": {
			Key: "获取仪表盘",
			API: "/api/w8t/dashboard/getDashboard",
		},
		"getTenantList": {
			Key: "查看租户",
			API: "/api/w8t/tenant/getTenantList",
		},
		"hisEvent": {
			Key: "查看历史告警",
			API: "/api/w8t/event/hisEvent",
		},
		"listDashboard": {
			Key: "查看仪表盘",
			API: "/api/w8t/dashboard/listDashboard",
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
			Key: "查看通知对象",
			API: "/api/w8t/notice/noticeList",
		},
		"noticeSearch": {
			Key: "搜索通知对象",
			API: "/api/w8t/notice/noticeSearch",
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
			Key: "查看通知模版",
			API: "/api/w8t/noticeTemplate/noticeTemplateList",
		},
		"noticeTemplateUpdate": {
			Key: "更新通知模版",
			API: "/api/w8t/noticeTemplate/noticeTemplateUpdate",
		},
		"noticeUpdate": {
			Key: "更新通知对象",
			API: "/api/w8t/notice/noticeUpdate",
		},
		"permsList": {
			Key: "查看用户权限",
			API: "/api/w8t/permissions/permsList",
		},
		"register": {
			Key: "用户注册",
			API: "/api/system/register",
		},
		"roleCreate": {
			Key: "创建用户角色",
			API: "/api/w8t/role/roleCreate",
		},
		"roleDelete": {
			Key: "删除用户角色",
			API: "/api/w8t/role/roleDelete",
		},
		"roleList": {
			Key: "查看用户角色",
			API: "/api/w8t/role/roleList",
		},
		"roleUpdate": {
			Key: "更新用户角色",
			API: "/api/w8t/role/roleUpdate",
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
			Key: "查看告警规则组",
			API: "/api/w8t/ruleGroup/ruleGroupList",
		},
		"ruleGroupUpdate": {
			Key: "更新告警规则组",
			API: "/api/w8t/ruleGroup/ruleGroupUpdate",
		},
		"ruleList": {
			Key: "查看告警规则",
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
			Key: "查看规则模版组",
			API: "/api/w8t/ruleTmplGroup/ruleTmplGroupList",
		},
		"ruleTmplList": {
			Key: "查看规则模版",
			API: "/api/w8t/ruleTmpl/ruleTmplList",
		},
		"ruleUpdate": {
			Key: "更新告警规则",
			API: "/api/w8t/rule/ruleUpdate",
		},
		"searchDashboard": {
			Key: "搜索仪表盘",
			API: "/api/w8t/dashboard/searchDashboard",
		},
		"searchDutyUser": {
			Key: "搜索值班用户",
			API: "/api/w8t/user/searchDutyUser",
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
			Key: "查看静默规则",
			API: "/api/w8t/silence/silenceList",
		},
		"silenceUpdate": {
			Key: "更新静默规则",
			API: "/api/w8t/silence/silenceUpdate",
		},
		"updateDashboard": {
			Key: "更新仪表盘",
			API: "/api/w8t/dashboard/updateDashboard",
		},
		"updateTenant": {
			Key: "更新租户信息",
			API: "/api/w8t/tenant/updateTenant",
		},
		"userChangePass": {
			Key: "修改用户密码",
			API: "/api/w8t/user/userChangePass",
		},
		"userDelete": {
			Key: "删除用户",
			API: "/api/w8t/user/userDelete",
		},
		"userList": {
			Key: "查看用户列表",
			API: "/api/w8t/user/userList",
		},
		"userUpdate": {
			Key: "更新用户信息",
			API: "/api/w8t/user/userUpdate",
		},
		"getJaegerService": {
			Key: "获取Jaeger服务列表",
			API: "/api/w8t/c/getJaegerService",
		},
		"searchUser": {
			Key: "搜索用户",
			API: "/api/w8t/user/searchUser",
		},
		"searchNoticeTmpl": {
			Key: "搜索通知模版",
			API: "/api/w8t/noticeTemplate/searchNoticeTmpl",
		},
		"saveSystemSetting": {
			Key: "编辑系统配置",
			API: "/api/w8t/setting/saveSystemSetting",
		},
		"getSystemSetting": {
			Key: "获取系统配置",
			API: "/api/w8t/setting/getSystemSetting",
		},
		"promQuery": {
			Key: "Prometheus指标查询",
			API: "/api/w8t/datasource/promQuery",
		},
		"getTenant": {
			Key: "获取租户详细信息",
			API: "/api/w8t/tenant/getTenant",
		},
		"addUsersToTenant": {
			Key: "向租户添加成员",
			API: "/api/w8t/tenant/addUsersToTenant",
		},
		"delUsersOfTenant": {
			Key: "删除租户成员",
			API: "/api/w8t/tenant/delUsersOfTenant",
		},
		"getUsersForTenant": {
			Key: "获取租户成员列表",
			API: "/api/w8t/tenant/getUsersForTenant",
		},
		"changeTenantUserRole": {
			Key: "修改租户成员角色",
			API: "/api/w8t/tenant/changeTenantUserRole",
		},
		"getResourceList": {
			Key: "获取Kubernetes资源列表",
			API: "/api/w8t/kubernetes/getResourceList",
		},
		"getReasonList": {
			Key: "获取Kubernetes事件类型列表",
			API: "/api/w8t/kubernetes/getReasonList",
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
			Key: "获取拨测规则列表",
			API: "/api/w8t/probing/listProbing",
		},
		"searchProbing": {
			Key: "获取拨测规则信息",
			API: "/api/w8t/probing/searchProbing",
		},
		"onceProbing": {
			Key: "一次性拨测任务",
			API: "/api/w8t/probing/onceProbing",
		},
		"listFolder": {
			Key: "获取仪表盘目录列表",
			API: "/api/w8t/dashboard/listFolder",
		},
		"getFolder": {
			Key: "获取仪表盘目录详情",
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
			Key: "获取仪表盘图表列表",
			API: "/api/w8t/dashboard/listGrafanaDashboards",
		},
		"getDashboardFullUrl": {
			Key: "获取仪表盘完整URL",
			API: "/api/w8t/dashboard/getDashboardFullUrl",
		},
		"createSubscribe": {
			Key: "创建告警订阅",
			API: "/api/w8t/subscribe/createSubscribe",
		},
		"deleteSubscribe": {
			Key: "删除告警订阅",
			API: "/api/w8t/subscribe/deleteSubscribe",
		},
		"listSubscribe": {
			Key: "获取告警订阅",
			API: "/api/w8t/subscribe/listSubscribe",
		},
		"getSubscribe": {
			Key: "搜索告警订阅",
			API: "/api/w8t/subscribe/getSubscribe",
		},
		"noticeRecordList": {
			Key: "获取通知记录列表",
			API: "/api/w8t/notice/noticeRecordList",
		},
		"noticeRecordMetric": {
			Key: "获取通知记录指标",
			API: "/api/w8t/notice/noticeRecordMetric",
		},
		"dataSourcePing": {
			Key: "数据源连接测试",
			API: "/api/w8t/datasource/dataSourcePing",
		},
		"faultCenterList": {
			Key: "获取故障中心列表",
			API: "/api/w8t/faultCenter/faultCenterList",
		},
		"faultCenterSearch": {
			Key: "查询故障中心",
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
		"esSearch": {
			Key: "ElasticSearch搜索",
			API: "/api/w8t/datasource/esSearch",
		},
	}
}
