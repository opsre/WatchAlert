package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
)

type DutyCalendarController struct{}

/*
值班表 API
/api/w8t/calendar
*/
func (dc DutyCalendarController) API(gin *gin.RouterGroup) {
	calendarA := gin.Group("calendar")
	calendarA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		calendarA.POST("calendarCreate", dc.Create)
		calendarA.POST("calendarUpdate", dc.Update)
	}

	calendarB := gin.Group("calendar")
	calendarB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		calendarB.GET("calendarSearch", dc.Search)
		calendarB.GET("getCalendarUsers", dc.GetCalendarUsers)
	}
}

func (dc DutyCalendarController) Create(ctx *gin.Context) {
	r := new(models.DutyScheduleCreate)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DutyCalendarService.CreateAndUpdate(r)
	})
}

func (dc DutyCalendarController) Update(ctx *gin.Context) {
	r := new(models.DutySchedule)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DutyCalendarService.Update(r)
	})
}

func (dc DutyCalendarController) Search(ctx *gin.Context) {
	r := new(models.DutyScheduleQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DutyCalendarService.Search(r)
	})
}

func (dc DutyCalendarController) GetCalendarUsers(ctx *gin.Context) {
	r := new(models.DutyScheduleQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DutyCalendarService.GetCalendarUsers(r)
	})
}
