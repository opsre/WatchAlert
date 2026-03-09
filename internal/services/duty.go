package services

import (
	"fmt"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"
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
	r := req.(*types.RequestDutyManagementQuery)
	data, err := dms.ctx.DB.Duty().List(r.TenantId)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (dms *dutyManageService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestDutyManagementCreate)
	ok := dms.ctx.DB.Duty().GetQuota(r.TenantId)
	if !ok {
		return nil, fmt.Errorf("创建失败, 配额不足")
	}

	err := dms.ctx.DB.Duty().Create(models.DutyManagement{
		TenantId:    r.TenantId,
		ID:          "dt-" + tools.RandId(),
		Name:        r.Name,
		Manager:     r.Manager,
		Description: r.Description,
		CurDutyUser: r.CurDutyUser,
		UpdateBy:    r.UpdateBy,
		UpdateAt:    time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (dms *dutyManageService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestDutyManagementUpdate)
	err := dms.ctx.DB.Duty().Update(models.DutyManagement{
		TenantId:    r.TenantId,
		ID:          r.ID,
		Name:        r.Name,
		Manager:     r.Manager,
		Description: r.Description,
		CurDutyUser: r.CurDutyUser,
		UpdateBy:    r.UpdateBy,
		UpdateAt:    time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (dms *dutyManageService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestDutyManagementQuery)
	err := dms.ctx.DB.Duty().Delete(r.TenantId, r.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (dms *dutyManageService) Get(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestDutyManagementQuery)
	data, err := dms.ctx.DB.Duty().Get(r.TenantId, r.ID)
	if err != nil {
		return nil, err
	}
	return data, nil
}
