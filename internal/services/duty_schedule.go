package services

import (
	"context"
	"fmt"
	"sync"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"

	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logc"
	"golang.org/x/sync/errgroup"
)

type dutyCalendarService struct {
	ctx *ctx.Context
}

type InterDutyCalendarService interface {
	CreateAndUpdate(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Search(req interface{}) (interface{}, interface{})
	GetCalendarUsers(req interface{}) (interface{}, interface{})
	GenerateNextYearScheduleCronjob(ctx context.Context)
}

func newInterDutyCalendarService(ctx *ctx.Context) InterDutyCalendarService {
	return &dutyCalendarService{
		ctx: ctx,
	}
}

// CreateAndUpdate 创建和更新值班表
func (dms dutyCalendarService) CreateAndUpdate(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestDutyCalendarCreate)
	curYear, curMonth, _ := tools.ParseTime(r.Month)

	var data = models.DutyCalendarInfo{
		TenantId:   r.TenantId,
		DutyId:     r.DutyId,
		DutyPeriod: r.DutyPeriod,
		Month:      r.Month,
		UserGroup:  r.UserGroup,
		DateType:   r.DateType,
	}

	_, err := dms.ctx.DB.DutyCalendar().GetCalendarInfo(r.DutyId)
	if err != nil {
		err := dms.ctx.DB.DutyCalendar().CreateCalendarInfo(data)
		if err != nil {
			logc.Errorf(dms.ctx.Ctx, "创建值班信息表失败: %v", err)
			return nil, err
		}
	} else {
		err := dms.ctx.DB.DutyCalendar().UpdateCalendarInfo(data)
		if err != nil {
			logc.Errorf(dms.ctx.Ctx, "更新值班信息表失败: %v", err)
			return nil, err
		}
	}

	dutyScheduleList, err := dms.generateDutySchedule(*r, curYear, curMonth)
	if err != nil {
		return nil, fmt.Errorf("生成值班表数据失败: %w", err)
	}

	if err := dms.updateDutyScheduleInDB(dutyScheduleList, r.TenantId); err != nil {
		logc.Errorf(dms.ctx.Ctx, "值班表数据入库失败: %w", err)
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

func (dms dutyCalendarService) generateDutySchedule(dutyInfo types.RequestDutyCalendarCreate, year int, month time.Month) ([]models.DutySchedule, error) {
	dutyDays := dms.calculateDutyDays(dutyInfo.DateType, dutyInfo.DutyPeriod)
	timeC := dms.generateDutyDates(year, month)
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
	case "month":
		return 31 * dutyPeriod
	case "year":
		return 365 * dutyPeriod
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

	for {
		// 数据消费完成后退出
		if len(timeC) == 0 {
			break
		}

		// 每个用户组独立计数，循环开始时重置
		var count int

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
		dutyScheduleInfo := dms.ctx.DB.DutyCalendar().GetCalendarData(schedule.DutyId, schedule.Time)

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

// GenerateNextYearScheduleCronjob 每年12月1日生成下一年的值班表
func (dms dutyCalendarService) GenerateNextYearScheduleCronjob(ctx context.Context) {
	c := cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
	))

	// 每年12月1日 00:00 执行
	entryID, err := c.AddFunc("0 0 1 12 *", func() {
		taskCtx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		logc.Infof(taskCtx, "触发年度值班表生成任务")
		dms.generateNextYearSchedule(taskCtx)
	})

	if err != nil {
		logc.Error(ctx, "添加年度值班表生成定时任务失败: %v", err)
		return
	}

	logc.Infof(ctx, "启动年度值班表生成定时任务, EntryID: %d, Cron: 0 0 1 12 *", entryID)

	c.Start()

	<-ctx.Done()
	logc.Infof(ctx, "停止年度值班表生成定时任务")
	c.Stop()
}

// generateNextYearSchedule 根据历史数据生成下一年的值班表
func (dms dutyCalendarService) generateNextYearSchedule(ctx context.Context) {
	now := time.Now().UTC()
	currentYear := now.Year()
	nextYear := currentYear + 1

	// 获取所有租户
	tenants, err := dms.ctx.DB.Tenant().List("")
	if err != nil {
		logc.Errorf(ctx, "获取租户列表失败: %s", err.Error())
		return
	}

	g, subCtx := errgroup.WithContext(ctx)
	// 最多同时处理 5 个租户
	g.SetLimit(5)

	for _, tenant := range tenants {
		tenant := tenant
		g.Go(func() error {
			infos, err := dms.ctx.DB.DutyCalendar().ListCalendarInfo(tenant.ID)
			if err != nil {
				logc.Errorf(subCtx, "租户 %s 获取当前年份值班表失败: %s", tenant.ID, err.Error())
				return nil
			}

			for _, info := range infos {
				dutyScheduleList, err := dms.generateDutySchedule(types.RequestDutyCalendarCreate{
					TenantId:   info.TenantId,
					DutyId:     info.DutyId,
					DutyPeriod: info.DutyPeriod,
					Month:      info.Month,
					UserGroup:  info.UserGroup,
					DateType:   info.DateType,
				}, nextYear, time.January)
				if err != nil {
					logc.Errorf(subCtx, "生成值班表数据失败, TenantId: %s, DutyId: %s, Err: %w", info.TenantId, info.DutyId, err)
					continue
				}

				if err := dms.updateDutyScheduleInDB(dutyScheduleList, info.TenantId); err != nil {
					logc.Errorf(subCtx, "值班表数据入库失败, TenantId: %s, DutyId: %s, Err: %w", info.TenantId, info.DutyId, err)
				}
			}
			return nil
		})
	}

	_ = g.Wait()

	logc.Infof(ctx, "年度值班表生成任务完成")
}
