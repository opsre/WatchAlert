package repo

import (
	"fmt"
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	RuleGroupRepo struct {
		entryRepo
	}

	InterRuleGroupRepo interface {
		List(tenantId, query string, page models.Page) ([]models.RuleGroups, int64, error)
		Create(req models.RuleGroups) error
		Update(req models.RuleGroups) error
		Delete(tenantId, id string) error
	}
)

func newRuleGroupInterface(db *gorm.DB, g InterGormDBCli) InterRuleGroupRepo {
	return &RuleGroupRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (r RuleGroupRepo) List(tenantId, query string, page models.Page) ([]models.RuleGroups, int64, error) {
	var (
		data  []models.RuleGroups
		db    = r.db.Model(&models.RuleGroups{})
		count int64
	)

	pageIndexInt := page.Index
	pageSizeInt := page.Size

	db.Where("tenant_id = ?", tenantId)

	if query != "" {
		db.Where("id LIKE ? OR name LIKE ? OR description LIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	db.Count(&count)

	db.Limit(int(pageSizeInt)).Offset(int((pageIndexInt - 1) * pageSizeInt))

	err := db.Find(&data).Error
	if err != nil {
		return nil, 0, err
	}

	for k, v := range data {
		var resRules []models.AlertRule
		r.db.Model(&models.AlertRule{}).Where("tenant_id = ? AND rule_group_id = ?", tenantId, v.ID).Find(&resRules)
		data[k].Number = len(resRules)
	}
	return data, count, nil
}

func (r RuleGroupRepo) Create(req models.RuleGroups) error {
	var resGroup models.RuleGroups
	r.db.Model(&models.RuleGroups{}).Where("name = ?", req.Name).First(&resGroup)
	if resGroup.Name != "" {
		return fmt.Errorf("规则组名称已存在")
	}

	err := r.g.Create(models.RuleGroups{}, req)
	if err != nil {
		return err
	}

	return nil
}

func (r RuleGroupRepo) Update(req models.RuleGroups) error {
	u := Updates{
		Table: &models.RuleGroups{},
		Where: map[string]interface{}{
			"tenant_id = ?": req.TenantId,
			"id = ?":        req.ID,
		},
		Updates: req,
	}

	err := r.g.Updates(u)
	if err != nil {
		return err
	}

	return nil
}

func (r RuleGroupRepo) Delete(tenantId, id string) error {
	var ruleNum int64
	r.db.Model(&models.AlertRule{}).Where("tenant_id = ? AND rule_group_id = ?", tenantId, id).
		Count(&ruleNum)
	if ruleNum != 0 {
		return fmt.Errorf("无法删除规则组 %s, 因为规则组不为空", id)
	}

	d := Delete{
		Table: models.RuleGroups{},
		Where: map[string]interface{}{
			"tenant_id = ?": tenantId,
			"id = ?":        id,
		},
	}

	err := r.g.Delete(d)
	if err != nil {
		return err
	}

	return nil
}
