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
		Delete(tenantId, ruleId string) error
		List(tenantId, ruleType, query string) ([]models.ProbingRule, error)
		Search(tenantId, ruleId string) (models.ProbingRule, error)
		AddRecord(history models.ProbingHistory) error
		GetRecord(ruleId string, dateRange int64) ([]models.ProbingHistory, error)
		DeleteRecord() error
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

func (p ProbingRepo) Delete(tenantId, ruleId string) error {
	del := Delete{
		Table: &models.ProbingRule{},
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

func (p ProbingRepo) List(tenantId, ruleType, query string) ([]models.ProbingRule, error) {
	var (
		data []models.ProbingRule
		db   = p.db.Model(&models.ProbingRule{})
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

func (p ProbingRepo) Search(tenantId, ruleId string) (models.ProbingRule, error) {
	var (
		data models.ProbingRule
		db   = p.db.Model(&models.ProbingRule{})
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

func (p ProbingRepo) AddRecord(history models.ProbingHistory) error {
	err := p.g.Create(models.ProbingHistory{}, history)
	if err != nil {
		logc.Errorf(context.Background(), err.Error())
		return err
	}
	return nil
}

func (p ProbingRepo) GetRecord(ruleId string, dateRange int64) ([]models.ProbingHistory, error) {
	var (
		data []models.ProbingHistory
		db   = p.db.Model(&models.ProbingHistory{})
	)

	// 计算起始时间戳（秒）
	now := time.Now().Unix()
	startTime := now - dateRange

	db.Where("rule_id = ?", ruleId).
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

func (p ProbingRepo) ChangeState(tenantId, ruleId string, state *bool) error {
	return p.db.Model(&models.ProbingRule{}).Where("tenant_id = ? AND rule_id = ?", tenantId, ruleId).Update("enabled", state).Error
}
