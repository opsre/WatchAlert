package services

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"sync"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"
)

type dutyCalendarService struct {
	ctx *ctx.Context
}

type InterDutyCalendarService interface {
	CreateAndUpdate(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Search(req interface{}) (interface{}, interface{})
	GetCalendarUsers(req interface{}) (interface{}, interface{})
}

func newInterDutyCalendarService(ctx *ctx.Context) InterDutyCalendarService {
	return &dutyCalendarService{
		ctx: ctx,
	}
}

// CreateAndUpdate 创建和更新值班表
func (dms dutyCalendarService) CreateAndUpdate(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestDutyCalendarCreate)
	dutyScheduleList, err := dms.generateDutySchedule(*r)
	if err != nil {
		return nil, fmt.Errorf("生成值班表失败: %w", err)
	}

	if err := dms.updateDutyScheduleInDB(dutyScheduleList, r.TenantId); err != nil {
		logc.Errorf(dms.ctx.Ctx, err.Error())
	}
	return nil, nil
}

// Update 更新值班表
func (dms dutyCalendarService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestDutyCalendarUpdate)
	err := dms.ctx.DB.DutyCalendar().Update(models.DutySchedule{
		TenantId: r.TenantId,
		DutyId:   r.DutyId,
		Time:     r.Time,
		Status:   r.Status,
		Users:    r.Users,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Search 查询值班表
func (dms dutyCalendarService) Search(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestDutyCalendarQuery)
	data, err := dms.ctx.DB.DutyCalendar().Search(r.TenantId, r.DutyId, r.Time)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (dms dutyCalendarService) GetCalendarUsers(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestDutyCalendarQuery)
	data, err := dms.ctx.DB.DutyCalendar().GetCalendarUsers(r.TenantId, r.DutyId)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (dms dutyCalendarService) generateDutySchedule(dutyInfo types.RequestDutyCalendarCreate) ([]models.DutySchedule, error) {
	curYear, curMonth, _ := tools.ParseTime(dutyInfo.Month)
	dutyDays := dms.calculateDutyDays(dutyInfo.DateType, dutyInfo.DutyPeriod)
	timeC := dms.generateDutyDates(curYear, curMonth)
	dutyScheduleList := dms.createDutyScheduleList(dutyInfo, timeC, dutyDays)

	return dutyScheduleList, nil
}

// 计算值班天数
func (dms dutyCalendarService) calculateDutyDays(dateType string, dutyPeriod int) int {
	switch dateType {
	case "day":
		return dutyPeriod
	case "week":
		return 7 * dutyPeriod
	default:
		return 0
	}
}

// 生成值班日期
func (dms dutyCalendarService) generateDutyDates(year int, startMonth time.Month) <-chan string {
	timeC := make(chan string, 370)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer close(timeC)
		defer wg.Done()
		for month := startMonth; month <= 12; month++ {
			daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
			for day := 1; day <= daysInMonth; day++ {
				date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
				if date.Month() != month {
					break
				}
				timeC <- date.Format("2006-1-2")
			}
		}
	}()

	// 等待所有日期生产完成
	wg.Wait()
	return timeC
}

// 创建值班表
func (dms dutyCalendarService) createDutyScheduleList(dutyInfo types.RequestDutyCalendarCreate, timeC <-chan string, dutyDays int) []models.DutySchedule {
	var dutyScheduleList []models.DutySchedule
	var count int

	for {
		// 数据消费完成后退出
		if len(timeC) == 0 {
			break
		}

		for _, users := range dutyInfo.UserGroup {
			for day := 1; day <= dutyDays; day++ {
				date, ok := <-timeC
				if !ok {
					return dutyScheduleList
				}

				dutyScheduleList = append(dutyScheduleList, models.DutySchedule{
					DutyId: dutyInfo.DutyId,
					Time:   date,
					Users:  users,
					Status: dutyInfo.Status,
				})

				if dutyInfo.DateType == "week" && tools.IsEndOfWeek(date) {
					count++
					if count == dutyInfo.DutyPeriod {
						count = 0
						break
					}
				}
			}
		}
	}

	return dutyScheduleList
}

// 更新库表
func (dms dutyCalendarService) updateDutyScheduleInDB(dutyScheduleList []models.DutySchedule, tenantId string) error {
	for _, schedule := range dutyScheduleList {
		schedule.TenantId = tenantId
		dutyScheduleInfo := dms.ctx.DB.DutyCalendar().GetCalendarInfo(schedule.DutyId, schedule.Time)

		var err error
		if dutyScheduleInfo.Time != "" {
			err = dms.ctx.DB.DutyCalendar().Update(schedule)
		} else {
			err = dms.ctx.DB.DutyCalendar().Create(schedule)
		}

		if err != nil {
			return fmt.Errorf("更新/创建值班系统失败: %w", err)
		}
	}
	return nil
}
