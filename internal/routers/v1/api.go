package v1

import (
	"github.com/gin-gonic/gin"
)

func Router(engine *gin.Engine) {

	v1 := engine.Group("api")
	{
		system := v1.Group("system")
		{
			DashboardInfo.API(v1)
			system.POST("register", Auth.Register)
			system.POST("login", Auth.Login)
			system.GET("checkUser", Auth.CheckUser)
			system.GET("checkNoticeStatus", Notice.Check)
			system.GET("userInfo", Auth.Get)
		}

		w8t := v1.Group("w8t")
		{
			Auth.API(w8t)
			Permissions.API(w8t)
			AlertEvent.API(w8t)
			Role.API(w8t)
			Dashboard.API(w8t)
			Datasource.API(w8t)
			RuleGroup.API(w8t)
			Rule.API(w8t)
			Silence.API(w8t)
			Notice.API(w8t)
			NoticeTemplate.API(w8t)
			Tenant.API(w8t)
			RuleTmplGroup.API(w8t)
			RuleTmpl.API(w8t)
			Duty.API(w8t)
			DutyCalendar.API(w8t)
			AuditLog.API(w8t)
			ClientApi.API(w8t)
			AWSCloudWatch.API(w8t)
			AWSRds.API(w8t)
			Setting.API(w8t)
			KubeEvent.API(w8t)
			Subscribe.API(w8t)
			Probing.API(w8t)
			FaultCenter.API(w8t)
		}

	}

}
