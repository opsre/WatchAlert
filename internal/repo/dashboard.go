package repo

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	DashboardRepo struct {
		entryRepo
	}

	InterDashboardRepo interface {
		ListDashboardFolder(tenantId, query string) ([]models.DashboardFolders, error)
		GetDashboardFolder(tenantId, id string) (models.DashboardFolders, error)
		CreateDashboardFolder(fd models.DashboardFolders) error
		UpdateDashboardFolder(fd models.DashboardFolders) error
		DeleteDashboardFolder(tenantId, id string) error
	}
)

func newDashboardInterface(db *gorm.DB, g InterGormDBCli) InterDashboardRepo {
	return &DashboardRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (dr DashboardRepo) ListDashboardFolder(tenantId, query string) ([]models.DashboardFolders, error) {
	var (
		data []models.DashboardFolders
		db   = dr.db.Model(&models.DashboardFolders{})
	)

	db.Where("tenant_id = ?", tenantId)
	if query != "" {
		db.Where("name LIKE ? OR grafana_host LIKE ? OR grafana_folder_id LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	err := db.Find(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}

func (dr DashboardRepo) GetDashboardFolder(tenantId, id string) (models.DashboardFolders, error) {
	var (
		data models.DashboardFolders
		db   = dr.db.Model(&models.DashboardFolders{})
	)

	db.Where("tenant_id = ? and id = ?", tenantId, id)
	err := db.First(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}

func (dr DashboardRepo) CreateDashboardFolder(fd models.DashboardFolders) error {
	err := dr.g.Create(&models.DashboardFolders{}, fd)
	if err != nil {
		logc.Error(context.Background(), err)
		return err
	}
	return nil
}

func (dr DashboardRepo) UpdateDashboardFolder(fd models.DashboardFolders) error {
	u := Updates{
		Table: &models.DashboardFolders{},
		Where: map[string]interface{}{
			"tenant_id = ?": fd.TenantId,
			"id = ?":        fd.ID,
		},
		Updates: fd,
	}
	err := dr.g.Updates(u)
	if err != nil {
		logc.Error(context.Background(), err)
		return err
	}
	return nil
}

func (dr DashboardRepo) DeleteDashboardFolder(tenantId, id string) error {
	d := Delete{
		Table: &models.DashboardFolders{},
		Where: map[string]interface{}{
			"tenant_id = ?": tenantId,
			"id = ?":        id,
		},
	}

	err := dr.g.Delete(d)
	if err != nil {
		logc.Error(context.Background(), err)
		return err
	}

	return nil
}
