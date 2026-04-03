package api

import (
	"watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"

	"github.com/gin-gonic/gin"
)

type settingsController struct{}

var SettingsController = new(settingsController)

func (settingsController settingsController) API(gin *gin.RouterGroup) {
	a := gin.Group("setting")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.AuditingLog(),
	)
	{
		a.POST("saveSystemSetting", settingsController.Save)
	}

	b := gin.Group("setting")
	b.Use(
		middleware.Auth(),
		middleware.Permission(),
	)
	{
		b.POST("syncLdapUser", settingsController.SyncLdapUser)
		b.GET("getSystemSetting", settingsController.Get)
	}
}

func (settingsController settingsController) Save(ctx *gin.Context) {
	r := new(models.Settings)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.SettingService.Save(r)
	})
}

func (settingsController settingsController) Get(ctx *gin.Context) {
	Service(ctx, func() (interface{}, interface{}) {
		return services.SettingService.Get()
	})
}

func (settingsController settingsController) SyncLdapUser(ctx *gin.Context) {
	Service(ctx, func() (interface{}, interface{}) {
		return nil, services.LdapService.SyncNow()
	})
}
