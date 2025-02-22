package services

import (
	"time"
	"watchAlert/alert"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
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
	r := req.(*models.FaultCenter)
	r.ID = "fc-" + tools.RandId()
	r.CreateAt = time.Now().Unix()
	err = f.ctx.DB.FaultCenter().Create(*r)
	if err != nil {
		return nil, err
	}

	f.ctx.Redis.FaultCenter().PushFaultCenterInfo(*r)
	alert.ConsumerWork.Submit(*r)

	return nil, nil
}

func (f faultCenterService) Update(req interface{}) (data interface{}, err interface{}) {
	r := req.(*models.FaultCenter)
	err = f.ctx.DB.FaultCenter().Update(*r)
	if err != nil {
		return nil, err
	}

	f.ctx.Redis.FaultCenter().PushFaultCenterInfo(*r)
	alert.ConsumerWork.Stop(r.ID)
	alert.ConsumerWork.Submit(*r)

	return nil, nil
}

func (f faultCenterService) Delete(req interface{}) (data interface{}, err interface{}) {
	r := req.(*models.FaultCenterQuery)
	err = f.ctx.DB.FaultCenter().Delete(*r)
	if err != nil {
		return nil, err
	}

	f.ctx.Redis.FaultCenter().RemoveFaultCenterInfo(models.BuildCacheInfoKey(r.TenantId, r.ID))
	alert.ConsumerWork.Stop(r.ID)

	return nil, nil
}

func (f faultCenterService) List(req interface{}) (data interface{}, err interface{}) {
	r := req.(*models.FaultCenterQuery)
	data, err = f.ctx.DB.FaultCenter().List(*r)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return data, nil
	}

	faultCenters := data.([]models.FaultCenter)
	for index, fc := range data.([]models.FaultCenter) {
		events, err := f.ctx.Redis.Event().GetAllEventsForFaultCenter(fc.GetFaultCenterKey())
		if err != nil {
			return nil, err
		}

		for _, event := range events {
			switch event.Status {
			case 1:
				faultCenters[index].CurrentAlertNumber++
			case 2:
				faultCenters[index].CurrentMuteNumber++
			case 3:
				faultCenters[index].CurrentRecoverNumber++
			}
		}
	}

	return faultCenters, nil
}

func (f faultCenterService) Get(req interface{}) (data interface{}, err interface{}) {
	r := req.(*models.FaultCenterQuery)
	data, err = f.ctx.DB.FaultCenter().Get(*r)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (f faultCenterService) Reset(req interface{}) (data interface{}, err interface{}) {
	r := req.(*models.FaultCenter)
	err = f.ctx.DB.FaultCenter().Reset(*r)
	if err != nil {
		return nil, err
	}

	alert.ConsumerWork.Stop(r.ID)
	data, err = f.ctx.DB.FaultCenter().Get(models.FaultCenterQuery{ID: r.ID})
	alert.ConsumerWork.Submit(data.(models.FaultCenter))

	return nil, nil
}
