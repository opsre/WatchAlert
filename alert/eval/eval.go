package eval

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"runtime/debug"
	"time"
	"watchAlert/alert/process"
	"watchAlert/alert/storage"
	"watchAlert/internal/global"
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
		Recover(rule models.AlertRule, curKeys []string)
		GC(rule models.AlertRule, curFiringKeys, curPendingKeys []string)
		RePushTask()
	}

	// AlertRule 告警规则
	AlertRule struct {
		ctx                   *ctx.Context
		watchCtxMap           map[string]context.CancelFunc
		alarmRecoverWaitStore storage.AlarmRecoverWaitStore
	}
)

func NewAlertRuleEval(ctx *ctx.Context, alarmRecoverWaitStore storage.AlarmRecoverWaitStore) AlertRuleEval {
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
	timer := time.NewTicker(time.Second * time.Duration(rule.EvalInterval))
	defer func() {
		timer.Stop()
		if r := recover(); r != nil {
			logc.Errorf(t.ctx.Ctx, "Recovered from panic in rule eval: %v\nStack: %s", r, debug.Stack())
		}
	}()

	for {
		select {
		case <-timer.C:
			// 检查规则是否仍然启用
			if !t.isRuleEnabled(rule.RuleId) {
				logc.Infof(t.ctx.Ctx, "RuleId: %s 已禁用，退出评估循环", rule.RuleId)
				return
			}

			for _, dsId := range rule.DatasourceIdList {
				curFiringKeys := []string{}
				curPendingKeys := []string{}

				instance, err := t.ctx.DB.Datasource().GetInstance(dsId)
				if err != nil {
					logc.Warnf(t.ctx.Ctx, "获取数据源失败: %s (dsId=%s)", err.Error(), dsId)
					continue
				}

				// 检查数据源健康状态
				if !provider.CheckDatasourceHealth(instance) {
					logc.Warnf(t.ctx.Ctx, "数据源 %s 不健康，跳过处理", dsId)
					continue
				}

				// 处理不同类型的数据源
				switch rule.DatasourceType {
				case "Prometheus", "VictoriaMetrics":
					curFiringKeys, curPendingKeys = metrics(t.ctx, dsId, instance.Type, rule)
				case "AliCloudSLS", "Loki", "ElasticSearch":
					curFiringKeys = logs(t.ctx, dsId, instance.Type, rule)
				case "Jaeger":
					curFiringKeys = traces(t.ctx, dsId, instance.Type, rule)
				case "CloudWatch":
					curFiringKeys = cloudWatch(t.ctx, dsId, rule)
				case "KubernetesEvent":
					curFiringKeys = kubernetesEvent(t.ctx, dsId, rule)
				}

				// 记录规则评估日志
				if len(curFiringKeys) == 0 && len(curPendingKeys) == 0 {
					logc.Debug(t.ctx.Ctx, "数据源 %s 规则 %s (RuleId: %s) 没有触发状态", dsId, rule.RuleName, rule.RuleId)
				} else {
					logc.Infof(t.ctx.Ctx, "数据源 %s 规则评估 -> %v", dsId, tools.JsonMarshal(rule))
				}

				// 立即执行恢复和 GC 逻辑
				t.Recover(rule, curFiringKeys)
				t.GC(rule, curFiringKeys, curPendingKeys)
			}

			// 重新创建定时器，保证新的 `EvalInterval` 生效
			timer.Stop()
			timer = time.NewTicker(time.Second * time.Duration(rule.EvalInterval))

		case <-ctx.Done():
			logc.Infof(t.ctx.Ctx, "停止 RuleId: %s, RuleName: %s 的 Watch 协程", rule.RuleId, rule.RuleName)
			return
		}
	}
}

func (t *AlertRule) Recover(rule models.AlertRule, curKeys []string) {
	firingKeys, err := ctx.Redis.Rule().GetAlertFiringCacheKeys(models.AlertRuleQuery{
		TenantId:         rule.TenantId,
		RuleId:           rule.RuleId,
		DatasourceIdList: rule.DatasourceIdList,
	})
	if err != nil {
		return
	}
	// 获取已恢复告警的keys
	recoverKeys := tools.GetSliceDifference(firingKeys, curKeys)
	if recoverKeys == nil {
		return
	}

	curTime := time.Now().Unix()
	for _, key := range recoverKeys {
		event := ctx.Redis.Event().GetCache(key)
		if event.IsRecovered == true {
			return
		}

		if _, exists := t.alarmRecoverWaitStore.Get(key); !exists {
			// 如果没有，则记录当前时间
			t.alarmRecoverWaitStore.Set(key, curTime)
			continue
		}

		// 判断是否在等待时间范围内
		wTime, _ := t.alarmRecoverWaitStore.Get(key)
		rt := time.Unix(wTime, 0).Add(time.Minute * time.Duration(global.Config.Server.AlarmConfig.RecoverWait)).Unix()
		if rt > curTime {
			continue
		}

		event.IsRecovered = true
		event.RecoverTime = curTime
		event.LastSendTime = 0

		ctx.Redis.Event().SetCache("Firing", event, 0)

		// 触发恢复删除带恢复中的 key
		t.alarmRecoverWaitStore.Remove(key)
	}
}

func (t *AlertRule) GC(rule models.AlertRule, curFiringKeys, curPendingKeys []string) {
	go process.GcPendingCache(t.ctx, rule, curPendingKeys)
	go process.GcRecoverWaitCache(t.ctx, t.alarmRecoverWaitStore, rule, curFiringKeys)
}

func (t *AlertRule) RePushTask() {
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
