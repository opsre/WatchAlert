package services

import (
	"time"
	"watchAlert/alert"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/client"
	"watchAlert/pkg/tools"
)

type (
	faultCenterService struct {
		ctx *ctx.Context
	}

	InterFaultCenterService interface {
		Create(req interface{}) (data interface{}, err interface{})
		Update(req interface{}) (data interface{}, err interface{})
		Delete(req interface{}) (data interface{}, err interface{})
		List(req interface{}) (data interface{}, err interface{})
		Get(req interface{}) (data interface{}, err interface{})
		Reset(req interface{}) (data interface{}, err interface{})
		Slo(req interface{}) (data interface{}, err interface{})
	}
)

func newInterFaultCenterService(ctx *ctx.Context) InterFaultCenterService {
	return &faultCenterService{
		ctx: ctx,
	}
}

func (f faultCenterService) Create(req interface{}) (data interface{}, err interface{}) {
	r := req.(*types.RequestFaultCenterCreate)
	fc := models.FaultCenter{
		TenantId:             r.TenantId,
		ID:                   "fc-" + tools.RandId(),
		Name:                 r.Name,
		Description:          r.Description,
		NoticeIds:            r.NoticeIds,
		NoticeRoutes:         r.NoticeRoutes,
		RepeatNoticeInterval: r.RepeatNoticeInterval,
		RecoverNotify:        r.RecoverNotify,
		AggregationType:      r.AggregationType,
		CreateAt:             time.Now().Unix(),
		RecoverWaitTime:      r.RecoverWaitTime,
		IsUpgradeEnabled:     r.IsUpgradeEnabled,
		UpgradableSeverity:   r.UpgradableSeverity,
		UpgradeStrategy:      r.UpgradeStrategy,
	}

	err = f.ctx.DB.FaultCenter().Create(fc)
	if err != nil {
		return nil, err
	}

	f.ctx.Redis.FaultCenter().PushFaultCenterInfo(fc)

	// 判断当前节点角色
	if alert.IsLeader() {
		// Leader: 直接启动消费协程
		alert.ConsumerWork.Submit(fc)
	} else {
		// Follower: 发布消息通知 Leader
		tools.PublishReloadMessage(f.ctx.Ctx, client.Redis, tools.ChannelFaultCenterReload, tools.ReloadMessage{
			Action:   tools.ActionCreate,
			ID:       fc.ID,
			TenantID: fc.TenantId,
			Name:     fc.Name,
		})
	}

	return nil, nil
}

func (f faultCenterService) Update(req interface{}) (data interface{}, err interface{}) {
	r := req.(*types.RequestFaultCenterUpdate)
	fc := models.FaultCenter{
		TenantId:             r.TenantId,
		ID:                   r.ID,
		Name:                 r.Name,
		Description:          r.Description,
		NoticeIds:            r.NoticeIds,
		NoticeRoutes:         r.NoticeRoutes,
		RepeatNoticeInterval: r.RepeatNoticeInterval,
		RecoverNotify:        r.RecoverNotify,
		AggregationType:      r.AggregationType,
		CreateAt:             r.CreateAt,
		RecoverWaitTime:      r.RecoverWaitTime,
		IsUpgradeEnabled:     r.IsUpgradeEnabled,
		UpgradableSeverity:   r.UpgradableSeverity,
		UpgradeStrategy:      r.UpgradeStrategy,
	}

	err = f.ctx.DB.FaultCenter().Update(fc)
	if err != nil {
		return nil, err
	}

	f.ctx.Redis.FaultCenter().PushFaultCenterInfo(fc)

	// 判断当前节点角色
	if alert.IsLeader() {
		// Leader: 直接重启消费协程
		alert.ConsumerWork.Stop(r.ID)
		alert.ConsumerWork.Submit(fc)
	} else {
		// Follower: 发布消息通知 Leader
		tools.PublishReloadMessage(f.ctx.Ctx, client.Redis, tools.ChannelFaultCenterReload, tools.ReloadMessage{
			Action:   tools.ActionUpdate,
			ID:       fc.ID,
			TenantID: fc.TenantId,
			Name:     fc.Name,
		})
	}

	return nil, nil
}

func (f faultCenterService) Delete(req interface{}) (data interface{}, err interface{}) {
	r := req.(*types.RequestFaultCenterQuery)
	err = f.ctx.DB.FaultCenter().Delete(r.TenantId, r.ID)
	if err != nil {
		return nil, err
	}

	f.ctx.Redis.FaultCenter().RemoveFaultCenterInfo(models.BuildFaultCenterInfoCacheKey(r.TenantId, r.ID))

	// 判断当前节点角色
	if alert.IsLeader() {
		// Leader: 直接停止消费协程
		alert.ConsumerWork.Stop(r.ID)
	} else {
		// Follower: 发布消息通知 Leader
		tools.PublishReloadMessage(f.ctx.Ctx, client.Redis, tools.ChannelFaultCenterReload, tools.ReloadMessage{
			Action:   tools.ActionDelete,
			ID:       r.ID,
			TenantID: r.TenantId,
			Name:     r.ID,
		})
	}

	return nil, nil
}

