package repo

import (
	"context"
	"fmt"
	"time"
	"watchAlert/internal/models"

	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
)

type (
	NoticeRepo struct {
		entryRepo
	}

	InterNoticeRepo interface {
		Get(tenantId, id string) (models.AlertNotice, error)
		GetQuota(id string) bool
		List(tenantId, noticeTmplId, query string) ([]models.AlertNotice, error)
		Create(r models.AlertNotice) error
		Update(r models.AlertNotice) error
		Delete(tenantId, id string) error
		AddRecord(r models.NoticeRecord) error
		ListRecord(tenantId, eventId, severity, status, query string, page models.Page) (models.ResponseNoticeRecords, error)
		CountRecord(r models.CountRecord) (int64, error)
		DeleteRecord() error
	}
)

func newNoticeInterface(db *gorm.DB, g InterGormDBCli) InterNoticeRepo {
	return &NoticeRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (nr NoticeRepo) GetQuota(id string) bool {
	var (
		db     = nr.db.Model(&models.Tenant{})
		data   models.Tenant
		Number int64
	)

	db.Where("id = ?", id)
	db.Find(&data)

	nr.db.Model(&models.AlertNotice{}).Where("tenant_id = ?", id).Count(&Number)

	if Number < data.NoticeNumber {
		return true
	}

	return false
}

func (nr NoticeRepo) Get(tenantId, id string) (models.AlertNotice, error) {
	var alertNoticeData models.AlertNotice
	db := nr.db.Model(&models.AlertNotice{}).Where("tenant_id = ? AND uuid = ?", tenantId, id)
	err := db.First(&alertNoticeData).Error
	if err != nil {
		return alertNoticeData, err
	}

	return alertNoticeData, nil
}

func (nr NoticeRepo) List(tenantId, noticeTmplId, query string) ([]models.AlertNotice, error) {
	var alertNoticeObject []models.AlertNotice
	db := nr.db.Model(&models.AlertNotice{})

	if tenantId != "" {
		db.Where("tenant_id = ?", tenantId)
	}
	if noticeTmplId != "" {
		db.Where("notice_tmpl_id = ?", noticeTmplId)
	}
	if query != "" {
		db.Where("uuid LIKE ? OR name LIKE ? OR env LIKE ? OR notice_type LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	err := db.Find(&alertNoticeObject).Error
	if err != nil {
		return nil, err
	}

	return alertNoticeObject, nil
}

func (nr NoticeRepo) Create(r models.AlertNotice) error {
	err := nr.g.Create(models.AlertNotice{}, r)
	if err != nil {
		return err
	}
	return nil
}

func (nr NoticeRepo) Update(r models.AlertNotice) error {
	u := Updates{
		Table: models.AlertNotice{},
		Where: map[string]interface{}{
			"tenant_id = ?": r.TenantId,
			"uuid = ?":      r.Uuid,
		},
		Updates: r,
	}
	err := nr.g.Updates(u)
	if err != nil {
		return err
	}
	return nil
}

func (nr NoticeRepo) Delete(tenantId, id string) error {
	var (
		ruleNum1, ruleNum2 int64
		db                 = nr.db.Model(&models.AlertRule{})
	)

	db.Where("notice_id = ?", id).Count(&ruleNum1)
	db.Where("notice_group LIKE ?", "%"+id+"%").Count(&ruleNum2)

	if ruleNum1 != 0 || ruleNum2 != 0 {
		return fmt.Errorf("无法删除通知对象 %s, 因为已有告警规则绑定", id)
	}

	d := Delete{
		Table: models.AlertNotice{},
		Where: map[string]interface{}{
			"tenant_id = ?": tenantId,
			"uuid = ?":      id,
		},
	}
	err := nr.g.Delete(d)
	if err != nil {
		return err
	}
	return nil
}

// AddRecord 添加通知记录
func (nr NoticeRepo) AddRecord(r models.NoticeRecord) error {
	err := nr.g.Create(models.NoticeRecord{}, r)
	if err != nil {
		return err
	}
	return nil
}

func (nr NoticeRepo) ListRecord(tenantId, eventId, severity, status, query string, page models.Page) (models.ResponseNoticeRecords, error) {
	var (
		records []models.NoticeRecord
		count   int64
		db      = nr.db.Model(&models.NoticeRecord{})
	)

	db.Where("tenant_id = ?", tenantId)
	if eventId != "" {
		db.Where("event_id = ?", eventId)
	}
	if severity != "" {
		db.Where("severity = ?", severity)
	}
	if status != "" {
		db.Where("status = ?", status)
	}
	if query != "" {
		db.Where("rule_name LIKE ? OR alarm_msg LIKE ? OR err_msg LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	if err := db.Count(&count).Error; err != nil {
		return models.ResponseNoticeRecords{}, err
	}

	err := db.Limit(int(page.Size)).Offset(int((page.Index - 1) * page.Size)).Order("create_at DESC").Find(&records).Error
	if err != nil {
		return models.ResponseNoticeRecords{}, err
	}

	return models.ResponseNoticeRecords{
		List: records,
		Page: models.Page{
			Index: page.Index,
			Size:  page.Size,
			Total: count,
		},
	}, nil
}

func (nr NoticeRepo) CountRecord(r models.CountRecord) (int64, error) {
	var count int64
	db := nr.db.Model(&models.NoticeRecord{})
	db.Where("tenant_id = ?", r.TenantId)
	if r.Date != "" {
		db.Where("date = ?", r.Date)
	}
	if r.Severity != "" {
		db.Where("severity = ?", r.Severity)
	}
	err := db.Count(&count).Error
	if err != nil {
		return count, err
	}

	return count, nil
}

func (nr NoticeRepo) DeleteRecord() error {
	var saveDays int64 = 3600 * 24 * 7

	now := time.Now().Unix()
	startTime := now - saveDays

	del := Delete{
		Table: &models.NoticeRecord{},
		Where: map[string]interface{}{
			"create_at < ?": startTime,
		},
	}
	err := nr.g.Delete(del)
	if err != nil {
		logc.Errorf(context.Background(), err.Error())
		return err
	}
	return nil
}
