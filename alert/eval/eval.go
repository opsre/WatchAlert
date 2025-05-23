package eval

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logc"
	"runtime/debug"
	"strings"
	"time"
	"watchAlert/alert/process"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/provider"
	"watchAlert/pkg/tools"

	"golang.org/x/sync/errgroup"
)

type (
	// AlertRuleEval 告警规则评估
	AlertRuleEval interface {
		Submit(rule models.AlertRule)
		Stop(ruleId string)
		Eval(ctx context.Context, rule models.AlertRule)
		Recover(tenantId, ruleId string, eventCacheKey models.AlertEventCacheKey, faultCenterInfoKey models.FaultCenterInfoCacheKey, curFingerprints []string)
		RestartAllEvals()
	}

	// AlertRule 告警规则
	AlertRule struct {
		ctx         *ctx.Context
		watchCtxMap map[string]context.CancelFunc
	}
)

func NewAlertRuleEval(ctx *ctx.Context) AlertRuleEval {
	return &AlertRule{
		ctx:         ctx,
		watchCtxMap: make(map[string]context.CancelFunc),
	}
}

func (t *AlertRule) Submit(rule models.AlertRule) {
	t.ctx.Mux.Lock()
	defer t.ctx.Mux.Unlock()

	c, cancel := context.WithCancel(context.Background())
	t.watchCtxMap[rule.RuleId] = cancel
	go t.Eval(c, rule)
}

func (t *AlertRule) Stop(ruleId string) {
	t.ctx.Mux.Lock()
	defer t.ctx.Mux.Unlock()

	if cancel, exists := t.watchCtxMap[ruleId]; exists {
		cancel()
		delete(t.watchCtxMap, ruleId)
	}
}

func (t *AlertRule) Restart(rule models.AlertRule) {
	t.Stop(rule.RuleId)
	t.Submit(rule)
}

func (t *AlertRule) Eval(ctx context.Context, rule models.AlertRule) {
	timer := time.NewTicker(t.getEvalTimeDuration(rule.EvalTimeType, rule.EvalInterval))
	defer func() {
		timer.Stop()
		if r := recover(); r != nil {
			// 获取调用栈信息
			stack := debug.Stack()
			logc.Error(t.ctx.Ctx, fmt.Sprintf("Recovered from rule eval goroutine panic: %s, RuleName: %s, RuleId: %s\n%s", r, rule.RuleName, rule.RuleId, stack))
			t.Restart(rule)
		}
	}()

	for {
		select {
		case <-timer.C:
			// 在规则评估前检查是否仍然启用，避免不必要的操作
			if !t.isRuleEnabled(rule.RuleId) {
				return
			}

			var curFingerprints []string
			for _, dsId := range rule.DatasourceIdList {
				instance, err := t.ctx.DB.Datasource().GetInstance(dsId)
				if err != nil {
					logc.Error(t.ctx.Ctx, err.Error())
					continue
				}

				ok, _ := provider.CheckDatasourceHealth(instance)
				if !ok {
					continue
				}

				var fingerprints []string

				switch rule.DatasourceType {
				case "Prometheus", "VictoriaMetrics":
					fingerprints = metrics(t.ctx, dsId, instance.Type, rule)
				case "AliCloudSLS", "Loki", "ElasticSearch", "VictoriaLogs":
					fingerprints = logs(t.ctx, dsId, instance.Type, rule)
				case "Jaeger":
					fingerprints = traces(t.ctx, dsId, instance.Type, rule)
				case "CloudWatch":
					fingerprints = cloudWatch(t.ctx, dsId, rule)
				case "KubernetesEvent":
					fingerprints = kubernetesEvent(t.ctx, dsId, rule)
				default:
					continue
				}
				// 追加当前数据源的指纹到总列表
				curFingerprints = append(curFingerprints, fingerprints...)
			}
			logc.Infof(t.ctx.Ctx, fmt.Sprintf("规则评估 -> %v", tools.JsonMarshal(rule)))
			t.Recover(rule.TenantId, rule.RuleId, models.BuildAlertEventCacheKey(rule.TenantId, rule.FaultCenterId), models.BuildFaultCenterInfoCacheKey(rule.TenantId, rule.FaultCenterId), curFingerprints)
			t.GC(t.ctx, rule, curFingerprints)

		case <-ctx.Done():
			logc.Infof(t.ctx.Ctx, fmt.Sprintf("停止 RuleId: %v, RuleName: %s 的 Watch 协程", rule.RuleId, rule.RuleName))
			return
		}
		timer.Reset(t.getEvalTimeDuration(rule.EvalTimeType, rule.EvalInterval))
	}
}

// getEvalTimeDuration 获取评估时间
func (t *AlertRule) getEvalTimeDuration(evalTimeType string, evalInterval int64) time.Duration {
	switch evalTimeType {
	case "millisecond":
		return time.Millisecond * time.Duration(evalInterval)
	default:
		return time.Second * time.Duration(evalInterval)
	}
}

