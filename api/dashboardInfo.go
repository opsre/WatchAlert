package api

import (
	"watchAlert/internal/ctx"
	"watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logc"
)

type dashboardInfoController struct{}

var DashboardInfoController = new(dashboardInfoController)

func (dashboardInfoController dashboardInfoController) API(gin *gin.RouterGroup) {
	system := gin.Group("system")
	system.Use(
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		system.GET("getDashboardInfo", dashboardInfoController.GetDashboardInfo)
	}
}

func (dashboardInfoController dashboardInfoController) GetDashboardInfo(context *gin.Context) {
	var c = ctx.DO()

	tid, _ := context.Get("TenantID")
	tidString := tid.(string)

	faultCenter, err := c.DB.FaultCenter().Get(tidString, context.Query("faultCenterId"), "")
	if err != nil {
		logc.Error(c.Ctx, err.Error())
		return
	}

	response.Success(context, types.ResponseDashboardInfo{
		CountAlertRules:   getRuleNumber(c, tidString),
		FaultCenterNumber: getFaultCenterNumber(c, tidString),
		UserNumber:        getUserNumber(c),
		CurAlertList:      getAlertList(c, faultCenter),
		AlarmDistribution: types.AlarmDistribution{
			P0: getAlarmDistribution(c, faultCenter, "P0"),
			P1: getAlarmDistribution(c, faultCenter, "P1"),
			P2: getAlarmDistribution(c, faultCenter, "P2"),
		},
	}, "success")
}

func getRuleNumber(ctx *ctx.Context, tenantId string) int64 {
	list, _, err := ctx.DB.Rule().List(tenantId, "", "", "", "", models.Page{
		Index: 0,
		Size:  10000,
	})
	if err != nil {
		return 0
	}
	return int64(len(list))
}

// getFaultCenterNumber 获取故障中心总数
func getFaultCenterNumber(ctx *ctx.Context, tenantId string) int64 {
	list, err := ctx.DB.FaultCenter().List(tenantId, "")
	if err != nil {
		logc.Error(ctx.Ctx, err.Error())
		return 0
	}
	return int64(len(list))
}

// getUserNumber 获取用户总数
func getUserNumber(ctx *ctx.Context) int64 {
	list, err := ctx.DB.User().List("", "")
	if err != nil {
		logc.Error(ctx.Ctx, err.Error())
		return 0
	}
	return int64(len(list))
}

// getAlertList 获取当前告警 annotations 列表

func getAlertList(ctx *ctx.Context, faultCenter models.FaultCenter) []types.AlertList {
	events, err := ctx.Redis.Alert().GetAllEvents(models.BuildAlertEventCacheKey(faultCenter.TenantId, faultCenter.ID))
	if err != nil {
		return nil
	}

	var list []types.AlertList
	var uniq = make(map[string]struct{})
	for _, event := range events {
		if _, ok := uniq[event.RuleName]; ok {
			continue
		}

		list = append(list, types.AlertList{Severity: event.Severity, RuleName: event.RuleName, FaultCenterId: event.FaultCenterId, TiggerTime: event.FirstTriggerTime})
		uniq[event.RuleName] = struct{}{}
	}
	return list
}

// getAlarmDistribution 获取告警分布
func getAlarmDistribution(ctx *ctx.Context, faultCenter models.FaultCenter, severity string) int64 {
	events, err := ctx.Redis.Alert().GetAllEvents(models.BuildAlertEventCacheKey(faultCenter.TenantId, faultCenter.ID))
	if err != nil {
		return 0
	}

	var number int64
	for _, event := range events {
		if event.Severity == severity {
			number++
		}
	}
	return number
}
