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
	CountAlertRules   int64                    `json:"countAlertRules"`
	CurMutes          int64                    `json:"curMutes"`
	CurAlerts         int                      `json:"curAlerts"`
	CurAlertList      []string                 `json:"curAlertList"`
	AlarmDistribution AlarmDistribution        `json:"alarmDistribution"`
	ServiceResource   []models.ServiceResource `json:"serviceResource"`
}

type AlarmDistribution struct {
	P0 int `json:"P0"`
	P1 int `json:"P1"`
	P2 int `json:"P2"`
}

func (di DashboardInfoController) GetDashboardInfo(context *gin.Context) {
	var (
		// 规则总数
		countAlertRules int64
		// 当前告警
		keys []string
		c    = ctx.DO()
	)

	tid, _ := context.Get("TenantID")
	tidString := tid.(string)
	// 告警分布
	alarmDistribution := make(map[string]int)
	c.DB.DB().Model(&models.AlertRule{}).Where("tenant_id = ?", tidString).Count(&countAlertRules)

	cursor := uint64(0)
	pattern := tidString + ":" + models.FiringAlertCachePrefix + "*"
	// 每次获取的键数量
	count := int64(100)

	for {
		var curKeys []string
		var err error

		curKeys, cursor, err = c.Redis.Redis().Scan(cursor, pattern, count).Result()
		if err != nil {
			break
		}

		keys = append(keys, curKeys...)

		if cursor == 0 {
			break
		}
	}

	var curAlertList []string
	for _, v := range keys {
		alarmDistribution[c.Redis.Event().GetCache(v).Severity] += 1
		if len(curAlertList) >= 5 {
			continue
		}
		curAlertList = append(curAlertList, c.Redis.Event().GetCache(v).Annotations)
	}

	res, err := c.DB.Silence().List(models.AlertSilenceQuery{TenantId: tidString, Status: 0})
	if err != nil {
		logc.Errorf(c.Ctx, "List silence failed, err: %s", err.Error())
		return
	}

	var resource []models.ServiceResource
	c.DB.DB().Model(&models.ServiceResource{}).Find(&resource)

	response.Success(context, ResponseDashboardInfo{
		CountAlertRules: countAlertRules,
		CurMutes:        res.Total,
		CurAlerts:       len(keys),
		CurAlertList:    curAlertList,
		AlarmDistribution: AlarmDistribution{
			P0: alarmDistribution["P0"],
			P1: alarmDistribution["P1"],
			P2: alarmDistribution["P2"],
		},
		ServiceResource: resource,
	}, "success")

}
