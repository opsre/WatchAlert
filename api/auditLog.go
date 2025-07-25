package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
)

type AuditLogController struct{}

func (ac AuditLogController) API(gin *gin.RouterGroup) {
	auditLog := gin.Group("auditLog")
	auditLog.Use(
		middleware.Cors(),
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		auditLog.GET("listAuditLog", ac.List)
		auditLog.GET("searchAuditLog", ac.Search)
	}
}

func (ac AuditLogController) List(ctx *gin.Context) {
	r := new(types.RequestAuditLogQuery)
	BindQuery(ctx, r)
	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)
	Service(ctx, func() (interface{}, interface{}) {
		return services.AuditLogService.List(r)
	})
}

func (ac AuditLogController) Search(ctx *gin.Context) {
	r := new(types.RequestAuditLogQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.AuditLogService.Search(r)
	})
}
