package services

import (
	"fmt"
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
	Get(req interface{}) (interface{}, interface{})
	List(req interface{}) (interface{}, interface{})
}

func newInterSilenceService(ctx *ctx.Context) InterSilenceService {
	return &alertSilenceService{
		ctx: ctx,
	}
}

func (ass alertSilenceService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*models.AlertSilences)
	createAt := time.Now().Unix()
	silenceEvent := models.AlertSilences{
		TenantId:       r.TenantId,
		Id:             "s-" + tools.RandId(),
		Fingerprint:    r.Fingerprint,
		Datasource:     r.Datasource,
		DatasourceType: r.DatasourceType,
		StartsAt:       r.StartsAt,
		EndsAt:         r.EndsAt,
		CreateBy:       r.CreateBy,
		CreateAt:       createAt,
		UpdateAt:       createAt,
		Comment:        r.Comment,
	}

	event, ok := ass.ctx.Redis.Silence().GetCache(models.AlertSilenceQuery{
		TenantId:    silenceEvent.TenantId,
		Fingerprint: silenceEvent.Fingerprint,
	})
	if ok && event != "" {
		return nil, fmt.Errorf("静默消息已存在, ID:%s", silenceEvent.Id)
	}

	muteAt := r.EndsAt - createAt
	duration := time.Duration(muteAt) * time.Second
	ass.ctx.Redis.Silence().SetCache(silenceEvent, duration)

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
	muteAt := r.EndsAt - r.StartsAt
	duration := time.Duration(muteAt) * time.Second
	ass.ctx.Redis.Silence().SetCache(*r, duration)

	err := ass.ctx.DB.Silence().Update(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ass alertSilenceService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*models.AlertSilenceQuery)
	err := ass.ctx.Redis.Silence().DelCache(*r)
	if err != nil {
		return nil, err
	}

	err = ass.ctx.DB.Silence().Delete(*r)
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

func (ass alertSilenceService) Get(req interface{}) (interface{}, interface{}) {
	r := req.(*models.AlertSilenceQuery)
	data, err := ass.ctx.DB.Silence().Get(*r)
	if err != nil {
		return nil, err
	}

	return data, nil
}
