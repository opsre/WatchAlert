package eval

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/provider"

	"github.com/zeromicro/go-zero/core/logc"
)

type (
	// RecordingRuleEval 记录规则评估
	RecordingRuleEval interface {
		Submit(rule models.RecordingRule)
		Stop(ruleId string)
		Eval(ctx context.Context, rule models.RecordingRule)
		RestartAllEvals()
		StopAllEvals()
	}

	// RecordingRule 记录规则
	RecordingRule struct {
		ctx *ctx.Context
	}
)

func NewRecordingRuleEval(ctx *ctx.Context) RecordingRuleEval {
	return &RecordingRule{
		ctx: ctx,
	}
}

// Submit 提交记录规则评估任务
func (t *RecordingRule) Submit(rule models.RecordingRule) {
	t.ctx.Mux.Lock()
	defer t.ctx.Mux.Unlock()

	c, cancel := context.WithCancel(context.Background())
	t.ctx.ContextMap[rule.RuleId] = cancel
	go t.Eval(c, rule)
}

// Stop 停止记录规则评估任务
func (t *RecordingRule) Stop(ruleId string) {
	t.ctx.Mux.Lock()
	defer t.ctx.Mux.Unlock()

	if cancel, exists := t.ctx.ContextMap[ruleId]; exists {
		cancel()
		delete(t.ctx.ContextMap, ruleId)
	}
}

// Eval 执行记录规则评估
func (t *RecordingRule) Eval(ctx context.Context, rule models.RecordingRule) {
	err := rule.Validate()
	if err != nil {
		logc.Errorf(t.ctx.Ctx, "Recording rule validation failed, MetricName: %s, RuleId: %s, Error: %v", rule.MetricName, rule.RuleId, err)
		return
	}

	taskChan := make(chan struct{}, TaskChannelBufferSize)
	timer := time.NewTicker(t.getEvalTimeDuration(rule.EvalInterval))
	defer func() {
		timer.Stop()
		if r := recover(); r != nil {
			// 获取调用栈信息
			stack := debug.Stack()
			logc.Errorf(t.ctx.Ctx, "Recovered from recording rule eval goroutine panic: %s, MetricName: %s, RuleId: %s\n%s", r, rule.MetricName, rule.RuleId, stack)
			t.Restart(rule)
		}
	}()

	for {
		select {
		case <-timer.C:
			// 处理任务信号量
			taskChan <- struct{}{}
			logc.Infof(t.ctx.Ctx, fmt.Sprintf("Handle recoding eval task, RuleId: %v, MetricName: %s", rule.RuleId, rule.MetricName))
			t.executeTask(rule, taskChan)
		case <-ctx.Done():
			logc.Infof(t.ctx.Ctx, fmt.Sprintf("Stop recoding eval task, RuleId: %v, MetricName: %s", rule.RuleId, rule.MetricName))
			return
		}
		timer.Reset(t.getEvalTimeDuration(rule.EvalInterval))
	}
}

// executeTask 执行评估任务
func (t *RecordingRule) executeTask(rule models.RecordingRule, taskChan chan struct{}) {
	defer func() {
		// 释放任务信号量
		<-taskChan
	}()

	// 在规则评估前检查是否仍然启用
	if !t.isRuleEnabled(rule.RuleId) {
		return
	}

	// 处理数据源
	t.processSingleDatasource(rule)
}

// processSingleDatasource 处理数据源
func (t *RecordingRule) processSingleDatasource(rule models.RecordingRule) {
	instance, err := t.ctx.DB.Datasource().GetInstance(rule.DatasourceId)
	if err != nil {
		logc.Errorf(t.ctx.Ctx, "Failed to get datasource instance %s: %v", rule.DatasourceId, err)
		return
	}

	// 检查数据源健康状态
	if ok, _ := provider.CheckDatasourceHealth(instance); !ok {
		logc.Errorf(t.ctx.Ctx, "Datasource %s is unhealthy", rule.DatasourceId)
		return
	}

	// 检查数据源是否启用
	if !*instance.Enabled {
		logc.Errorf(t.ctx.Ctx, "Datasource %s is disabled", rule.DatasourceId)
		return
	}

	// 调用处理器
	switch rule.DatasourceType {
	default:
		t.processPrometheus(rule)
	}
}

