package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	RuleTmplGroupRepo struct {
		entryRepo
	}

	InterRuleTmplGroupRepo interface {
		List(r models.RuleTemplateGroupQuery) ([]models.RuleTemplateGroup, error)
		Create(r models.RuleTemplateGroup) error
		Update(r models.RuleTemplateGroup) error
		Delete(r models.RuleTemplateGroupQuery) error
	}
)

func newRuleTmplGroupInterface(db *gorm.DB, g InterGormDBCli) InterRuleTmplGroupRepo {
	return &RuleTmplGroupRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (rtg RuleTmplGroupRepo) List(r models.RuleTemplateGroupQuery) ([]models.RuleTemplateGroup, error) {
	var data []models.RuleTemplateGroup
	db := rtg.db.Model(&models.RuleTemplateGroup{})
	db.Where("type = ?", r.Type)
	if r.Query != "" {
		db.Where("name LIKE ? OR description LIKE ?",
			"%"+r.Query+"%", "%"+r.Query+"%")
	}
	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}

	for k, v := range data {
		var ruleCount int64
		rtdb := rtg.db.Model(&models.RuleTemplate{})
		rtdb.Where("type = ?", r.Type)
		rtdb.Where("rule_group_name = ?", v.Name).Count(&ruleCount)
		data[k].Number = int(ruleCount)
	}

	return data, nil
}

func (rtg RuleTmplGroupRepo) Create(r models.RuleTemplateGroup) error {
	err := rtg.g.Create(models.RuleTemplateGroup{}, r)
	if err != nil {
		return err
	}

	return nil
}

func (rtg RuleTmplGroupRepo) Update(r models.RuleTemplateGroup) error {
	u := Updates{
		Table: models.RuleTemplateGroup{},
		Where: map[string]interface{}{
			"name = ?": r.Name,
		},
		Updates: r,
	}
	err := rtg.g.Updates(u)
	if err != nil {
		return err
	}

	return nil
}

func (rtg RuleTmplGroupRepo) Delete(r models.RuleTemplateGroupQuery) error {
	d := Delete{
		Table: &models.RuleTemplateGroup{},
		Where: map[string]interface{}{
			"name = ?": r.Name,
		},
	}

	err := rtg.g.Delete(d)
	if err != nil {
		return err
	}

	return nil
}
