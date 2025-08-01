package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
)

type dutyCalendarController struct{}

var DutyCalendarController = new(dutyCalendarController)

/*
值班表 API
/api/w8t/calendar
*/
func (dutyCalendarController dutyCalendarController) API(gin *gin.RouterGroup) {
	a := gin.Group("calendar")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		a.POST("calendarCreate", dutyCalendarController.Create)
		a.POST("calendarUpdate", dutyCalendarController.Update)
	}

	b := gin.Group("calendar")
	b.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		b.GET("calendarSearch", dutyCalendarController.Search)
	}

	c := gin.Group("calendar")
	c.Use(
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		c.GET("getCalendarUsers", dutyCalendarController.GetCalendarUsers)
	}
}

func (dutyCalendarController dutyCalendarController) Create(ctx *gin.Context) {
	r := new(types.RequestDutyCalendarCreate)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DutyCalendarService.CreateAndUpdate(r)
	})
}

func (dutyCalendarController dutyCalendarController) Update(ctx *gin.Context) {
	r := new(types.RequestDutyCalendarUpdate)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DutyCalendarService.Update(r)
	})
}

func (dutyCalendarController dutyCalendarController) Search(ctx *gin.Context) {
	r := new(types.RequestDutyCalendarQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DutyCalendarService.Search(r)
	})
}

func (dutyCalendarController dutyCalendarController) GetCalendarUsers(ctx *gin.Context) {
	r := new(types.RequestDutyCalendarQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DutyCalendarService.GetCalendarUsers(r)
	})
}
