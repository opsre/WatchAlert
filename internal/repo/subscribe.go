package repo

import (
	"watchAlert/internal/models"

	"gorm.io/gorm"
)

type (
	subscribeRepo struct {
		entryRepo
	}

	InterSubscribeRepo interface {
		List(tenantId, ruleId, query string) ([]models.AlertSubscribe, error)
		Get(tenantId, sid, userId, ruleId string) (models.AlertSubscribe, bool, error)
		Create(r models.AlertSubscribe) error
		Delete(tenantId, sid string) error
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

func (s subscribeRepo) List(tenantId, ruleId, query string) ([]models.AlertSubscribe, error) {
	var (
		data []models.AlertSubscribe
		db   = s.db.Model(models.AlertSubscribe{})
	)

	db.Where("s_tenant_id = ?", tenantId)
	if ruleId != "" {
		db.Where("s_rule_id = ?", ruleId)
	}
	if query != "" {
		db.Where("s_rule_id LIKE ? or s_rule_name LIKE ? or s_rule_type LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s subscribeRepo) Get(tenantId, sid, userId, ruleId string) (models.AlertSubscribe, bool, error) {
	var (
		data models.AlertSubscribe
		db   = s.db.Model(models.AlertSubscribe{})
	)
	db.Where("s_tenant_id = ?", tenantId)
	if sid != "" {
		db.Where("s_id = ?", sid)

	}
	if userId != "" {
		db.Where("s_user_id = ?", userId)

	}
	if ruleId != "" {
		db.Where("s_rule_id = ?", ruleId)

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

func (s subscribeRepo) Delete(tenantId, sid string) error {
	err := s.g.Delete(Delete{
		Table: models.AlertSubscribe{},
		Where: map[string]interface{}{
			"s_tenant_id": tenantId,
			"s_id":        sid,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
