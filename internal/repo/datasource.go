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
		List(r models.DatasourceQuery) ([]models.AlertDataSource, error)
		Search(r models.DatasourceQuery) ([]models.AlertDataSource, error)
		Get(r models.DatasourceQuery) (models.AlertDataSource, error)
		Create(r models.AlertDataSource) error
		Update(r models.AlertDataSource) error
		Delete(r models.DatasourceQuery) error
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

func (ds DatasourceRepo) List(r models.DatasourceQuery) ([]models.AlertDataSource, error) {
	var data []models.AlertDataSource

	db := ds.db.Model(&models.AlertDataSource{})
	if r.TenantId != "" {
		db.Where("tenant_id = ?", r.TenantId)
	}
	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (ds DatasourceRepo) Search(r models.DatasourceQuery) ([]models.AlertDataSource, error) {
	var db = ds.db.Model(&models.AlertDataSource{})
	var data []models.AlertDataSource

	db.Where("tenant_id = ?", r.TenantId)
	if r.Id != "" {
		db.Where("id = ?", r.Id)
	}
	if r.Type != "" {
		db.Where("type = ?", r.Type)
	}
	if r.Query != "" {
		db.Where("type LIKE ? OR id LIKE ? OR name LIKE ? OR description LIKE ?", "%"+r.Query+"%", "%"+r.Query+"%", "%"+r.Query+"%", "%"+r.Query+"%")
	}

	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (ds DatasourceRepo) Get(r models.DatasourceQuery) (models.AlertDataSource, error) {
	db := ds.db.Model(&models.AlertDataSource{})
	db.Where("id = ?", r.Id)

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
			"id = ?":        r.Id,
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

func (ds DatasourceRepo) Delete(r models.DatasourceQuery) error {
	var ruleNum int64
	ds.DB().Model(&models.AlertRule{}).Where("tenant_id = ? AND datasource_id_list LIKE ?", r.TenantId, "%"+r.Id+"%").Count(&ruleNum)
	if ruleNum != 0 {
		return fmt.Errorf("无法删除数据源 %s, 因为已有告警规则绑定", r.Id)
	}

	data := Delete{
		Table: models.AlertDataSource{},
		Where: map[string]interface{}{
			"tenant_id = ?": r.TenantId,
			"id = ?":        r.Id,
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
