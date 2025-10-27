package services

import (
	"fmt"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"
)

type tenantService struct {
	ctx *ctx.Context
}

type InterTenantService interface {
	Create(req interface{}) (data interface{}, err interface{})
	Update(req interface{}) (data interface{}, err interface{})
	Delete(req interface{}) (data interface{}, err interface{})
	List(req interface{}) (data interface{}, err interface{})
	Get(req interface{}) (data interface{}, err interface{})
	AddUsersToTenant(req interface{}) (data interface{}, err interface{})
	DelUsersOfTenant(req interface{}) (data interface{}, err interface{})
	GetUsersForTenant(req interface{}) (data interface{}, err interface{})
	ChangeTenantUserRole(req interface{}) (data interface{}, err interface{})
}

func newInterTenantService(ctx *ctx.Context) InterTenantService {
	return &tenantService{
		ctx: ctx,
	}
}

func (ts tenantService) Create(req interface{}) (data interface{}, err interface{}) {
	r := req.(*types.RequestTenantCreate)
	tenant := models.Tenant{
		ID:               "tid-" + tools.RandId(),
		Name:             r.Name,
		UserId:           r.UserId,
		UpdateAt:         time.Now().Unix(),
		Manager:          r.Manager,
		Description:      r.Description,
		RuleNumber:       r.RuleNumber,
		UserNumber:       r.UserNumber,
		DutyNumber:       r.DutyNumber,
		NoticeNumber:     r.NoticeNumber,
		RemoveProtection: r.GetRemoveProtection(),
	}

	err = ts.ctx.DB.Tenant().Create(tenant)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ts tenantService) Update(req interface{}) (data interface{}, err interface{}) {
	r := req.(*types.RequestTenantUpdate)
	tenant := models.Tenant{
		ID:               r.ID,
		Name:             r.Name,
		UserId:           r.UserId,
		UpdateAt:         time.Now().Unix(),
		Manager:          r.Manager,
		Description:      r.Description,
		RuleNumber:       r.RuleNumber,
		UserNumber:       r.UserNumber,
		DutyNumber:       r.DutyNumber,
		NoticeNumber:     r.NoticeNumber,
		RemoveProtection: r.GetRemoveProtection(),
	}

	err = ts.ctx.DB.Tenant().Update(tenant)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ts tenantService) Delete(req interface{}) (data interface{}, err interface{}) {
	r := req.(*types.RequestTenantQuery)

	var t models.Tenant
	ts.ctx.DB.DB().Model(&models.Tenant{}).Where("id = ?", r.ID).Find(&t)

	if *t.GetRemoveProtection() {
		return nil, fmt.Errorf("删除失败, 删除保护已开启 关闭后再删除")
	}

	err = ts.ctx.DB.Tenant().Delete(r.ID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ts tenantService) List(req interface{}) (data interface{}, err interface{}) {
	r := req.(*types.RequestTenantQuery)
	data, err = ts.ctx.DB.Tenant().List(r.UserID)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (ts tenantService) Get(req interface{}) (data interface{}, err interface{}) {
	r := req.(*types.RequestTenantQuery)
	data, err = ts.ctx.DB.Tenant().Get(r.ID)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (ts tenantService) AddUsersToTenant(req interface{}) (data interface{}, err interface{}) {
	r := req.(*types.RequestTenantAddUsers)
	err = ts.ctx.DB.Tenant().AddTenantLinkedUsers(r.ID, r.Users, r.UserRole)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (ts tenantService) DelUsersOfTenant(req interface{}) (data interface{}, err interface{}) {
	r := req.(*types.RequestTenantQuery)
	err = ts.ctx.DB.Tenant().RemoveTenantLinkedUsers(r.ID, r.UserID)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (ts tenantService) GetUsersForTenant(req interface{}) (data interface{}, err interface{}) {
	r := req.(*types.RequestTenantQuery)
	data, err = ts.ctx.DB.Tenant().GetTenantLinkedUsers(r.ID)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (ts tenantService) ChangeTenantUserRole(req interface{}) (data interface{}, err interface{}) {
	r := req.(*types.RequestTenantChangeUserRole)
	err = ts.ctx.DB.Tenant().ChangeTenantUserRole(r.ID, r.UserID, r.UserRole)
	if err != nil {
		return nil, err
	}
	return nil, err
}
