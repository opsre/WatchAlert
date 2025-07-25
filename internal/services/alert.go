package services

import (
	"github.com/zeromicro/go-zero/core/logc"
	"sync"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
)

type (
	alertService struct {
		ctx *ctx.Context
	}

	InterAlertService interface {
		RePushRule(ctx *ctx.Context, rule chan *models.AlertRule)
	}
)

func newInterAlertService(ctx *ctx.Context) InterAlertService {
	return &alertService{
		ctx: ctx,
	}
}

func (as alertService) RePushRule(ctx *ctx.Context, alertRule chan *models.AlertRule) {

	var (
		ruleList []models.AlertRule
		// 创建一个通道用于接收处理结果
		resultCh = make(chan error)
		// 使用 WaitGroup 来等待所有规则的处理完成
		wg sync.WaitGroup
	)
	ctx.DB.DB().Where("enabled = ?", "1").Find(&ruleList)

	// 并发处理规则
	for _, rule := range ruleList {
		wg.Add(1)
		go func(rule models.AlertRule) {
			defer wg.Done()

			alertRule <- &rule

			resultCh <- nil
		}(rule)
	}

	// 等待所有规则的处理完成
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// 处理结果
	for result := range resultCh {
		if result != nil {
			logc.Error(ctx.Ctx, result)
		}
	}

}
