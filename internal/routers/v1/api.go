package v1

import (
	"github.com/gin-gonic/gin"
	"watchAlert/api"
)

func Router(engine *gin.Engine) {
	v1 := engine.Group("api")
	{
		system := v1.Group("system")
		{
			api.DashboardInfoController.API(v1)
			system.POST("register", api.UserController.Register)
			system.POST("login", api.UserController.Login)
			system.GET("checkUser", api.UserController.CheckUser)
			system.GET("checkNoticeStatus", api.NoticeController.Check)
			system.GET("userInfo", api.UserController.GetUserInfo)
		}

		w8t := v1.Group("w8t")
		{
			api.UserController.API(w8t)
			api.UserPermissionsController.API(w8t)
			api.AlertEventController.API(w8t)
			api.UserRoleController.API(w8t)
			api.DashboardController.API(w8t)
			api.DatasourceController.API(w8t)
			api.RuleGroupController.API(w8t)
			api.RuleController.API(w8t)
			api.SilenceController.API(w8t)
			api.NoticeController.API(w8t)
			api.NoticeTemplateController.API(w8t)
			api.TenantController.API(w8t)
			api.RuleTmplGroupController.API(w8t)
			api.RuleTmplController.API(w8t)
			api.DutyController.API(w8t)
			api.DutyCalendarController.API(w8t)
			api.AuditLogController.API(w8t)
			api.ClientController.API(w8t)
			api.AWSCloudWatchController.API(w8t)
			api.AWSCloudWatchRDSController.API(w8t)
			api.SettingsController.API(w8t)
			api.KubernetesTypesController.API(w8t)
			api.SubscribeController.API(w8t)
			api.ProbingController.API(w8t)
			api.FaultCenterController.API(w8t)
			api.AiController.API(w8t)
		}
	}
}