func (t *AlertRule) Recover(tenantId, ruleId string, eventCacheKey models.AlertEventCacheKey, faultCenterInfoKey models.FaultCenterInfoCacheKey, curFingerprints []string) {
	// 获取所有的故障中心告警事件
	events, err := t.ctx.Redis.Alert().GetAllEvents(eventCacheKey)
	if err != nil {
		logc.Errorf(t.ctx.Ctx, "Failed to get all events: %v", err)
		return
	}

	// 提取当前规则相关的指纹
	var fingerprints []string
	for fingerprint, event := range events {
		if strings.Contains(event.RuleId, ruleId) {
			fingerprints = append(fingerprints, fingerprint)
		}
	}

	// 计算需要恢复的指纹列表
	recoverFingerprints := tools.GetSliceDifference(fingerprints, curFingerprints)
	if len(recoverFingerprints) == 0 {
		t.handlePendingRecovery(tenantId, ruleId, events)
		return
	}

	// 处理恢复逻辑
	curTime := time.Now().Unix()
	recoverWaitTime := t.getRecoverWaitTime(faultCenterInfoKey)

	for _, fingerprint := range recoverFingerprints {
		event, ok := events[fingerprint]
		if !ok {
			continue
		}

		// 获取待恢复状态的时间戳
		wTime, err := t.ctx.Redis.PendingRecover().Get(tenantId, ruleId, fingerprint)
		if err == redis.Nil {
			// 如果不存在，则记录当前时间
			t.ctx.Redis.PendingRecover().Set(tenantId, ruleId, fingerprint, curTime)
			continue
		} else if err != nil {
			logc.Errorf(t.ctx.Ctx, "Failed to get pending recovery time for fingerprint %s: %v", fingerprint, err)
			continue
		}

		// 判断是否在等待时间内
		recoverThreshold := wTime + int64(recoverWaitTime)*60
		if recoverThreshold >= curTime {
			// 进入待恢复状态
			if err := event.TransitionStatus(models.StatePendingRecovery); err != nil {
				logc.Errorf(t.ctx.Ctx, "Failed to transition to pending recovery state for fingerprint %s: %v", fingerprint, err)
				continue
			}
		} else if curTime >= recoverThreshold && event.Status == models.StatePendingRecovery { // 当前时间超过预期等待时间，并且状态是 PendingRecovery 时才执行恢复逻辑
			// 已恢复状态
			if err := event.TransitionStatus(models.StateRecovered); err != nil {
				logc.Errorf(t.ctx.Ctx, "Failed to transition to recovered state for fingerprint %s: %v", fingerprint, err)
				continue
			}
			t.ctx.Redis.PendingRecover().Delete(tenantId, ruleId, fingerprint)
		}

		// 更新告警事件
		t.ctx.Redis.Alert().PushAlertEvent(event)
	}
}

// 处理待恢复状态的事件
func (t *AlertRule) handlePendingRecovery(tenantId, ruleId string, events map[string]*models.AlertCurEvent) {
	fs := t.ctx.Redis.PendingRecover().List(tenantId, ruleId)
	for fingerprint := range fs {
		event, ok := events[fingerprint]
		if !ok {
			t.ctx.Redis.PendingRecover().Delete(tenantId, ruleId, fingerprint)
			continue
		}

		// 转换为告警状态
		if err := event.TransitionStatus(models.StateAlerting); err != nil {
			logc.Errorf(t.ctx.Ctx, "Failed to transition to alerting state for fingerprint %s: %v", fingerprint, err)
			continue
		}

		// 更新告警事件并删除待恢复状态
		t.ctx.Redis.Alert().PushAlertEvent(event)
		t.ctx.Redis.PendingRecover().Delete(tenantId, ruleId, fingerprint)
	}
}

// 获取恢复等待时间
func (t *AlertRule) getRecoverWaitTime(faultCenterInfoKey models.FaultCenterInfoCacheKey) int64 {
	faultCenter := t.ctx.Redis.FaultCenter().GetFaultCenterInfo(faultCenterInfoKey)
	if faultCenter.RecoverWaitTime == 0 {
		return 1
	}
	return faultCenter.RecoverWaitTime
}

func (t *AlertRule) GC(ctx *ctx.Context, rule models.AlertRule, curFingerprints []string) {
	go process.GcRecoverWaitCache(ctx, rule, curFingerprints)
}

func (t *AlertRule) RestartAllEvals() {
	ruleList, err := t.getRuleList()
	if err != nil {
		logc.Error(t.ctx.Ctx, err.Error())
		return
	}

	g := new(errgroup.Group)
	for _, rule := range ruleList {
		rule := rule
		g.Go(func() error {
			t.Submit(rule)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		logc.Error(t.ctx.Ctx, err.Error())
	}
}

// isRuleEnabled 检查规则是否启用
func (t *AlertRule) isRuleEnabled(ruleId string) bool {
	// 直接检查数据库或缓存中的当前启用状态
	return *t.ctx.DB.Rule().GetRuleObject(ruleId).Enabled
}

func (t *AlertRule) getRuleList() ([]models.AlertRule, error) {
	var ruleList []models.AlertRule
	if err := t.ctx.DB.DB().Where("enabled = ?", "1").Find(&ruleList).Error; err != nil {
		return ruleList, fmt.Errorf("获取 Rule List 失败, err: %s", err.Error())
	}
	return ruleList, nil
}
