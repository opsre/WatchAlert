package services

import (
	"time"
	"watchAlert/internal/ctx"
	models "watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

type userRoleService struct {
	ctx *ctx.Context
}

type InterUserRoleService interface {
	List(req interface{}) (interface{}, interface{})
	Create(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Delete(req interface{}) (interface{}, interface{})
}

func newInterUserRoleService(ctx *ctx.Context) InterUserRoleService {
	return &userRoleService{
		ctx: ctx,
	}
}

func (ur userRoleService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*models.UserRoleQuery)

	data, err := ur.ctx.DB.UserRole().List(*r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (ur userRoleService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*models.UserRole)

	r.ID = "ur-" + tools.RandId()
	r.CreateAt = time.Now().Unix()

	err := ur.ctx.DB.UserRole().Create(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ur userRoleService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*models.UserRole)

	err := ur.ctx.DB.UserRole().Update(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ur userRoleService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*models.UserRoleQuery)
	err := ur.ctx.DB.UserRole().Delete(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
