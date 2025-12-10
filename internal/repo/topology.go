package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	topologyRepo struct {
		entryRepo
	}

	InterTopologyRepo interface {
		Create(params models.Topology) error
		Update(params models.Topology) error
		Delete(tenantId, id string) error
		List(tenantId, query string) ([]models.TopologyList, error)
		Get(tenantId, id string) (models.Topology, error)
		GetDetail(tenantId, id string) (models.Topology, error)
	}
)

func newInterTopologyRepo(db *gorm.DB, g InterGormDBCli) InterTopologyRepo {
	return &topologyRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (t topologyRepo) Create(params models.Topology) error {
	err := t.g.Create(&models.Topology{}, params)
	if err != nil {
		return err
	}
	return nil
}

func (t topologyRepo) Update(params models.Topology) error {
	u := Updates{
		Table: &models.Topology{},
		Where: map[string]interface{}{
			"tenant_id = ?": params.TenantId,
			"id = ?":        params.ID,
		},
		Updates: params,
	}
	err := t.g.Updates(u)
	if err != nil {
		return err
	}
	return nil
}

func (t topologyRepo) Delete(tenantId, id string) error {
	del := Delete{
		Table: &models.Topology{},
		Where: map[string]interface{}{
			"tenant_id = ?": tenantId,
			"id = ?":        id,
		},
	}
	err := t.g.Delete(del)
	if err != nil {
		return err
	}
	return nil
}

func (t topologyRepo) List(tenantId, query string) ([]models.TopologyList, error) {
	var (
		data []models.TopologyList
		db   = t.db.Model(&models.Topology{})
	)

	if tenantId != "" {
		db.Where("tenant_id = ?", tenantId)
	}
	if query != "" {
		db.Where("name LIKE ? OR id LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (t topologyRepo) Get(tenantId, id string) (models.Topology, error) {
	var (
		data models.Topology
		db   = t.db.Model(&models.Topology{})
	)

	if tenantId != "" {
		db.Where("tenant_id = ?", tenantId)
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

// GetDetail 获取拓扑的完整信息，包括nodes和edges
func (t topologyRepo) GetDetail(tenantId, id string) (models.Topology, error) {
	return t.Get(tenantId, id)
}