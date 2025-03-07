package v1

import (
	"watchAlert/api"
)

var (
	Notice         = api.ApiGroupApp.NoticeController
	Silence        = api.ApiGroupApp.SilenceController
	Datasource     = api.ApiGroupApp.DatasourceController
	Duty           = api.ApiGroupApp.DutyController
	DutyCalendar   = api.ApiGroupApp.DutyCalendarController
	Rule           = api.ApiGroupApp.RuleController
	Auth           = api.ApiGroupApp.UserController
	AlertEvent     = api.ApiGroupApp.AlertEventController
	Role           = api.ApiGroupApp.UserRoleController
	Permissions    = api.ApiGroupApp.UserPermissionsController
	NoticeTemplate = api.ApiGroupApp.NoticeTemplateController
	RuleGroup      = api.ApiGroupApp.RuleGroupController
	RuleTmplGroup  = api.ApiGroupApp.RuleTmplGroupController
	RuleTmpl       = api.ApiGroupApp.RuleTmplController
	DashboardInfo  = api.ApiGroupApp.DashboardInfoController
	Tenant         = api.ApiGroupApp.TenantController
	Dashboard      = api.ApiGroupApp.DashboardController
	AuditLog       = api.ApiGroupApp.AuditLogController
	ClientApi      = api.ApiGroupApp.ClientController
	AWSCloudWatch  = api.ApiGroupApp.AWSCloudWatchController
	AWSRds         = api.ApiGroupApp.AWSCloudWatchRDSController
	Setting        = api.ApiGroupApp.SettingsController
	KubeEvent      = api.ApiGroupApp.KubernetesTypesController
	Subscribe      = api.ApiGroupApp.SubscribeController
	Probing        = api.ApiGroupApp.ProbingController
	FaultCenter    = api.ApiGroupApp.FaultCenterController
	Ai             = api.ApiGroupApp.AiController
)
