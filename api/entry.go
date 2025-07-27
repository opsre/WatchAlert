package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/pkg/response"
)

type ApiGroup struct {
	NoticeController
	DutyController
	DutyCalendarController
	DatasourceController
	SilenceController
	RuleController
	UserController
	AlertEventController
	UserRoleController
	UserPermissionsController
	NoticeTemplateController
	RuleGroupController
	RuleTmplGroupController
	RuleTmplController
	DashboardInfoController
	TenantController
	DashboardController
	AuditLogController
	ClientController
	AWSCloudWatchController
	AWSCloudWatchRDSController
	SettingsController
	KubernetesTypesController
	SubscribeController
	ProbingController
	FaultCenterController
	AiController
}

var ApiGroupApp = new(ApiGroup)

func Service(ctx *gin.Context, fu func() (interface{}, interface{})) {
	data, err := fu()
	if err != nil {
		response.Fail(ctx, err.(error).Error(), "failed")
		ctx.Abort()
		return
	} else {
		response.Success(ctx, data, "success")
	}
}

func BindJson(ctx *gin.Context, req interface{}) {
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		ctx.Abort()
		return
	}
}

func BindQuery(ctx *gin.Context, req interface{}) {
	err := ctx.ShouldBindQuery(req)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		ctx.Abort()
		return
	}
}
