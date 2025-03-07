package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	AiRepo struct {
		entryRepo
	}
	InterAiRepo interface {
		Get(ruleId string) (models.AiContentRecord, bool, error)
		Create(data models.AiContentRecord) error
		Update(data models.AiContentRecord) error
	}
)

func newAiRepoInterface(db *gorm.DB, g InterGormDBCli) InterAiRepo {
	return &AiRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (a AiRepo) Get(ruleId string) (models.AiContentRecord, bool, error) {
	var (
		db   = a.DB().Model(&models.AiContentRecord{})
		data models.AiContentRecord
	)

	db.Where("rule_id = ?", ruleId)
	if err := db.First(&data).Error; err != nil {
		return data, false, err
	}

	return data, true, nil
}

func (a AiRepo) Create(data models.AiContentRecord) error {
	err := a.g.Create(&models.AiContentRecord{}, &data)
	if err != nil {
		return err
	}

	return nil
}

func (a AiRepo) Update(data models.AiContentRecord) error {
	err := a.g.Updates(Updates{
		Table: &models.AiContentRecord{},
		Where: map[string]interface{}{
			"rule_id = ?": data.RuleId,
		},
		Updates: &data,
	})
	if err != nil {
		return err
	}

	return nil
}
