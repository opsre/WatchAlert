package services

import (
	"time"
	models "watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/tools"
)

type alertSilenceService struct {
	alertEvent models.AlertCurEvent
	ctx        *ctx.Context
}

type InterSilenceService interface {
	Create(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Delete(req interface{}) (interface{}, interface{})
	List(req interface{}) (interface{}, interface{})
}

func newInterSilenceService(ctx *ctx.Context) InterSilenceService {
	return &alertSilenceService{
		ctx: ctx,
	}
}

func (ass alertSilenceService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*models.AlertSilences)
	updateAt := time.Now().Unix()
	silenceEvent := models.AlertSilences{
		TenantId:      r.TenantId,
		Name:          r.Name,
		Id:            "s-" + tools.RandId(),
		StartsAt:      r.StartsAt,
		EndsAt:        r.EndsAt,
		UpdateAt:      updateAt,
		UpdateBy:      r.UpdateBy,
		FaultCenterId: r.FaultCenterId,
		Labels:        r.Labels,
		Comment:       r.Comment,
		Status:        1,
	}

	if r.StartsAt > updateAt {
		r.Status = 0
	}

	ass.ctx.Cache.Silence().PushMuteToFaultCenter(silenceEvent)
	err := ass.ctx.DB.Silence().Create(silenceEvent)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ass alertSilenceService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*models.AlertSilences)
	updateAt := time.Now().Unix()
	r.UpdateAt = updateAt

	if r.StartsAt > updateAt {
		r.Status = 0
	} else {
		r.Status = 1
	}

	ass.ctx.Cache.Silence().PushMuteToFaultCenter(*r)
	err := ass.ctx.DB.Silence().Update(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ass alertSilenceService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*models.AlertSilenceQuery)
	ass.ctx.Cache.Silence().RemoveMuteFromFaultCenter(r.TenantId, r.FaultCenterId, r.Id)
	err := ass.ctx.DB.Silence().Delete(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ass alertSilenceService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*models.AlertSilenceQuery)
	data, err := ass.ctx.DB.Silence().List(*r)
	if err != nil {
		return nil, err
	}

	return data, nil
}
