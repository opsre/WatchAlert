package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
)

type auditLogController struct{}

var AuditLogController = new(auditLogController)

func (auditLogController auditLogController) API(gin *gin.RouterGroup) {
	a := gin.Group("auditLog")
	a.Use(
		middleware.Cors(),
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		a.GET("listAuditLog", auditLogController.List)
		a.GET("searchAuditLog", auditLogController.Search)
	}
}

func (auditLogController auditLogController) List(ctx *gin.Context) {
	r := new(types.RequestAuditLogQuery)
	BindQuery(ctx, r)
	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)
	Service(ctx, func() (interface{}, interface{}) {
		return services.AuditLogService.List(r)
	})
}

func (auditLogController auditLogController) Search(ctx *gin.Context) {
	r := new(types.RequestAuditLogQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.AuditLogService.Search(r)
	})
}
