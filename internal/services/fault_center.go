package services

import (
	"time"
	"watchAlert/alert"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
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
	alert.ConsumerWork.Submit(fc)

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
	alert.ConsumerWork.Stop(r.ID)
	alert.ConsumerWork.Submit(fc)

	return nil, nil
}

func (f faultCenterService) Delete(req interface{}) (data interface{}, err interface{}) {
	r := req.(*types.RequestFaultCenterQuery)
	err = f.ctx.DB.FaultCenter().Delete(r.TenantId, r.ID)
	if err != nil {
		return nil, err
	}

	f.ctx.Redis.FaultCenter().RemoveFaultCenterInfo(models.BuildFaultCenterInfoCacheKey(r.TenantId, r.ID))
	alert.ConsumerWork.Stop(r.ID)

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

	alert.ConsumerWork.Stop(r.ID)
	data, err = f.ctx.DB.FaultCenter().Get(r.TenantId, r.ID, r.Name)
	f.ctx.Redis.FaultCenter().PushFaultCenterInfo(data.(models.FaultCenter))
	alert.ConsumerWork.Submit(data.(models.FaultCenter))

	return nil, nil
}
