package services

import (
	"fmt"
	"gorm.io/gorm"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

type (
	alertSubscribeService struct {
		ctx *ctx.Context
	}

	InterAlertSubscribeService interface {
		List(req interface{}) (interface{}, interface{})
		Get(req interface{}) (interface{}, interface{})
		Create(req interface{}) (interface{}, interface{})
		Delete(req interface{}) (interface{}, interface{})
	}
)

func newInterAlertSubscribe(ctx *ctx.Context) InterAlertSubscribeService {
	return alertSubscribeService{
		ctx: ctx,
	}
}

func (s alertSubscribeService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*models.AlertSubscribeQuery)
	list, err := s.ctx.DB.Subscribe().List(*r)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s alertSubscribeService) Get(req interface{}) (interface{}, interface{}) {
	r := req.(*models.AlertSubscribeQuery)
	get, _, err := s.ctx.DB.Subscribe().Get(*r)
	if err != nil {
		return nil, err
	}

	return get, nil
}

func (s alertSubscribeService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*models.AlertSubscribe)
	_, b, err := s.ctx.DB.Subscribe().Get(models.AlertSubscribeQuery{STenantId: r.STenantId, SUserId: r.SUserId, SRuleId: r.SRuleId})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if b {
		return nil, fmt.Errorf("用户已订阅该规则, 请勿重复创建!")
	}

	r.SId = "as-" + tools.RandId()
	r.SCreateAt = time.Now().Unix()
	err = s.ctx.DB.Subscribe().Create(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s alertSubscribeService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*models.AlertSubscribeQuery)
	err := s.ctx.DB.Subscribe().Delete(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
