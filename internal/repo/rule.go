package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	RuleRepo struct {
		entryRepo
	}

	InterRuleRepo interface {
		GetQuota(id string) bool
		Get(tenantId, ruleGroupId, ruleId string) (models.AlertRule, error)
		List(tenantId, ruleGroupId, datasourceType, query, status string, page models.Page) ([]models.AlertRule, int64, error)
		Create(r models.AlertRule) error
		Update(r models.AlertRule) error
		Delete(tenantId, ruleId string) error
		GetRuleIsExist(ruleId string) bool
		GetRuleObject(ruleId string) models.AlertRule
		ChangeStatus(tenantId, ruleGroupId, ruleId string, state *bool) error
	}
)

func newRuleInterface(db *gorm.DB, g InterGormDBCli) InterRuleRepo {
	return &RuleRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (rr RuleRepo) GetQuota(id string) bool {
	var (
		db     = rr.db.Model(&models.Tenant{})
		data   models.Tenant
		Number int64
	)

	db.Where("id = ?", id)
	db.Find(&data)

	rr.db.Model(&models.AlertRule{}).Where("tenant_id = ?", id).Count(&Number)

	if Number < data.RuleNumber {
		return true
	}

	return false
}

func (rr RuleRepo) Get(tenantId, ruleGroupId, ruleId string) (models.AlertRule, error) {
	var data models.AlertRule

	db := rr.db.Model(&models.AlertRule{})
	db.Where("tenant_id = ? AND rule_group_id = ? AND rule_id = ?", tenantId, ruleGroupId, ruleId)
	err := db.First(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}

func (rr RuleRepo) List(tenantId, ruleGroupId, datasourceType, query, status string, page models.Page) ([]models.AlertRule, int64, error) {
	var (
		data  []models.AlertRule
		count int64
	)

	db := rr.db.Model(&models.AlertRule{})
	db.Where("tenant_id = ?", tenantId)
	if ruleGroupId != "" {
		db.Where("rule_group_id = ?", ruleGroupId)
	}

	if datasourceType != "" {
		db.Where("datasource_type = ?", datasourceType)
	}

	if query != "" {
		db.Where("rule_id LIKE ? OR rule_name LIKE ? OR description LIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	if status != "all" {
		switch status {
		case "enabled":
			db.Where("enabled = ?", true)
		case "disabled":
			db.Where("enabled = ?", false)
		}
	}

	db.Count(&count)

	db.Limit(int(page.Size)).Offset(int((page.Index - 1) * page.Size))

	err := db.Find(&data).Error

	if err != nil {
		return nil, 0, err
	}

	return data, count, nil
}

func (rr RuleRepo) Create(r models.AlertRule) error {
	err := rr.g.Create(models.AlertRule{}, r)
	if err != nil {
		return err
	}

	return nil
}

func (rr RuleRepo) Update(r models.AlertRule) error {
	u := Updates{
		Table: &models.AlertRule{},
		Where: map[string]interface{}{
			"tenant_id = ?": r.TenantId,
			"rule_id = ?":   r.RuleId,
		},
		Updates: r,
	}

	err := rr.g.Updates(u)
	if err != nil {
		return err
	}

	return nil
}

func (rr RuleRepo) Delete(tenantId, ruleId string) error {
	var alertRule models.AlertRule
	d := Delete{
		Table: alertRule,
		Where: map[string]interface{}{
			"tenant_id = ?": tenantId,
			"rule_id = ?":   ruleId,
		},
	}

	err := rr.g.Delete(d)
	if err != nil {
		return err
	}

	return nil
}

func (rr RuleRepo) GetRuleIsExist(ruleId string) bool {
	var ruleNum int64
	rr.DB().Model(&models.AlertRule{}).
		Where("rule_id = ? AND enabled = ?", ruleId, "1").
		Count(&ruleNum)
	if ruleNum > 0 {
		return true
	}

	return false
}

func (rr RuleRepo) GetRuleObject(ruleId string) models.AlertRule {
	var data models.AlertRule
	rr.DB().Model(&models.AlertRule{}).
		Where("rule_id = ?", ruleId).
		First(&data)

	return data
}

func (rr RuleRepo) ChangeStatus(tenantId, ruleGroupId, ruleId string, state *bool) error {
	return rr.DB().Model(&models.AlertRule{}).
		Where("tenant_id = ? AND rule_group_id = ? AND rule_id = ?", tenantId, ruleGroupId, ruleId).
		Update("enabled", state).Error
}