// processPrometheus 处理 Prometheus 数据源
func (t *RecordingRule) processPrometheus(rule models.RecordingRule) {
	instance, err := t.ctx.DB.Datasource().GetInstance(rule.DatasourceId)
	if err != nil {
		logc.Errorf(t.ctx.Ctx, "Failed to get datasource instance %s: %v", rule.DatasourceId, err)
		return
	}

	// 创建 Prometheus 客户端
	cli, err := provider.NewPrometheusClient(instance)
	if err != nil {
		logc.Errorf(t.ctx.Ctx, "Failed to create Prometheus client: %v", err)
		return
	}

	// 执行 PromQL 查询
	results, err := cli.Query(rule.PromQL)
	if err != nil {
		logc.Errorf(t.ctx.Ctx, "Failed to execute PromQL query: %v", err)
		return
	}

	newResults := make([]provider.Metrics, len(results))
	for i, result := range results {
		newResults[i] = provider.Metrics{
			Name:   rule.MetricName,
			Help:   result.Help,
			Labels: result.Labels,
			Value:  result.Value,
		}
		delete(newResults[i].Labels, "__name__")
	}

	// 将结果写入 Prometheus 远程写入端点
	err = cli.Write(context.Background(), newResults, rule.Labels)
	if err != nil {
		logc.Errorf(t.ctx.Ctx, "Failed to write recording rule result: %v", err)
		return
	}
}

// Restart 重启记录规则评估
func (t *RecordingRule) Restart(rule models.RecordingRule) {
	t.Stop(rule.RuleId)
	t.Submit(rule)
}

// RestartAllEvals 重启所有记录规则评估器
func (t *RecordingRule) RestartAllEvals() {
	ruleList, err := t.getRuleList()
	if err != nil {
		logc.Error(t.ctx.Ctx, fmt.Sprintf("Failed to get recording rule list: %v", err))
		return
	}

	count := len(ruleList)
	if count == 0 {
		return
	}

	logc.Info(t.ctx.Ctx, fmt.Sprintf("获取到 %d 个状态为启用的记录规则", count))

	// 使用工作池限制并发数量
	const maxWorkers = 10
	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, maxWorkers)

	wg.Add(count)
	for _, rule := range ruleList {
		rule := rule
		go func() {
			semaphore <- struct{}{}
			defer func() {
				wg.Done()
				<-semaphore
			}()

			t.Submit(rule)
		}()
	}

	wg.Wait()
	close(semaphore)
	logc.Info(t.ctx.Ctx, "所有记录规则评估器启动成功！")
}

// isRuleEnabled 检查记录规则是否启用
func (t *RecordingRule) isRuleEnabled(ruleId string) bool {
	// 直接检查数据库或缓存中的当前启用状态
	e := t.ctx.DB.RecordingRule().GetRuleObject(ruleId).Enabled
	if e == nil {
		return false
	}

	return *e
}

// getRuleList 获取记录规则列表
func (t *RecordingRule) getRuleList() ([]models.RecordingRule, error) {
	var ruleList []models.RecordingRule
	if err := t.ctx.DB.DB().Where("enabled = ?", "1").Find(&ruleList).Error; err != nil {
		return nil, fmt.Errorf("获取 Recording Rule List 失败: %w", err)
	}
	return ruleList, nil
}

// StopAllEvals 停止所有记录规则评估器
func (t *RecordingRule) StopAllEvals() {
	t.ctx.Mux.Lock()
	defer t.ctx.Mux.Unlock()

	count := len(t.ctx.ContextMap)
	if count == 0 {
		return
	}

	logc.Infof(t.ctx.Ctx, "停止 %d 个记录规则评估器...", count)

	// 取消所有评估任务
	for ruleId, cancel := range t.ctx.ContextMap {
		cancel()
		delete(t.ctx.ContextMap, ruleId)
	}

	logc.Infof(t.ctx.Ctx, "所有记录规则评估器已停止")
}

// getEvalTimeDuration 获取评估时间间隔
func (t *RecordingRule) getEvalTimeDuration(evalInterval int64) time.Duration {
	return time.Duration(evalInterval) * time.Second
}
