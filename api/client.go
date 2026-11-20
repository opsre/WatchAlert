package api

import (
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"

	"github.com/gin-gonic/gin"
)

type clientController struct{}

var ClientController = new(clientController)

func (clientController clientController) API(gin *gin.RouterGroup) {
	a := gin.Group("c")
	a.Use(
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		a.GET("getJaegerService", clientController.GetJaegerService)
		a.GET("getTraceService", clientController.GetTraceService)
	}
}

func (clientController clientController) GetJaegerService(ctx *gin.Context) {
	r := new(types.RequestDatasourceQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.ClientService.GetJaegerService(r)
	})
}

func (clientController clientController) GetTraceService(ctx *gin.Context) {
	r := new(types.RequestDatasourceQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.ClientService.GetTraceService(r)
	})
}
