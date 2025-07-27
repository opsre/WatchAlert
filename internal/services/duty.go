package services

import (
	"fmt"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
)

type dutyManageService struct {
	ctx *ctx.Context
}

type InterDutyManageService interface {
	List(req interface{}) (interface{}, interface{})
	Create(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Delete(req interface{}) (interface{}, interface{})
	Get(req interface{}) (interface{}, interface{})
}

func newInterDutyManageService(ctx *ctx.Context) InterDutyManageService {
	return &dutyManageService{
		ctx: ctx,
	}
}

func (dms *dutyManageService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*models.DutyManagementQuery)
	data, err := dms.ctx.DB.Duty().List(*r)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (dms *dutyManageService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*models.DutyManagement)
	ok := dms.ctx.DB.Duty().GetQuota(r.TenantId)
	if !ok {
		return models.DutyManagement{}, fmt.Errorf("创建失败, 配额不足")
	}

	err := dms.ctx.DB.Duty().Create(*r)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (dms *dutyManageService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*models.DutyManagement)
	err := dms.ctx.DB.Duty().Update(*r)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (dms *dutyManageService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*models.DutyManagementQuery)
	err := dms.ctx.DB.Duty().Delete(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (dms *dutyManageService) Get(req interface{}) (interface{}, interface{}) {
	r := req.(*models.DutyManagementQuery)
	data, err := dms.ctx.DB.Duty().Get(*r)
	if err != nil {
		return nil, err
	}
	return data, nil
}
