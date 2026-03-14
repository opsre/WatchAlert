package repo

import (
	"watchAlert/internal/models"

	"gorm.io/gorm"
)

type (
	RecordingRuleRepo struct {
		entryRepo
	}

	InterRecordingRuleRepo interface {
		Get(tenantId, ruleId string) (models.RecordingRule, error)
		List(tenantId, datasourceType, query, status string, page models.Page) ([]models.RecordingRule, int64, error)
		Create(r models.RecordingRule) error
		Update(r models.RecordingRule) error
		Delete(tenantId, ruleId string) error
		GetRuleObject(ruleId string) models.RecordingRule
		ChangeStatus(tenantId, ruleId string, state *bool) error
	}
)

func newRecordingRuleInterface(db *gorm.DB, g InterGormDBCli) InterRecordingRuleRepo {
	return &RecordingRuleRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (rr RecordingRuleRepo) Get(tenantId, ruleId string) (models.RecordingRule, error) {
	var data models.RecordingRule

	db := rr.db.Model(&models.RecordingRule{})
	db.Where("tenant_id = ? AND rule_id = ?", tenantId, ruleId)
	err := db.First(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}

func (rr RecordingRuleRepo) List(tenantId, datasourceType, query, status string, page models.Page) ([]models.RecordingRule, int64, error) {
	var (
		data  []models.RecordingRule
		count int64
	)

	db := rr.db.Model(&models.RecordingRule{})
	db.Where("tenant_id = ?", tenantId)
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

func (rr RecordingRuleRepo) Create(r models.RecordingRule) error {
	err := rr.g.Create(models.RecordingRule{}, r)
	if err != nil {
		return err
	}

	return nil
}

func (rr RecordingRuleRepo) Update(r models.RecordingRule) error {
	u := Updates{
		Table: &models.RecordingRule{},
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

func (rr RecordingRuleRepo) Delete(tenantId, ruleId string) error {
	var recordingRule models.RecordingRule
	d := Delete{
		Table: recordingRule,
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

func (rr RecordingRuleRepo) GetRuleObject(ruleId string) models.RecordingRule {
	var data models.RecordingRule
	rr.DB().Model(&models.RecordingRule{}).
		Where("rule_id = ?", ruleId).
		First(&data)

	return data
}

func (rr RecordingRuleRepo) ChangeStatus(tenantId, ruleId string, state *bool) error {
	return rr.DB().Model(&models.RecordingRule{}).
		Where("tenant_id = ? AND rule_id = ?", tenantId, ruleId).
		Update("enabled", state).Error
}
