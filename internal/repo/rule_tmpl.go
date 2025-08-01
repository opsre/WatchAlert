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
		List(tmplGroup, tmplType, query string) ([]models.RuleTemplate, error)
		Create(r models.RuleTemplate) error
		Update(r models.RuleTemplate) error
		Delete(tmplGroupName, tmplName string) error
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

func (rt RuleTmplRepo) List(tmplGroup, tmplType, query string) ([]models.RuleTemplate, error) {
	var data []models.RuleTemplate
	db := rt.db.Model(&models.RuleTemplate{}).Where("rule_group_name = ?", tmplGroup)
	db.Where("type = ?", tmplType)
	if query != "" {
		db.Where("rule_name LIKE ? OR datasource_type LIKE ?",
			"%"+query+"%", "%"+query+"%")
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

func (rt RuleTmplRepo) Delete(tmplGroupName, tmplName string) error {
	d := Delete{
		Table: models.RuleTemplate{},
		Where: map[string]interface{}{
			"rule_group_name = ?": tmplGroupName,
			"rule_name = ?":       tmplName,
		},
	}

	err := rt.g.Delete(d)
	if err != nil {
		return err
	}

	return nil
}
