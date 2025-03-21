package eval

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"runtime/debug"
	"strings"
	"time"
	"watchAlert/alert/process"
	"watchAlert/alert/storage"
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
		Recover(ruleId, faultCenterKey string, faultCenterInfoKey string, curKeys []string)
		RestartAllEvals()
	}

	// AlertRule 告警规则
	AlertRule struct {
		ctx                   *ctx.Context
		watchCtxMap           map[string]context.CancelFunc
		alarmRecoverWaitStore *storage.AlarmRecoverWaitStore
	}
)

func NewAlertRuleEval(ctx *ctx.Context, alarmRecoverWaitStore *storage.AlarmRecoverWaitStore) AlertRuleEval {
	return &AlertRule{
		ctx:                   ctx,
		watchCtxMap:           make(map[string]context.CancelFunc),
		alarmRecoverWaitStore: alarmRecoverWaitStore,
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

func (t *AlertRule) Eval(ctx context.Context, rule models.AlertRule) {
	timer := time.NewTicker(t.getEvalTimeDuration(rule.EvalTimeType, rule.EvalInterval))
	defer func() {
		timer.Stop()
		if r := recover(); r != nil {
			// 获取调用栈信息
			stack := debug.Stack()
			logc.Error(t.ctx.Ctx, fmt.Sprintf("Recovered from rule eval goroutine panic: %s, RuleName: %s, RuleId: %s\n%s", r, rule.RuleName, rule.RuleId, stack))
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
				case "AliCloudSLS", "Loki", "ElasticSearch":
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
			t.Recover(rule.RuleId, models.BuildCacheEventKey(rule.TenantId, rule.FaultCenterId), models.BuildCacheInfoKey(rule.TenantId, rule.FaultCenterId), curFingerprints)
			t.GC(rule, curFingerprints)

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

func (t *AlertRule) Recover(RuleId, faultCenterKey string, faultCenterInfoKey string, curFingerprints []string) {
	// 获取所有的故障中心的告警事件
	events, err := t.ctx.Redis.Event().GetAllEventsForFaultCenter(faultCenterKey)
	if err != nil {
		return
	}

	// 只获取当前规则的事件
	var currentRuleEvents = make(map[string]models.AlertCurEvent)
	for fingerprint, event := range events {
		if strings.Contains(event.RuleId, RuleId) {
			currentRuleEvents[fingerprint] = event
		}
	}
	events = currentRuleEvents

	// 提取事件中的告警指纹
	fingerprints := make([]string, 0)
	for fingerprint := range events {
		fingerprints = append(fingerprints, fingerprint)
	}

	// 获取已恢复告警的keys
	recoverFingerprints := tools.GetSliceDifference(fingerprints, curFingerprints)
	if recoverFingerprints == nil {
		return
	}

	curTime := time.Now().Unix()
	for _, fingerprint := range recoverFingerprints {
		event := events[fingerprint]
		if event.IsRecovered == true {
			return
		}

		// 如果是 预告警 状态的事件，触发了恢复逻辑，但它并非是真正触发告警而恢复，所以只需要删除历史事件即可，无需继续处理恢复逻辑。
		if event.Status == 0 {
			t.ctx.Redis.Event().RemoveEventFromFaultCenter(event.TenantId, event.FaultCenterId, event.Fingerprint)
			t.alarmRecoverWaitStore.Remove(RuleId, fingerprint)
			continue
		}

		// 调整为待恢复状态
		event.Status = 3
		t.ctx.Redis.Event().PushEventToFaultCenter(&event)

		// 判断是否在等待时间范围内
		wTime, exists := t.alarmRecoverWaitStore.Get(RuleId, fingerprint)
		if !exists {
			// 如果没有，则记录当前时间
			t.alarmRecoverWaitStore.Set(RuleId, fingerprint, curTime)
			continue
		}

		rt := time.Unix(wTime, 0).Add(time.Minute * time.Duration(t.getRecoverWaitTime(faultCenterInfoKey))).Unix()
		if rt > curTime {
			continue
		}

		// 已恢复状态
		event.Status = 4
		event.IsRecovered = true
		event.RecoverTime = curTime
		event.LastSendTime = 0
		t.ctx.Redis.Event().PushEventToFaultCenter(&event)
		// 触发恢复删除带恢复中的 key
		t.alarmRecoverWaitStore.Remove(RuleId, fingerprint)
	}
}

func (t *AlertRule) getRecoverWaitTime(faultCenterInfoKey string) int64 {
	faultCenter := t.ctx.Redis.FaultCenter().GetFaultCenterInfo(faultCenterInfoKey)
	if faultCenter.RecoverWaitTime == 0 {
		return 1
	}

	return faultCenter.RecoverWaitTime
}

func (t *AlertRule) GC(rule models.AlertRule, curFingerprints []string) {
	go process.GcRecoverWaitCache(t.alarmRecoverWaitStore, rule, curFingerprints)
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
