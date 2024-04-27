INSERT INTO `user_roles` (`id`, `name`, `description`, `permissions`, `create_at`) VALUES ('ur-colpotjadq770urjfrhg', 'admin', 'system', '[{\"key\":\"更新告警规则组\",\"api\":\"/api/w8t/ruleGroup/ruleGroupUpdate\"},{\"key\":\"搜索值班用户\",\"api\":\"/api/w8t/user/searchDutyUser\"},{\"key\":\"修改用户密码\",\"api\":\"/api/w8t/user/userChangePass\"},{\"key\":\"获取Jaeger服务列表\",\"api\":\"/api/w8t/c/getJaegerService\"},{\"key\":\"删除数据源\",\"api\":\"/api/w8t/datasource/dataSourceDelete\"},{\"key\":\"查看历史告警\",\"api\":\"/api/w8t/event/hisEvent\"},{\"key\":\"创建规则模版\",\"api\":\"/api/w8t/ruleTmpl/ruleTmplCreate\"},{\"key\":\"删除规则模版\",\"api\":\"/api/w8t/ruleTmpl/ruleTmplDelete\"},{\"key\":\"查看静默规则\",\"api\":\"/api/w8t/silence/silenceList\"},{\"key\":\"搜索用户\",\"api\":\"/api/w8t/user/searchUser\"},{\"key\":\"更新日历表\",\"api\":\"/api/w8t/calendar/calendarUpdate\"},{\"key\":\"更新通知模版\",\"api\":\"/api/w8t/noticeTemplate/noticeTemplateUpdate\"},{\"key\":\"更新仪表盘\",\"api\":\"/api/w8t/dashboard/updateDashboard\"},{\"key\":\"更新告警规则\",\"api\":\"/api/w8t/rule/ruleUpdate\"},{\"key\":\"更新静默规则\",\"api\":\"/api/w8t/silence/silenceUpdate\"},{\"key\":\"查看告警规则组\",\"api\":\"/api/w8t/ruleGroup/ruleGroupList\"},{\"key\":\"查看当前告警事件\",\"api\":\"/api/w8t/event/curEvent\"},{\"key\":\"创建值班表\",\"api\":\"/api/w8t/dutyManage/dutyManageCreate\"},{\"key\":\"获取仪表盘\",\"api\":\"/api/w8t/dashboard/getDashboard\"},{\"key\":\"搜索日历表\",\"api\":\"/api/w8t/calendar/calendarSearch\"},{\"key\":\"查看数据源\",\"api\":\"/api/w8t/datasource/dataSourceList\"},{\"key\":\"搜索仪表盘\",\"api\":\"/api/w8t/dashboard/searchDashboard\"},{\"key\":\"创建静默规则\",\"api\":\"/api/w8t/silence/silenceCreate\"},{\"key\":\"删除静默规则\",\"api\":\"/api/w8t/silence/silenceDelete\"},{\"key\":\"更新通知对象\",\"api\":\"/api/w8t/sender/noticeUpdate\"},{\"key\":\"用户注册\",\"api\":\"/api/system/register\"},{\"key\":\"搜索通知模版\",\"api\":\"/api/w8t/noticeTemplate/searchNoticeTmpl\"},{\"key\":\"创建数据源\",\"api\":\"/api/w8t/datasource/dataSourceCreate\"},{\"key\":\"搜索通知对象\",\"api\":\"/api/w8t/sender/noticeSearch\"},{\"key\":\"更新用户信息\",\"api\":\"/api/w8t/user/userUpdate\"},{\"key\":\"创建租户\",\"api\":\"/api/w8t/tenant/createTenant\"},{\"key\":\"删除通知对象\",\"api\":\"/api/w8t/notice/noticeDelete\"},{\"key\":\"更新值班表\",\"api\":\"/api/w8t/dutyManage/dutyManageDelete\"},{\"key\":\"查看值班表\",\"api\":\"/api/w8t/dutyManage/dutyManageList\"},{\"key\":\"删除规则模版组\",\"api\":\"/api/w8t/ruleTmplGroup/ruleTmplGroupDelete\"},{\"key\":\"查看用户列表\",\"api\":\"/api/w8t/user/userList\"},{\"key\":\"发布日历表\",\"api\":\"/api/w8t/calendar/calendarCreate\"},{\"key\":\"获取数据源\",\"api\":\"/api/w8t/datasource/dataSourceGet\"},{\"key\":\"创建通知模版\",\"api\":\"/api/w8t/noticeTemplate/noticeTemplateCreate\"},{\"key\":\"删除通知模版\",\"api\":\"/api/w8t/noticeTemplate/noticeTemplateDelete\"},{\"key\":\"更新租户信息\",\"api\":\"/api/w8t/tenant/updateTenant\"},{\"key\":\"删除用户\",\"api\":\"/api/w8t/user/userDelete\"},{\"key\":\"搜索数据源\",\"api\":\"/api/w8t/datasource/dataSourceSearch\"},{\"key\":\"创建通知对象\",\"api\":\"/api/w8t/notice/noticeCreate\"},{\"key\":\"查看用户角色\",\"api\":\"/api/w8t/role/roleList\"},{\"key\":\"查看告警规则\",\"api\":\"/api/w8t/rule/ruleList\"},{\"key\":\"创建仪表盘\",\"api\":\"/api/w8t/dashboard/createDashboard\"},{\"key\":\"更新数据源\",\"api\":\"/api/w8t/datasource/dataSourceUpdate\"},{\"key\":\"更新用户角色\",\"api\":\"/api/w8t/role/roleUpdate\"},{\"key\":\"删除告警规则\",\"api\":\"/api/w8t/rule/ruleDelete\"},{\"key\":\"查看通知对象\",\"api\":\"/api/w8t/notice/noticeList\"},{\"key\":\"查看通知模版\",\"api\":\"/api/w8t/noticeTemplate/noticeTemplateList\"},{\"key\":\"更新值班表\",\"api\":\"/api/w8t/dutyManage/dutyManageUpdate\"},{\"key\":\"创建告警规则组\",\"api\":\"/api/w8t/ruleGroup/ruleGroupCreate\"},{\"key\":\"创建规则模版组\",\"api\":\"/api/w8t/ruleTmplGroup/ruleTmplGroupCreate\"},{\"key\":\"查看规则模版组\",\"api\":\"/api/w8t/ruleTmplGroup/ruleTmplGroupList\"},{\"key\":\"搜索告警规则\",\"api\":\"/api/w8t/rule/ruleSearch\"},{\"key\":\"删除租户\",\"api\":\"/api/w8t/tenant/deleteTenant\"},{\"key\":\"创建用户角色\",\"api\":\"/api/w8t/role/roleCreate\"},{\"key\":\"删除用户角色\",\"api\":\"/api/w8t/role/roleDelete\"},{\"key\":\"创建告警规则\",\"api\":\"/api/w8t/rule/ruleCreate\"},{\"key\":\"查看仪表盘\",\"api\":\"/api/w8t/dashboard/listDashboard\"},{\"key\":\"查看用户权限\",\"api\":\"/api/w8t/permissions/permsList\"},{\"key\":\"查看租户\",\"api\":\"/api/w8t/tenant/getTenantList\"},{\"key\":\"删除告警规则组\",\"api\":\"/api/w8t/ruleGroup/ruleGroupDelete\"},{\"key\":\"查看规则模版\",\"api\":\"/api/w8t/ruleTmpl/ruleTmplList\"},{\"key\":\"删除仪表盘\",\"api\":\"/api/w8t/dashboard/deleteDashboard\"},{\"key\":\"搜索值班表\",\"api\":\"/api/w8t/dutyManage/dutyManageSearch\"}]', 1714134134);
INSERT INTO `user_roles` (`id`, `name`, `description`, `permissions`, `create_at`) VALUES ('ur-colpq8jadq771uj7j7b0', '只读', 'Read', '[{\"key\":\"获取Jaeger服务列表\",\"api\":\"/api/w8t/c/getJaegerService\"},{\"key\":\"获取数据源\",\"api\":\"/api/w8t/datasource/dataSourceGet\"},{\"key\":\"获取仪表盘\",\"api\":\"/api/w8t/dashboard/getDashboard\"},{\"key\":\"搜索告警规则\",\"api\":\"/api/w8t/rule/ruleSearch\"},{\"key\":\"搜索通知对象\",\"api\":\"/api/w8t/sender/noticeSearch\"},{\"key\":\"搜索仪表盘\",\"api\":\"/api/w8t/dashboard/searchDashboard\"},{\"key\":\"搜索值班用户\",\"api\":\"/api/w8t/user/searchDutyUser\"},{\"key\":\"搜索数据源\",\"api\":\"/api/w8t/datasource/dataSourceSearch\"},{\"key\":\"搜索值班表\",\"api\":\"/api/w8t/dutyManage/dutyManageSearch\"},{\"key\":\"搜索用户\",\"api\":\"/api/w8t/user/searchUser\"},{\"key\":\"搜索通知模版\",\"api\":\"/api/w8t/noticeTemplate/searchNoticeTmpl\"},{\"key\":\"搜索日历表\",\"api\":\"/api/w8t/calendar/calendarSearch\"},{\"key\":\"查看租户\",\"api\":\"/api/w8t/tenant/getTenantList\"},{\"key\":\"查看仪表盘\",\"api\":\"/api/w8t/dashboard/listDashboard\"},{\"key\":\"查看用户权限\",\"api\":\"/api/w8t/permissions/permsList\"},{\"key\":\"查看通知对象\",\"api\":\"/api/w8t/notice/noticeList\"},{\"key\":\"查看值班表\",\"api\":\"/api/w8t/dutyManage/dutyManageList\"},{\"key\":\"查看用户角色\",\"api\":\"/api/w8t/role/roleList\"},{\"key\":\"查看告警规则\",\"api\":\"/api/w8t/rule/ruleList\"},{\"key\":\"查看数据源\",\"api\":\"/api/w8t/datasource/dataSourceList\"},{\"key\":\"查看静默规则\",\"api\":\"/api/w8t/silence/silenceList\"},{\"key\":\"查看用户列表\",\"api\":\"/api/w8t/user/userList\"},{\"key\":\"查看历史告警\",\"api\":\"/api/w8t/event/hisEvent\"},{\"key\":\"查看通知模版\",\"api\":\"/api/w8t/noticeTemplate/noticeTemplateList\"},{\"key\":\"查看规则模版\",\"api\":\"/api/w8t/ruleTmpl/ruleTmplList\"},{\"key\":\"查看当前告警事件\",\"api\":\"/api/w8t/event/curEvent\"},{\"key\":\"查看告警规则组\",\"api\":\"/api/w8t/ruleGroup/ruleGroupList\"},{\"key\":\"查看规则模版组\",\"api\":\"/api/w8t/ruleTmplGroup/ruleTmplGroupList\"}]', 1714134306);
INSERT INTO `user_roles` (`id`, `name`, `description`, `permissions`, `create_at`) VALUES ('ur-com7nvjadq7bufnlg130', '读写', 'ReadWrite', '[{\"key\":\"更新通知对象\",\"api\":\"/api/w8t/sender/noticeUpdate\"},{\"key\":\"更新告警规则\",\"api\":\"/api/w8t/rule/ruleUpdate\"},{\"key\":\"更新日历表\",\"api\":\"/api/w8t/calendar/calendarUpdate\"},{\"key\":\"更新告警规则组\",\"api\":\"/api/w8t/ruleGroup/ruleGroupUpdate\"},{\"key\":\"更新静默规则\",\"api\":\"/api/w8t/silence/silenceUpdate\"},{\"key\":\"更新用户角色\",\"api\":\"/api/w8t/role/roleUpdate\"},{\"key\":\"更新值班表\",\"api\":\"/api/w8t/dutyManage/dutyManageUpdate\"},{\"key\":\"更新租户信息\",\"api\":\"/api/w8t/tenant/updateTenant\"},{\"key\":\"更新仪表盘\",\"api\":\"/api/w8t/dashboard/updateDashboard\"},{\"key\":\"更新数据源\",\"api\":\"/api/w8t/datasource/dataSourceUpdate\"},{\"key\":\"更新通知模版\",\"api\":\"/api/w8t/noticeTemplate/noticeTemplateUpdate\"},{\"key\":\"更新用户信息\",\"api\":\"/api/w8t/user/userUpdate\"},{\"key\":\"获取仪表盘\",\"api\":\"/api/w8t/dashboard/getDashboard\"},{\"key\":\"查看仪表盘\",\"api\":\"/api/w8t/dashboard/listDashboard\"},{\"key\":\"查看通知模版\",\"api\":\"/api/w8t/noticeTemplate/noticeTemplateList\"},{\"key\":\"查看规则模版\",\"api\":\"/api/w8t/ruleTmpl/ruleTmplList\"},{\"key\":\"查看数据源\",\"api\":\"/api/w8t/datasource/dataSourceList\"},{\"key\":\"查看用户权限\",\"api\":\"/api/w8t/permissions/permsList\"},{\"key\":\"查看告警规则\",\"api\":\"/api/w8t/rule/ruleList\"},{\"key\":\"查看值班表\",\"api\":\"/api/w8t/dutyManage/dutyManageList\"},{\"key\":\"查看用户列表\",\"api\":\"/api/w8t/user/userList\"},{\"key\":\"查看静默规则\",\"api\":\"/api/w8t/silence/silenceList\"},{\"key\":\"查看当前告警事件\",\"api\":\"/api/w8t/event/curEvent\"},{\"key\":\"查看用户角色\",\"api\":\"/api/w8t/role/roleList\"},{\"key\":\"搜索仪表盘\",\"api\":\"/api/w8t/dashboard/searchDashboard\"},{\"key\":\"搜索值班表\",\"api\":\"/api/w8t/dutyManage/dutyManageSearch\"},{\"key\":\"获取Jaeger服务列表\",\"api\":\"/api/w8t/c/getJaegerService\"},{\"key\":\"搜索通知模版\",\"api\":\"/api/w8t/noticeTemplate/searchNoticeTmpl\"},{\"key\":\"搜索告警规则\",\"api\":\"/api/w8t/rule/ruleSearch\"},{\"key\":\"查看告警规则组\",\"api\":\"/api/w8t/ruleGroup/ruleGroupList\"},{\"key\":\"搜索数据源\",\"api\":\"/api/w8t/datasource/dataSourceSearch\"},{\"key\":\"用户注册\",\"api\":\"/api/system/register\"},{\"key\":\"搜索用户\",\"api\":\"/api/w8t/user/searchUser\"},{\"key\":\"查看历史告警\",\"api\":\"/api/w8t/event/hisEvent\"},{\"key\":\"发布日历表\",\"api\":\"/api/w8t/calendar/calendarCreate\"},{\"key\":\"查看规则模版组\",\"api\":\"/api/w8t/ruleTmplGroup/ruleTmplGroupList\"},{\"key\":\"搜索日历表\",\"api\":\"/api/w8t/calendar/calendarSearch\"},{\"key\":\"获取数据源\",\"api\":\"/api/w8t/datasource/dataSourceGet\"},{\"key\":\"查看租户\",\"api\":\"/api/w8t/tenant/getTenantList\"},{\"key\":\"查看通知对象\",\"api\":\"/api/w8t/notice/noticeList\"},{\"key\":\"搜索通知对象\",\"api\":\"/api/w8t/sender/noticeSearch\"},{\"key\":\"搜索值班用户\",\"api\":\"/api/w8t/user/searchDutyUser\"},{\"key\":\"修改用户密码\",\"api\":\"/api/w8t/user/userChangePass\"},{\"key\":\"创建通知对象\",\"api\":\"/api/w8t/notice/noticeCreate\"},{\"key\":\"创建规则模版\",\"api\":\"/api/w8t/ruleTmpl/ruleTmplCreate\"},{\"key\":\"创建静默规则\",\"api\":\"/api/w8t/silence/silenceCreate\"},{\"key\":\"创建告警规则\",\"api\":\"/api/w8t/rule/ruleCreate\"},{\"key\":\"创建仪表盘\",\"api\":\"/api/w8t/dashboard/createDashboard\"},{\"key\":\"创建规则模版组\",\"api\":\"/api/w8t/ruleTmplGroup/ruleTmplGroupCreate\"},{\"key\":\"创建租户\",\"api\":\"/api/w8t/tenant/createTenant\"},{\"key\":\"创建数据源\",\"api\":\"/api/w8t/datasource/dataSourceCreate\"},{\"key\":\"创建通知模版\",\"api\":\"/api/w8t/noticeTemplate/noticeTemplateCreate\"},{\"key\":\"创建告警规则组\",\"api\":\"/api/w8t/ruleGroup/ruleGroupCreate\"},{\"key\":\"创建值班表\",\"api\":\"/api/w8t/dutyManage/dutyManageCreate\"},{\"key\":\"创建用户角色\",\"api\":\"/api/w8t/role/roleCreate\"}]', 1714191358);

