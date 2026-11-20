package repo

import (
	"encoding/json"
	"watchAlert/internal/models"

	"gorm.io/gorm"
)

type (
	faultCenterRepo struct {
		entryRepo
	}

	InterFaultCenterRepo interface {
		Create(params models.FaultCenter) error
		Update(params models.FaultCenter) error
		Delete(tenantId, id string) error
		List(tenantId, query string) ([]models.FaultCenter, error)
		Get(tenantId, id, name string) (models.FaultCenter, error)
		Reset(tenantId, id, name, description, aggregationType string, aggregationLabels []string) error
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
		return err
	}
	return nil
}

func (f faultCenterRepo) Delete(tenantId, id string) error {
	del := Delete{
		Table: &models.FaultCenter{},
		Where: map[string]interface{}{
			"tenant_id = ?": tenantId,
			"id = ?":        id,
		},
	}
	err := f.g.Delete(del)
	if err != nil {
		return err
	}
	return nil
}

func (f faultCenterRepo) List(tenantId, query string) ([]models.FaultCenter, error) {
	var (
		data []models.FaultCenter
		db   = f.db.Model(&models.FaultCenter{})
	)

	if tenantId != "" {
		db.Where("tenant_id = ?", tenantId)
	}
	if query != "" {
		db.Where("name LIKE ? OR id LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (f faultCenterRepo) Get(tenantId, id, name string) (models.FaultCenter, error) {
	var (
		data models.FaultCenter
		db   = f.db.Model(&models.FaultCenter{})
	)

	if tenantId != "" {
		db.Where("tenant_id = ?", tenantId)
	}
	if name != "" {
		db.Where("name = ?", name)
	}
	if id != "" {
		db.Where("id = ?", id)
	}

	err := db.First(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func (f faultCenterRepo) Reset(tenantId, id, name, description, aggregationType string, aggregationLabels []string) error {
	updates := make(map[string]interface{})

	if name != "" {
		updates["name"] = name
	}
	if description != "" {
		updates["description"] = description
	}
	if aggregationType != "" {
		aggregationLabelsJSON, err := json.Marshal(aggregationLabels)
		if err != nil {
			return err
		}
		updates["aggregation_type"] = aggregationType
		updates["aggregation_labels"] = string(aggregationLabelsJSON)
	}

	if len(updates) == 0 {
		return nil
	}

	return f.g.Updates(Updates{
		Table: &models.FaultCenter{},
		Where: map[string]interface{}{
			"tenant_id = ?": tenantId,
			"id = ?":        id,
		},
		Updates: updates,
	})
}
