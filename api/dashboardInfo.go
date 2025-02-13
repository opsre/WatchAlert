package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logc"
	"watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/response"
)

type DashboardInfoController struct {
	models.AlertCurEvent
}

func (di DashboardInfoController) API(gin *gin.RouterGroup) {
	system := gin.Group("system")
	system.Use(
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		system.GET("getDashboardInfo", di.GetDashboardInfo)
	}
}

type ResponseDashboardInfo struct {
	CountAlertRules   int64             `json:"countAlertRules"`
	FaultCenterNumber int64             `json:"faultCenterNumber"`
	UserNumber        int64             `json:"userNumber"`
	CurAlertList      []string          `json:"curAlertList"`
	AlarmDistribution AlarmDistribution `json:"alarmDistribution"`
}

type AlarmDistribution struct {
	P0 int64 `json:"P0"`
	P1 int64 `json:"P1"`
	P2 int64 `json:"P2"`
}

func (di DashboardInfoController) GetDashboardInfo(context *gin.Context) {
	var c = ctx.DO()
	tid, _ := context.Get("TenantID")
	tidString := tid.(string)

	faultCenter, err := ctx.DB.FaultCenter().Get(models.FaultCenterQuery{
		TenantId: tidString,
		ID:       context.Query("faultCenterId"),
	})
	if err != nil {
		logc.Error(ctx.Ctx, err.Error())
		return
	}

	response.Success(context, ResponseDashboardInfo{
		CountAlertRules:   getRuleNumber(c, tidString),
		FaultCenterNumber: getFaultCenterNumber(c, tidString),
		UserNumber:        getUserNumber(c),
		CurAlertList:      getAlertList(c, faultCenter),
		AlarmDistribution: AlarmDistribution{
			P0: getAlarmDistribution(c, faultCenter, "P0"),
			P1: getAlarmDistribution(c, faultCenter, "P1"),
			P2: getAlarmDistribution(c, faultCenter, "P2"),
		},
	}, "success")
}

func getRuleNumber(ctx *ctx.Context, tenantId string) int64 {
	list, err := ctx.DB.Rule().List(models.AlertRuleQuery{
		TenantId: tenantId,
		Page: models.Page{
			Index: 0,
			Size:  10000,
		},
	})
	if err != nil {
		return 0
	}
	return int64(len(list.List))
}

// getFaultCenterNumber 获取故障中心总数
func getFaultCenterNumber(ctx *ctx.Context, tenantId string) int64 {
	list, err := ctx.DB.FaultCenter().List(models.FaultCenterQuery{TenantId: tenantId})
	if err != nil {
		logc.Error(ctx.Ctx, err.Error())
		return 0
	}
	return int64(len(list))
}

// getUserNumber 获取用户总数
func getUserNumber(ctx *ctx.Context) int64 {
	list, err := ctx.DB.User().List()
	if err != nil {
		logc.Error(ctx.Ctx, err.Error())
		return 0
	}
	return int64(len(list))
}

// getAlertList 获取当前告警 annotations 列表
func getAlertList(ctx *ctx.Context, faultCenter models.FaultCenter) []string {
	events, err := ctx.Redis.Event().GetAllEventsForFaultCenter(faultCenter.GetFaultCenterKey())
	if err != nil {
		return nil
	}

	var annotations []string
	for _, event := range events {
		annotations = append(annotations, event.Annotations)
	}
	return annotations
}

// getAlarmDistribution 获取告警分布
func getAlarmDistribution(ctx *ctx.Context, faultCenter models.FaultCenter, severity string) int64 {
	events, err := ctx.Redis.Event().GetAllEventsForFaultCenter(faultCenter.GetFaultCenterKey())
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
