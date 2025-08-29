package repo

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	faultCenterRepo struct {
		entryRepo
	}

	InterFaultCenterRepo interface {
		Create(params models.FaultCenter) error
		Update(params models.FaultCenter) error
		Delete(params models.FaultCenterQuery) error
		List(params models.FaultCenterQuery) ([]models.FaultCenter, error)
		Get(params models.FaultCenterQuery) (models.FaultCenter, error)
		Reset(params models.FaultCenter) error
	}
)

func newInterFaultCenterRepo(db *gorm.DB, g InterGormDBCli) InterFaultCenterRepo {
	return &faultCenterRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (f faultCenterRepo) Create(params models.FaultCenter) error {
	err := f.g.Create(&models.FaultCenter{}, params)
	if err != nil {
		logc.Error(context.Background(), err)
		return err
	}
	return nil
}

func (f faultCenterRepo) Update(params models.FaultCenter) error {
	u := Updates{
		Table: &models.FaultCenter{},
		Where: map[string]interface{}{
			"tenant_id = ?": params.TenantId,
			"id = ?":        params.ID,
		},
		Updates: params,
	}
	err := f.g.Updates(u)
	if err != nil {
		logc.Error(context.Background(), err)
		return err
	}
	return nil
}

func (f faultCenterRepo) Delete(params models.FaultCenterQuery) error {
	del := Delete{
		Table: &models.FaultCenter{},
		Where: map[string]interface{}{
			"tenant_id = ?": params.TenantId,
			"id = ?":        params.ID,
		},
	}
	err := f.g.Delete(del)
	if err != nil {
		logc.Error(context.Background(), err)
		return err
	}
	return nil
}

func (f faultCenterRepo) List(params models.FaultCenterQuery) ([]models.FaultCenter, error) {
	var db = f.db.Model(&models.FaultCenter{})
	var data []models.FaultCenter
	if params.TenantId != "" {
		db.Where("tenant_id = ?", params.TenantId)
	}

	if params.Query != "" {
		db.Where("name LIKE ? OR id LIKE ?", "%"+params.Query+"%", "%"+params.Query+"%")
	}

	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (f faultCenterRepo) Get(params models.FaultCenterQuery) (models.FaultCenter, error) {
	var db = f.db.Model(&models.FaultCenter{})
	var data models.FaultCenter
	if params.Name != "" {
		db.Where("name = ?", params.Name)
	}
	if params.ID != "" {
		db.Where("id = ?", params.ID)
	}
	err := db.First(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func (f faultCenterRepo) Reset(params models.FaultCenter) error {
	var update []string
	if params.Name != "" {
		update = []string{"name", params.Name}
	}

	if params.Description != "" {
		update = []string{"description", params.Description}
	}

	if params.AggregationType != "" {
		update = []string{"aggregation_type", params.AggregationType}
	}

	if update != nil {
		err := f.g.Update(Update{
			Table: &models.FaultCenter{},
			Where: map[string]interface{}{
				"id = ?": params.ID,
			},
			Update: update,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
