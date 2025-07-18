package repo

import (
	"fmt"
	"gorm.io/gorm"
	"time"
	models "watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

type (
	DutyRepo struct {
		entryRepo
	}
	InterDutyRepo interface {
		GetQuota(id string) bool
		List(r models.DutyManagementQuery) ([]models.DutyManagement, error)
		Create(r models.DutyManagement) error
		Update(r models.DutyManagement) error
		Delete(r models.DutyManagementQuery) error
		Get(r models.DutyManagementQuery) (models.DutyManagement, error)
	}
)

func newDutyInterface(db *gorm.DB, g InterGormDBCli) InterDutyRepo {
	return &DutyRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (d DutyRepo) GetQuota(id string) bool {
	var (
		db     = d.DB().Model(&models.Tenant{})
		data   models.Tenant
		Number int64
	)

	db.Where("id = ?", id)
	db.Find(&data)

	d.DB().Model(&models.DutyManagement{}).Where("tenant_id = ?", id).Count(&Number)

	if Number < data.DutyNumber {
		return true
	}

	return false
}

func (d DutyRepo) List(r models.DutyManagementQuery) ([]models.DutyManagement, error) {
	var data []models.DutyManagement

	db := d.db.Model(&models.DutyManagement{})
	db.Where("tenant_id = ?", r.TenantId)
	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}

	for index, value := range data {
		var dutySchedule models.DutySchedule
		d.DB().Model(models.DutySchedule{}).Where("duty_id = ? and time = ?", value.ID, time.Now().Format("2006-1-2")).Find(&dutySchedule)
		data[index].CurDutyUser = dutySchedule.Users
	}

	return data, nil
}

func (d DutyRepo) Create(r models.DutyManagement) error {
	nr := r
	nr.ID = "dt-" + tools.RandId()
	nr.CreateAt = time.Now().Unix()
	err := d.g.Create(&models.DutyManagement{}, nr)
	if err != nil {
		return err
	}
	return nil
}

func (d DutyRepo) Update(r models.DutyManagement) error {
	u := Updates{
		Table: models.DutyManagement{},
		Where: map[string]interface{}{
			"tenant_id = ?": r.TenantId,
			"id = ?":        r.ID,
		},
		Updates: r,
	}
	err := d.g.Updates(u)
	if err != nil {
		return err
	}
	return nil
}

func (d DutyRepo) Delete(r models.DutyManagementQuery) error {
	var noticeNum int64
	db := d.db.Model(&models.AlertNotice{})
	db.Where("tenant_id = ? AND duty_id = ?", r.TenantId, r.ID).Count(&noticeNum)
	if noticeNum != 0 {
		return fmt.Errorf("无法删除值班表 %s, 因为已有通知对象绑定", r.ID)
	}

	delDuty := Delete{
		Table: models.DutyManagement{},
		Where: map[string]interface{}{
			"tenant_id = ?": r.TenantId,
			"id = ?":        r.ID,
		},
	}
	err := d.g.Delete(delDuty)
	if err != nil {
		return err
	}

	delCalendar := Delete{
		Table: models.DutySchedule{},
		Where: map[string]interface{}{
			"tenant_id = ?": r.TenantId,
			"duty_id = ?":   r.ID,
		},
	}
	err = d.g.Delete(delCalendar)
	if err != nil {
		return err
	}

	return nil
}

func (d DutyRepo) Get(r models.DutyManagementQuery) (models.DutyManagement, error) {
	var data models.DutyManagement
	db := d.db.Model(&models.DutyManagement{})
	db.Where("tenant_id = ? AND id = ?", r.TenantId, r.ID)
	err := db.First(&data).Error
	if err != nil {
		return models.DutyManagement{}, err
	}

	return data, nil
}
