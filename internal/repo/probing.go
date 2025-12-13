package repo

import (
	"context"
	"watchAlert/internal/models"

	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
)

type (
	ProbingRepo struct {
		entryRepo
	}

	InterProbingRepo interface {
		Create(d models.ProbeRule) error
		Update(d models.ProbeRule) error
		Delete(tenantId, ruleId string) error
		List(tenantId, ruleType, query string) ([]models.ProbeRule, error)
		Search(tenantId, ruleId string) (models.ProbeRule, error)
		ChangeState(tenantId, ruleId string, state *bool) error
	}
)

func newProbingRepoInterface(db *gorm.DB, g InterGormDBCli) InterProbingRepo {
	return &ProbingRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (p ProbingRepo) Create(d models.ProbeRule) error {

	err := p.g.Create(models.ProbeRule{}, d)
	if err != nil {
		logc.Errorf(context.Background(), err.Error())
		return err
	}
	return nil
}

func (p ProbingRepo) Update(d models.ProbeRule) error {
	u := Updates{
		Table: &models.ProbeRule{},
		Where: map[string]interface{}{
			"tenant_id = ?": d.TenantId,
			"rule_id = ?":   d.RuleId,
		},
		Updates: d,
	}
	err := p.g.Updates(u)
	if err != nil {
		logc.Errorf(context.Background(), err.Error())
		return err
	}
	return nil
}

func (p ProbingRepo) Delete(tenantId, ruleId string) error {
	del := Delete{
		Table: &models.ProbeRule{},
		Where: map[string]interface{}{
			"tenant_id = ?": tenantId,
			"rule_id = ?":   ruleId,
		},
	}
	err := p.g.Delete(del)
	if err != nil {
		logc.Errorf(context.Background(), err.Error())
		return err
	}
	return nil
}

func (p ProbingRepo) List(tenantId, ruleType, query string) ([]models.ProbeRule, error) {
	var (
		data []models.ProbeRule
		db   = p.db.Model(&models.ProbeRule{})
	)

	db.Where("tenant_id = ?", tenantId)
	if ruleType != "" {
		db.Where("rule_type = ?", ruleType)
	}
	if query != "" {
		db.Where("probing_endpoint_config LIKE ?", "%"+query+"%")
	}

	err := db.Find(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return data, nil
		}
		return data, err
	}
	return data, nil
}

func (p ProbingRepo) Search(tenantId, ruleId string) (models.ProbeRule, error) {
	var (
		data models.ProbeRule
		db   = p.db.Model(&models.ProbeRule{})
	)

	if tenantId != "" {
		db.Where("tenant_id = ?", tenantId)
	}
	db.Where("rule_id = ? ", ruleId)

	err := db.First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return data, nil
		}
		return data, err
	}
	return data, nil
}

func (p ProbingRepo) ChangeState(tenantId, ruleId string, state *bool) error {
	return p.db.Model(&models.ProbeRule{}).Where("tenant_id = ? AND rule_id = ?", tenantId, ruleId).Update("enabled", state).Error
}