func (f faultCenterService) List(req interface{}) (data interface{}, err interface{}) {
	r := req.(*types.RequestFaultCenterQuery)
	data, err = f.ctx.DB.FaultCenter().List(r.TenantId, r.Query)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return data, nil
	}

	faultCenters := data.([]models.FaultCenter)
	for index, fc := range data.([]models.FaultCenter) {
		events, err := f.ctx.Redis.Alert().GetAllEvents(models.BuildAlertEventCacheKey(fc.TenantId, fc.ID))
		if err != nil {
			return nil, err
		}

		for _, event := range events {
			switch event.Status {
			case models.StatePreAlert:
				faultCenters[index].CurrentPreAlertNumber++
			case models.StateAlerting:
				faultCenters[index].CurrentAlertNumber++
			case models.StateSilenced:
				faultCenters[index].CurrentMuteNumber++
			case models.StatePendingRecovery:
				faultCenters[index].CurrentRecoverNumber++
			}
		}
	}

	return faultCenters, nil
}

func (f faultCenterService) Get(req interface{}) (data interface{}, err interface{}) {
	r := req.(*types.RequestFaultCenterQuery)
	data, err = f.ctx.DB.FaultCenter().Get(r.TenantId, r.ID, r.Name)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (f faultCenterService) Reset(req interface{}) (data interface{}, err interface{}) {
	r := req.(*types.RequestFaultCenterReset)
	err = f.ctx.DB.FaultCenter().Reset(r.TenantId, r.ID, r.Name, r.Description, r.AggregationType)
	if err != nil {
		return nil, err
	}

	data, err = f.ctx.DB.FaultCenter().Get(r.TenantId, r.ID, r.Name)
	if err != nil {
		return nil, err
	}
	f.ctx.Redis.FaultCenter().PushFaultCenterInfo(data.(models.FaultCenter))

	// 判断当前节点角色
	if alert.IsLeader() {
		// Leader: 直接重启消费协程
		alert.ConsumerWork.Stop(r.ID)
		alert.ConsumerWork.Submit(data.(models.FaultCenter))
	} else {
		// Follower: 发布消息通知 Leader
		tools.PublishReloadMessage(f.ctx.Ctx, client.Redis, tools.ChannelFaultCenterReload, tools.ReloadMessage{
			Action:   tools.ActionUpdate,
			ID:       r.ID,
			TenantID: r.TenantId,
			Name:     r.Name,
		})
	}

	return nil, nil
}

func (f faultCenterService) Slo(req interface{}) (data interface{}, err interface{}) {
	r := req.(*types.RequestFaultCenterQuery)
	// 拉取近7天的历史事件（一次性拉全量，后面按天聚合）
	eventsResp, eventErr := f.ctx.DB.Event().GetHistoryEvent(types.RequestAlertHisEventQuery{
		TenantId:      r.TenantId,
		FaultCenterId: r.ID,
		Page: models.Page{
			Index: 1,
			Size:  999999,
		},
	})
	if eventErr != nil {
		return nil, eventErr
	}
	eventsList := eventsResp.List

	// 准备返回的 7 天数组（按时间从旧到新）
	mttrArr := make([]float64, 0, 7)
	mttaArr := make([]float64, 0, 7)

	now := time.Now()
	// 生成最近 7 天，从 6 天前 到 今天（顺序：旧 -> 新）
	for i := 6; i >= 0; i-- {
		day := now.AddDate(0, 0, -i)
		start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location()).Unix()
		end := time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 0, day.Location()).Unix()

		mttr, _ := f.calculateMTTRRange(eventsList, start, end)
		mtta, _ := f.calculateMTTARange(eventsList, start, end)

		mttrArr = append(mttrArr, mttr)
		mttaArr = append(mttaArr, mtta)
	}

	return types.RequestFaultCenterSLO{
		MTTR: mttrArr,
		MTTA: mttaArr,
	}, nil
}

func (f faultCenterService) calculateMTTRRange(eventsList []models.AlertHisEvent, start, end int64) (float64, error) {
	var totalRepairTime float64
	var recoveredCount int

	for _, event := range eventsList {
		if event.RecoverTime > 0 && event.RecoverTime >= start && event.RecoverTime <= end {
			repairTime := float64(event.RecoverTime - event.FirstTriggerTime)
			if repairTime > 0 {
				totalRepairTime += repairTime
				recoveredCount++
			}
		}
	}

	if recoveredCount == 0 {
		return 0, nil
	}
	return totalRepairTime / float64(recoveredCount), nil
}

func (f faultCenterService) calculateMTTARange(eventsList []models.AlertHisEvent, start, end int64) (float64, error) {
	var totalRespTime float64
	var ackCount int

	for _, ev := range eventsList {
		ack := ev.ConfirmState.ConfirmActionTime
		first := ev.FirstTriggerTime
		if ack <= 0 || first <= 0 {
			continue
		}
		if ack < start || ack > end {
			continue
		}
		resp := float64(ack - first)
		if resp >= 0 {
			totalRespTime += resp
			ackCount++
		}
	}

	if ackCount == 0 {
		return 0, nil
	}
	return totalRespTime / float64(ackCount), nil
}
