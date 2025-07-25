package services

import (
	"watchAlert/internal/ctx"
)

type (
	userPermissionService struct {
		ctx *ctx.Context
	}

	InterUserPermissionService interface {
		List() (interface{}, interface{})
	}
)

func newInterUserPermissionService(ctx *ctx.Context) InterUserPermissionService {
	return &userPermissionService{
		ctx: ctx,
	}
}

func (up userPermissionService) List() (interface{}, interface{}) {
	data, err := up.ctx.DB.UserPermissions().List()
	if err != nil {
		return nil, err
	}

	return data, nil
}
