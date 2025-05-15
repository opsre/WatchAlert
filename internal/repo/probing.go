package repo

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
	"time"
	"watchAlert/internal/models"
)

type (
	ProbingRepo struct {
		entryRepo
	}

	InterProbingRepo interface {
		Create(d models.ProbingRule) error
		Update(d models.ProbingRule) error
		Delete(d models.ProbingRuleQuery) error
		List(d models.ProbingRuleQuery) ([]models.ProbingRule, error)
		Search(d models.ProbingRuleQuery) (models.ProbingRule, error)
		AddRecord(history models.ProbingHistory) error
		GetRecord(query models.ReqProbingHistory) ([]models.ProbingHistory, error)
		DeleteRecord() error
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

func (p ProbingRepo) Create(d models.ProbingRule) error {

	err := p.g.Create(models.ProbingRule{}, d)
	if err != nil {
		logc.Errorf(context.Background(), err.Error())
		return err
	}
	return nil
}

func (p ProbingRepo) Update(d models.ProbingRule) error {
	u := Updates{
		Table: &models.ProbingRule{},
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

func (p ProbingRepo) Delete(d models.ProbingRuleQuery) error {
	del := Delete{
		Table: &models.ProbingRule{},
		Where: map[string]interface{}{
			"tenant_id = ?": d.TenantId,
			"rule_id = ?":   d.RuleId,
		},
	}
	err := p.g.Delete(del)
	if err != nil {
		logc.Errorf(context.Background(), err.Error())
		return err
	}
	return nil
}

func (p ProbingRepo) List(d models.ProbingRuleQuery) ([]models.ProbingRule, error) {
	var (
		data []models.ProbingRule
		db   = p.db.Model(&models.ProbingRule{})
	)
	db.Where("tenant_id = ?", d.TenantId)

	if d.RuleType != "" {
		db.Where("rule_type = ?", d.RuleType)
	}

	if d.Query != "" {
		db.Where("probing_endpoint_config LIKE ?", "%"+d.Query+"%")
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

func (p ProbingRepo) Search(d models.ProbingRuleQuery) (models.ProbingRule, error) {
	var (
		data models.ProbingRule
		db   = p.db.Model(&models.ProbingRule{})
	)
	if d.TenantId != "" {
		db.Where("tenant_id = ?", d.TenantId)
	}

	db.Where("rule_id = ? ", d.RuleId)
	err := db.Find(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return data, nil
		}
		return data, err
	}
	return data, nil
}

func (p ProbingRepo) AddRecord(history models.ProbingHistory) error {
	err := p.g.Create(models.ProbingHistory{}, history)
	if err != nil {
		logc.Errorf(context.Background(), err.Error())
		return err
	}
	return nil
}

func (p ProbingRepo) GetRecord(query models.ReqProbingHistory) ([]models.ProbingHistory, error) {
	var (
		data []models.ProbingHistory
		db   = p.db.Model(&models.ProbingHistory{})
	)
	db.Where("rule_id = ? ", query.RuleId)
	// 计算起始时间戳（秒）
	now := time.Now().Unix()
	startTime := now - query.DateRange

	db.Where("rule_id = ?", query.RuleId).
		Where("timestamp BETWEEN ? AND ?", startTime, now)

	err := db.Find(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return data, nil
		}
		return data, err
	}

	return data, nil
}

func (p ProbingRepo) DeleteRecord() error {
	var saveDays int64 = 3600 * 24

	now := time.Now().Unix()
	startTime := now - saveDays

	del := Delete{
		Table: &models.ProbingHistory{},
		Where: map[string]interface{}{
			"timestamp < ?": startTime,
		},
	}
	err := p.g.Delete(del)
	if err != nil {
		logc.Errorf(context.Background(), err.Error())
		return err
	}
	return nil
}
