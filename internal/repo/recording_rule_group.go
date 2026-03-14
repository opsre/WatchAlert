package repo

import (
	"fmt"
	"watchAlert/internal/models"

	"gorm.io/gorm"
)

type (
	RecordingRuleGroupRepo struct {
		entryRepo
	}

	InterRecordingRuleGroupRepo interface {
		List(tenantId, query string, page models.Page) ([]models.RecordingRuleGroup, int64, error)
		Create(req *models.RecordingRuleGroup) error
		Update(req *models.RecordingRuleGroup) error
		Delete(tenantId string, id int64) error
		Get(tenantId string, id int64) (models.RecordingRuleGroup, error)
	}
)

func newRecordingRuleGroupInterface(db *gorm.DB, g InterGormDBCli) InterRecordingRuleGroupRepo {
	return &RecordingRuleGroupRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (r RecordingRuleGroupRepo) List(tenantId, query string, page models.Page) ([]models.RecordingRuleGroup, int64, error) {
	var (
		data  []models.RecordingRuleGroup
		db    = r.db.Model(&models.RecordingRuleGroup{})
		count int64
	)

	db.Where("tenant_id = ?", tenantId)

	if query != "" {
		db.Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	db.Count(&count)

	db.Limit(int(page.Size)).Offset(int((page.Index - 1) * page.Size))

	err := db.Find(&data).Error
	if err != nil {
		return nil, 0, err
	}

	return data, count, nil
}

func (r RecordingRuleGroupRepo) Create(req *models.RecordingRuleGroup) error {
	var resGroup models.RecordingRuleGroup
	r.db.Model(&models.RecordingRuleGroup{}).Where("tenant_id = ? AND name = ?", req.TenantId, req.Name).First(&resGroup)
	if resGroup.ID > 0 {
		return fmt.Errorf("规则组名称已存在")
	}

	err := r.g.Create(&models.RecordingRuleGroup{}, req)
	if err != nil {
		return err
	}

	return nil
}

func (r RecordingRuleGroupRepo) Update(req *models.RecordingRuleGroup) error {
	u := Updates{
		Table: &models.RecordingRuleGroup{},
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

func (r RecordingRuleGroupRepo) Delete(tenantId string, id int64) error {
	var ruleNum int64
	r.db.Model(&models.RecordingRule{}).Where("tenant_id = ? AND rule_group_id = ?", tenantId, id).Count(&ruleNum)
	if ruleNum != 0 {
		return fmt.Errorf("禁止删除, 规则组 %d 不为空", id)
	}

	d := Delete{
		Table: models.RecordingRuleGroup{},
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

func (r RecordingRuleGroupRepo) Get(tenantId string, id int64) (models.RecordingRuleGroup, error) {
	var data models.RecordingRuleGroup
	db := r.db.Model(&models.RecordingRuleGroup{})
	db.Where("tenant_id = ? AND id = ?", tenantId, id)
	err := db.First(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}
