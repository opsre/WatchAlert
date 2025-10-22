package services

import (
	"fmt"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"

	"gorm.io/gorm"
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
	r := req.(*types.RequestSubscribeQuery)
	list, err := s.ctx.DB.Subscribe().List(r.STenantId, r.SRuleId, r.Query)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s alertSubscribeService) Get(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestSubscribeQuery)
	get, _, err := s.ctx.DB.Subscribe().Get(r.STenantId, r.SId, r.SUserId, r.SRuleId)
	if err != nil {
		return nil, err
	}

	return get, nil
}

func (s alertSubscribeService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestSubscribeCreate)
	_, b, err := s.ctx.DB.Subscribe().Get(r.STenantId, "", r.SUserId, r.SRuleId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if b {
		return nil, fmt.Errorf("用户已订阅该规则, 请勿重复创建!")
	}

	subscribe := models.AlertSubscribe{
		SId:               "as-" + tools.RandId(),
		STenantId:         r.STenantId,
		SUserId:           r.SUserId,
		SUserEmail:        r.SUserEmail,
		SRuleId:           r.SRuleId,
		SRuleName:         r.SRuleName,
		SRuleType:         r.SRuleType,
		SRuleSeverity:     r.SRuleSeverity,
		SNoticeSubject:    r.SNoticeSubject,
		SNoticeTemplateId: r.SNoticeTemplateId,
		SFilter:           r.SFilter,
		SCreateAt:         time.Now().Unix(),
	}

	err = s.ctx.DB.Subscribe().Create(subscribe)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s alertSubscribeService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestSubscribeQuery)
	err := s.ctx.DB.Subscribe().Delete(r.STenantId, r.SId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
