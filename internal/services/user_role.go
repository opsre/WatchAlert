package services

import (
	"time"
	"watchAlert/internal/ctx"
	models "watchAlert/internal/models"
	"watchAlert/internal/types"
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
	data, err := ur.ctx.DB.UserRole().List()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (ur userRoleService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestUserRoleCreate)

	err := ur.ctx.DB.UserRole().Create(models.UserRole{
		ID:          "ur-" + tools.RandId(),
		Name:        r.Name,
		Description: r.Description,
		Permissions: r.Permissions,
		UpdateAt:    time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ur userRoleService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestUserRoleUpdate)

	err := ur.ctx.DB.UserRole().Update(models.UserRole{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Permissions: r.Permissions,
		UpdateAt:    time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ur userRoleService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestUserRoleQuery)
	err := ur.ctx.DB.UserRole().Delete(r.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
