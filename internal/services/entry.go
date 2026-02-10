package services

import (
	"watchAlert/internal/ctx"
	"watchAlert/pkg/aws/cloudwatch"
	"watchAlert/pkg/aws/region"
)

var (
	DatasourceService       InterDatasourceService
	AuditLogService         InterAuditLogService
	DashboardService        InterDashboardService
	DutyManageService       InterDutyManageService
	DutyCalendarService     InterDutyCalendarService
	EventService            InterEventService
	NoticeService           InterNoticeService
	NoticeTmplService       InterNoticeTmplService
	RuleService             InterRuleService
	RuleGroupService        InterRuleGroupService
	RuleTmplService         InterRuleTmplService
	SilenceService          InterSilenceService
	TenantService           InterTenantService
	UserService             InterUserService
	UserRoleService         InterUserRoleService
	AlertService            InterAlertService
	RuleTmplGroupService    InterRuleTmplGroupService
	UserPermissionService   InterUserPermissionService
	AWSRegionService        region.InterAwsRegionService
	AWSCloudWatchService    cloudwatch.InterAwsCloudWatchService
	AWSCloudWatchRdsService cloudwatch.InterAwsRdsService
	SettingService          InterSettingService
	ClientService           InterClientService
	LdapService             InterLdapService
	SubscribeService        InterAlertSubscribeService
	ProbingService          InterProbingService
	FaultCenterService      InterFaultCenterService
	AiService               InterAiService
	OidcService             InterOidcService
	TopologyService         InterTopologyService
	ApiKeyService           InterApiKeyService
)

func NewServices(ctx *ctx.Context) {
	DatasourceService = newInterDatasourceService(ctx)
	AuditLogService = newInterAuditLogService(ctx)
	DashboardService = newInterDashboardService(ctx)
	DutyManageService = newInterDutyManageService(ctx)
	DutyCalendarService = newInterDutyCalendarService(ctx)
	EventService = newInterEventService(ctx)
	NoticeService = newInterAlertNoticeService(ctx)
	NoticeTmplService = newInterNoticeTmplService(ctx)
	RuleService = newInterRuleService(ctx)
	RuleGroupService = newInterRuleGroupService(ctx)
	RuleTmplService = newInterRuleTmplService(ctx)
	RuleTmplGroupService = newInterRuleTmplGroupService(ctx)
	SilenceService = newInterSilenceService(ctx)
	TenantService = newInterTenantService(ctx)
	UserService = newInterUserService(ctx)
	UserRoleService = newInterUserRoleService(ctx)
	AlertService = newInterAlertService(ctx)
	UserPermissionService = newInterUserPermissionService(ctx)
	AWSRegionService = region.NewInterAwsRegionService(ctx)
	AWSCloudWatchService = cloudwatch.NewInterAwsCloudWatchService(ctx)
	AWSCloudWatchRdsService = cloudwatch.NewInterAWSRdsService(ctx)
	SettingService = newInterSettingService(ctx)
	ClientService = newInterClientService(ctx)
	LdapService = newInterLdapService(ctx)
	SubscribeService = newInterAlertSubscribe(ctx)
	ProbingService = newInterProbingService(ctx)
	FaultCenterService = newInterFaultCenterService(ctx)
	AiService = newInterAiService(ctx)
	OidcService = newInterOidcService(ctx)
	TopologyService = newInterTopologyService(ctx)
	ApiKeyService = newInterApiKeyService(ctx)
}
