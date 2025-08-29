package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	subscribeRepo struct {
		entryRepo
	}

	InterSubscribeRepo interface {
		List(r models.AlertSubscribeQuery) ([]models.AlertSubscribe, error)
		Get(r models.AlertSubscribeQuery) (models.AlertSubscribe, bool, error)
		Create(r models.AlertSubscribe) error
		Delete(r models.AlertSubscribeQuery) error
	}
)

func newInterSubscribeRepo(db *gorm.DB, g InterGormDBCli) InterSubscribeRepo {
	return &subscribeRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (s subscribeRepo) List(r models.AlertSubscribeQuery) ([]models.AlertSubscribe, error) {
	var (
		data []models.AlertSubscribe
		db   = s.db.Model(models.AlertSubscribe{})
	)

	db.Where("s_tenant_id = ?", r.STenantId)
	if r.SRuleId != "" {
		db.Where("s_rule_id = ?", r.SRuleId)
	}
	if r.Query != "" {
		db.Where("s_rule_id LIKE ? or s_rule_name LIKE ? or s_rule_type LIKE ?", "%"+r.Query+"%", "%"+r.Query+"%", "%"+r.Query+"%")
	}

	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s subscribeRepo) Get(r models.AlertSubscribeQuery) (models.AlertSubscribe, bool, error) {
	var (
		data models.AlertSubscribe
		db   = s.db.Model(models.AlertSubscribe{})
	)
	db.Where("s_tenant_id = ?", r.STenantId)
	if r.SId != "" {
		db.Where("s_id = ?", r.SId)

	}
	if r.SUserId != "" {
		db.Where("s_user_id = ?", r.SUserId)

	}
	if r.SRuleId != "" {
		db.Where("s_rule_id = ?", r.SRuleId)

	}
	err := db.First(&data).Error
	if err != nil {
		return data, false, err
	}

	return data, true, nil
}

func (s subscribeRepo) Create(r models.AlertSubscribe) error {
	err := s.g.Create(models.AlertSubscribe{}, r)
	if err != nil {
		return err
	}

	return nil
}

func (s subscribeRepo) Delete(r models.AlertSubscribeQuery) error {
	err := s.g.Delete(Delete{
		Table: models.AlertSubscribe{},
		Where: map[string]interface{}{
			"s_tenant_id": r.STenantId,
			"s_id":        r.SId,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
