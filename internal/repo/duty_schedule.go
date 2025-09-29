package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
	"time"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

type (
	DutyCalendarRepo struct {
		entryRepo
	}

	InterDutyCalendar interface {
		GetCalendarInfo(dutyId, time string) models.DutySchedule
		GetDutyUserInfo(dutyId, time string) ([]models.Member, bool)
		Create(r models.DutySchedule) error
		Update(r models.DutySchedule) error
		Search(tenantId, dutyId, time string) ([]models.DutySchedule, error)
		GetCalendarUsers(tenantId, dutyId string) ([][]models.DutyUser, error)
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
func (dc DutyCalendarRepo) GetDutyUserInfo(dutyId, time string) ([]models.Member, bool) {
	var users []models.Member
	schedule := dc.GetCalendarInfo(dutyId, time)
	for _, user := range schedule.Users {
		var userData models.Member
		db := dc.db.Model(models.Member{}).Where("user_id = ?", user.UserId)
		if err := db.First(&userData).Error; err != nil {
			logc.Error(context.Background(), "获取值班用户信息失败, msg: "+err.Error())
			continue
		}
		users = append(users, userData)
	}

	if users == nil {
		return users, false
	}

	return users, true
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

func (dc DutyCalendarRepo) Search(tenantId, dutyId, time string) ([]models.DutySchedule, error) {
	var dutyScheduleList []models.DutySchedule
	db := dc.db.Model(&models.DutySchedule{})

	db.Where("tenant_id = ? AND duty_id = ? AND time LIKE ?", tenantId, dutyId, time+"%")
	err := db.Find(&dutyScheduleList).Error
	if err != nil {
		return dutyScheduleList, err
	}

	return dutyScheduleList, nil
}

// GetCalendarUsers 获取值班用户
// 只获取当前月份到月底正在值班的用户，避免已经移除过的用户仍存在值班用户列表当中；
func (dc DutyCalendarRepo) GetCalendarUsers(tenantId, dutyId string) ([][]models.DutyUser, error) {
	var (
		entries      []models.DutySchedule
		groupedUsers [][]models.DutyUser
	)

	// 获取当前年月日
	now := time.Now().UTC()
	currentDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	// 计算当年12月31日
	endOfYear := time.Date(now.Year(), 12, 31, 0, 0, 0, 0, time.UTC)

	db := dc.db.Model(&models.DutySchedule{})
	db.Where("tenant_id = ? AND duty_id = ? AND status = ?", tenantId, dutyId, models.CalendarFormalStatus)
	db.Where("time >= ? AND time <= ?", currentDate, endOfYear)

	if err := db.Find(&entries).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get calendar users: %w", err)
	}

	user := make(map[string]struct{})
	for _, entry := range entries {
		key := tools.JsonMarshalToString(entry.Users)
		if _, ok := user[key]; ok {
			continue
		}

		groupedUsers = append(groupedUsers, entry.Users)
		user[key] = struct{}{}
	}

	return groupedUsers, nil
}
