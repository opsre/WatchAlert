package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	RuleTmplRepo struct {
		entryRepo
	}

	InterRuleTmplRepo interface {
		List(r models.RuleTemplateQuery) ([]models.RuleTemplate, error)
		Create(r models.RuleTemplate) error
		Update(r models.RuleTemplate) error
		Delete(r models.RuleTemplateQuery) error
	}
)

func newRuleTmplInterface(db *gorm.DB, g InterGormDBCli) InterRuleTmplRepo {
	return &RuleTmplRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (rt RuleTmplRepo) List(r models.RuleTemplateQuery) ([]models.RuleTemplate, error) {
	var data []models.RuleTemplate
	db := rt.db.Model(&models.RuleTemplate{}).Where("rule_group_name = ?", r.RuleGroupName)
	db.Where("type = ?", r.Type)
	if r.Query != "" {
		db.Where("rule_name LIKE ? OR datasource_type LIKE ?",
			"%"+r.Query+"%", "%"+r.Query+"%")
	}

	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (rt RuleTmplRepo) Create(r models.RuleTemplate) error {
	err := rt.g.Create(models.RuleTemplate{}, r)
	if err != nil {
		return err
	}

	return nil
}

func (rt RuleTmplRepo) Update(r models.RuleTemplate) error {
	u := Updates{
		Table: models.RuleTemplate{},
		Where: map[string]interface{}{
			"rule_name = ?": r.RuleName,
		},
		Updates: r,
	}
	err := rt.g.Updates(u)
	if err != nil {
		return err
	}

	return nil
}

func (rt RuleTmplRepo) Delete(r models.RuleTemplateQuery) error {
	d := Delete{
		Table: models.RuleTemplate{},
		Where: map[string]interface{}{
			"rule_group_name = ?": r.RuleGroupName,
			"rule_name = ?":       r.RuleName,
		},
	}

	err := rt.g.Delete(d)
	if err != nil {
		return err
	}

	return nil
}
