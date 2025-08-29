package repo

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
	"watchAlert/internal/models"
)

type (
	DutyCalendarRepo struct {
		entryRepo
	}

	InterDutyCalendar interface {
		GetCalendarInfo(dutyId, time string) models.DutySchedule
		GetDutyUserInfo(dutyId, time string) (models.Member, bool)
		Create(r models.DutySchedule) error
		Update(r models.DutySchedule) error
		Search(r models.DutyScheduleQuery) ([]models.DutySchedule, error)
		GetCalendarUsers(r models.DutyScheduleQuery) ([]models.Users, error)
	}
)

func newDutyCalendarInterface(db *gorm.DB, g InterGormDBCli) InterDutyCalendar {
	return &DutyCalendarRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

// GetCalendarInfo 获取值班表信息
func (dc DutyCalendarRepo) GetCalendarInfo(dutyId, time string) models.DutySchedule {
	var dutySchedule models.DutySchedule

	dc.db.Model(models.DutySchedule{}).
		Where("duty_id = ? AND time = ?", dutyId, time).
		First(&dutySchedule)

	return dutySchedule
}

// GetDutyUserInfo 获取值班用户信息
func (dc DutyCalendarRepo) GetDutyUserInfo(dutyId, time string) (models.Member, bool) {
	var user models.Member
	schedule := dc.GetCalendarInfo(dutyId, time)
	db := dc.db.Model(models.Member{}).
		Where("user_id = ?", schedule.UserId)
	if err := db.First(&user).Error; err != nil {
		return user, false
	}
	if user.JoinDuty == "true" {
		return user, true
	}

	return user, false
}

func (dc DutyCalendarRepo) Create(r models.DutySchedule) error {
	err := dc.g.Create(models.DutySchedule{}, r)
	if err != nil {
		return err
	}
	return nil
}

func (dc DutyCalendarRepo) Update(r models.DutySchedule) error {
	u := Updates{
		Table: models.DutySchedule{},
		Where: map[string]interface{}{
			"tenant_id = ?": r.TenantId,
			"duty_id = ?":   r.DutyId,
			"time = ?":      r.Time,
		},
		Updates: r,
	}

	err := dc.g.Updates(u)
	if err != nil {
		return err
	}
	return nil
}

func (dc DutyCalendarRepo) Search(r models.DutyScheduleQuery) ([]models.DutySchedule, error) {
	var dutyScheduleList []models.DutySchedule
	db := dc.db.Model(&models.DutySchedule{})

	if r.Time != "" {
		db.Where("tenant_id = ? AND duty_id = ? AND time = ?", r.TenantId, r.DutyId, r.Time).Find(&dutyScheduleList)
		return dutyScheduleList, nil
	}

	yearMonth := fmt.Sprintf("%v-%v-", r.Year, r.Month)
	db.Where("tenant_id = ? AND duty_id = ? AND time LIKE ?", r.TenantId, r.DutyId, yearMonth+"%")
	err := db.Find(&dutyScheduleList).Error
	if err != nil {
		return dutyScheduleList, err
	}

	return dutyScheduleList, nil
}

// GetCalendarUsers 获取值班用户
// 只获取当前月份到月底正在值班的用户，避免已经移除过的用户仍存在值班用户列表当中；
func (dc DutyCalendarRepo) GetCalendarUsers(r models.DutyScheduleQuery) ([]models.Users, error) {
	var users []models.Users

	// 获取当前年月日
	now := time.Now().UTC()
	currentDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	// 计算当年12月31日
	endOfYear := time.Date(now.Year(), 12, 31, 0, 0, 0, 0, time.UTC)

	db := dc.db.Model(&models.DutySchedule{})
	db.Where("tenant_id = ? AND duty_id = ?", r.TenantId, r.DutyId)
	db.Where("time >= ? AND time <= ?", currentDate, endOfYear)
	db.Distinct("user_id", "username")
	if err := db.Find(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get calendar users: %w", err)
	}

	return users, nil
}
