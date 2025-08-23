package repo

import (
	"fmt"
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	DatasourceRepo struct {
		entryRepo
	}

	InterDatasourceRepo interface {
		List(tenantId, datasourceId, datasourceType, query string) ([]models.AlertDataSource, error)
		Get(datasourceId string) (models.AlertDataSource, error)
		Create(r models.AlertDataSource) error
		Update(r models.AlertDataSource) error
		Delete(tenantId, datasourceId string) error
		GetInstance(datasourceId string) (models.AlertDataSource, error)
	}
)

func newDatasourceInterface(db *gorm.DB, g InterGormDBCli) InterDatasourceRepo {
	return &DatasourceRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (ds DatasourceRepo) List(tenantId, datasourceId, datasourceType, query string) ([]models.AlertDataSource, error) {
	var db = ds.db.Model(&models.AlertDataSource{})
	var data []models.AlertDataSource

	if tenantId != "" {
		db.Where("tenant_id = ?", tenantId)
	}
	if datasourceId != "" {
		db.Where("id = ?", datasourceId)
	}
	if datasourceType != "" {
		db.Where("type = ?", datasourceType)
	}
	if query != "" {
		db.Where("type LIKE ? OR id LIKE ? OR name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (ds DatasourceRepo) Get(datasourceId string) (models.AlertDataSource, error) {
	db := ds.db.Model(&models.AlertDataSource{})
	db.Where("id = ?", datasourceId)

	var data models.AlertDataSource
	err := db.First(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}

func (ds DatasourceRepo) Create(r models.AlertDataSource) error {
	err := ds.g.Create(models.AlertDataSource{}, r)
	if err != nil {
		return err
	}
	return nil
}

func (ds DatasourceRepo) Update(r models.AlertDataSource) error {
	data := Updates{
		Table: models.AlertDataSource{},
		Where: map[string]interface{}{
			"id = ?":        r.ID,
			"tenant_id = ?": r.TenantId,
		},
		Updates: r,
	}
	err := ds.g.Updates(data)
	if err != nil {
		return err
	}
	return nil
}

func (ds DatasourceRepo) Delete(tenantId, datasourceId string) error {
	var ruleNum int64
	ds.DB().Model(&models.AlertRule{}).Where("tenant_id = ? AND datasource_id_list LIKE ?", tenantId, "%"+datasourceId+"%").Count(&ruleNum)
	if ruleNum != 0 {
		return fmt.Errorf("无法删除数据源 %s, 因为已有告警规则绑定", datasourceId)
	}

	data := Delete{
		Table: models.AlertDataSource{},
		Where: map[string]interface{}{
			"tenant_id = ?": tenantId,
			"id = ?":        datasourceId,
		},
	}
	err := ds.g.Delete(data)
	if err != nil {
		return err
	}
	return nil
}

func (ds DatasourceRepo) GetInstance(datasourceId string) (models.AlertDataSource, error) {
	var data models.AlertDataSource
	var db = ds.DB().Model(&models.AlertDataSource{})
	db.Where("id = ?", datasourceId)
	err := db.First(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}
