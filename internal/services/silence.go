package services

import (
	"time"
	"watchAlert/internal/ctx"
	models "watchAlert/internal/models"
	"watchAlert/internal/types"
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
	r := req.(*types.RequestSilenceCreate)
	updateAt := time.Now().Unix()
	silence := models.AlertSilences{
		TenantId:      r.TenantId,
		Name:          r.Name,
		ID:            "s-" + tools.RandId(),
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

	ass.ctx.Redis.Silence().PushAlertMute(silence)
	err := ass.ctx.DB.Silence().Create(silence)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ass alertSilenceService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestSilenceUpdate)
	silence := models.AlertSilences{
		TenantId:      r.TenantId,
		Name:          r.Name,
		ID:            r.ID,
		StartsAt:      r.StartsAt,
		EndsAt:        r.EndsAt,
		UpdateAt:      time.Now().Unix(),
		UpdateBy:      r.UpdateBy,
		FaultCenterId: r.FaultCenterId,
		Labels:        r.Labels,
		Comment:       r.Comment,
		Status:        1,
	}

	if r.StartsAt > r.UpdateAt {
		r.Status = 0
	} else {
		r.Status = 1
	}

	ass.ctx.Redis.Silence().PushAlertMute(silence)
	err := ass.ctx.DB.Silence().Update(silence)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ass alertSilenceService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestSilenceQuery)
	ass.ctx.Redis.Silence().RemoveAlertMute(r.TenantId, r.FaultCenterId, r.ID)
	err := ass.ctx.DB.Silence().Delete(r.TenantId, r.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ass alertSilenceService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestSilenceQuery)
	data, count, err := ass.ctx.DB.Silence().List(r.TenantId, r.FaultCenterId, r.Query, r.Page)
	if err != nil {
		return nil, err
	}

	return types.ResponseSilenceList{
		List: data,
		Page: models.Page{
			Total: count,
			Index: r.Page.Index,
			Size:  r.Page.Size,
		},
	}, nil
}
