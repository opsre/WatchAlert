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
		List(groupType, query string) ([]models.RuleTemplateGroup, error)
		Create(r models.RuleTemplateGroup) error
		Update(r models.RuleTemplateGroup) error
		Delete(groupName string) error
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

func (rtg RuleTmplGroupRepo) List(groupType, query string) ([]models.RuleTemplateGroup, error) {
	var data []models.RuleTemplateGroup
	db := rtg.db.Model(&models.RuleTemplateGroup{})
	db.Where("type = ?", groupType)
	if query != "" {
		db.Where("name LIKE ? OR description LIKE ?",
			"%"+query+"%", "%"+query+"%")
	}
	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}

	for k, v := range data {
		var ruleCount int64
		rtdb := rtg.db.Model(&models.RuleTemplate{})
		rtdb.Where("type = ?", groupType)
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

func (rtg RuleTmplGroupRepo) Delete(groupName string) error {
	d := Delete{
		Table: &models.RuleTemplateGroup{},
		Where: map[string]interface{}{
			"name = ?": groupName,
		},
	}

	err := rtg.g.Delete(d)
	if err != nil {
		return err
	}

	return nil
}
