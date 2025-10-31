package alert

import (
	"context"
	"watchAlert/alert/consumer"
	"watchAlert/alert/eval"
	"watchAlert/alert/probing"
	"watchAlert/internal/ctx"
	"watchAlert/internal/global"
	"watchAlert/pkg/client"
	"watchAlert/pkg/tools"

	"github.com/zeromicro/go-zero/core/logc"
)

var (
	AlertRule    eval.AlertRuleEval
	ConsumerWork consumer.ConsumeInterface

	ProductProbing probing.ProductProbing
	ConsumeProbing probing.ConsumeProbing

	// Leader 选举器
	LeaderElector *tools.LeaderElector

	// 消息订阅取消函数
	subscriberCancels []context.CancelFunc

	// 选举开关
	leaderElectionEnabled bool
)

func Initialize(ctx *ctx.Context) {
	// 初始化告警规则评估任务
	AlertRule = eval.NewAlertRuleEval(ctx)
	ConsumerWork = consumer.NewConsumerWork(ctx)

	// 初始化拨测任务
	ConsumeProbing = probing.NewProbingConsumerTask(ctx)
	ProductProbing = probing.NewProbingTask(ctx)

	// 检查 Leader 选举是否启用
	leaderElectionEnabled = global.Config.Server.EnableElection

	if leaderElectionEnabled {
		// 启用 Leader 选举模式
		logc.Infof(ctx.Ctx, "Leader 选举已启用，开始选举流程...")
		LeaderElector = tools.NewLeaderElector(
			ctx.Ctx,
			client.Redis,
			loadRules,
			unloadRules,
		)
		// 启动 Leader 选举
		LeaderElector.Start()
	} else {
		loadRules()
	}
}

// loadRules 加载所有规则(成为 Leader 时调用)
func loadRules() {
	logc.Infof(ctx.Ctx, "本节点为 Leader 节点，开始加载规则...")

	// 重启所有告警规则评估器
	AlertRule.RestartAllEvals()

	// 重启所有故障中心消费者
	ConsumerWork.RestartAllConsumers()

	// 重启所有拨测任务
	ProductProbing.RePushRule(&ConsumeProbing)

	// 启动 Redis 消息订阅，监听规则变更
	startMessageSubscribers()
}

// startMessageSubscribers 启动消息订阅器
func startMessageSubscribers() {
	subscriberCancels = make([]context.CancelFunc, 0)

	// 订阅告警规则重载消息
	subCtx1, cancel1 := context.WithCancel(ctx.Ctx)
	subscriberCancels = append(subscriberCancels, cancel1)
	go tools.SubscribeReloadMessages(subCtx1, client.Redis, tools.ChannelRuleReload, handleRuleReload)

	// 订阅故障中心重载消息
	subCtx2, cancel2 := context.WithCancel(ctx.Ctx)
	subscriberCancels = append(subscriberCancels, cancel2)
	go tools.SubscribeReloadMessages(subCtx2, client.Redis, tools.ChannelFaultCenterReload, handleFaultCenterReload)

	// 订阅拨测规则重载消息
	subCtx3, cancel3 := context.WithCancel(ctx.Ctx)
	subscriberCancels = append(subscriberCancels, cancel3)
	go tools.SubscribeReloadMessages(subCtx3, client.Redis, tools.ChannelProbingReload, handleProbingReload)
}

// stopMessageSubscribers 停止消息订阅器
func stopMessageSubscribers() {
	for _, cancel := range subscriberCancels {
		cancel()
	}
	subscriberCancels = nil
	logc.Infof(ctx.Ctx, "消息订阅器已停止")
}

// handleRuleReload 处理告警规则重载消息
func handleRuleReload(msg tools.ReloadMessage) {

	// 从数据库获取规则
	rule := ctx.DB.Rule().GetRuleObject(msg.ID)
	if rule.RuleId == "" {
		logc.Errorf(ctx.Ctx, "规则不存在: %s", msg.ID)
		return
	}

	switch msg.Action {
	case tools.ActionCreate, tools.ActionEnable:
		if rule.Enabled != nil && *rule.Enabled {
			AlertRule.Submit(rule)
			logc.Infof(ctx.Ctx, "[Leader] 已启动规则评估: %s", msg.Name)
		}

	case tools.ActionUpdate:
		AlertRule.Stop(msg.ID)
		if rule.Enabled != nil && *rule.Enabled {
			AlertRule.Submit(rule)
			logc.Infof(ctx.Ctx, "[Leader] 已重启规则评估: %s", msg.Name)
		}

	case tools.ActionDelete, tools.ActionDisable:
		AlertRule.Stop(msg.ID)
		logc.Infof(ctx.Ctx, "[Leader] 已停止规则评估: %s", msg.Name)
	}
}

// handleFaultCenterReload 处理故障中心重载消息
func handleFaultCenterReload(msg tools.ReloadMessage) {
	fc, err := ctx.DB.FaultCenter().Get(msg.TenantID, msg.ID, "")
	if err != nil {
		logc.Errorf(ctx.Ctx, "故障中心不存在: %s, err: %v", msg.ID, err)
		return
	}

	switch msg.Action {
	case tools.ActionCreate, tools.ActionEnable:
		ConsumerWork.Submit(fc)
		logc.Infof(ctx.Ctx, "[Leader] 已启动故障中心消费: %s", msg.Name)

	case tools.ActionUpdate:
		ConsumerWork.Stop(msg.ID)
		ConsumerWork.Submit(fc)
		logc.Infof(ctx.Ctx, "[Leader] 已重启故障中心消费: %s", msg.Name)

	case tools.ActionDelete, tools.ActionDisable:
		ConsumerWork.Stop(msg.ID)
		logc.Infof(ctx.Ctx, "[Leader] 已停止故障中心消费: %s", msg.Name)
	}
}

// handleProbingReload 处理拨测规则重载消息
func handleProbingReload(msg tools.ReloadMessage) {
	rule, err := ctx.DB.Probing().Search(msg.TenantID, msg.ID)
	if err != nil {
		logc.Errorf(ctx.Ctx, "拨测规则不存在: %s, err: %v", msg.ID, err)
		return
	}
	switch msg.Action {
	case tools.ActionCreate, tools.ActionEnable:
		if rule.Enabled != nil && *rule.Enabled {
			ProductProbing.Add(rule)
			ConsumeProbing.Add(rule)
			logc.Infof(ctx.Ctx, "[Leader] 已启动拨测任务: %s", msg.Name)
		}

	case tools.ActionUpdate:
		ProductProbing.Stop(msg.ID)
		ConsumeProbing.Stop(msg.ID)
		if rule.Enabled != nil && *rule.Enabled {
			ProductProbing.Add(rule)
			ConsumeProbing.Add(rule)
			logc.Infof(ctx.Ctx, "[Leader] 已重启拨测任务: %s", msg.Name)
		}

	case tools.ActionDelete, tools.ActionDisable:
		ProductProbing.Stop(msg.ID)
		ConsumeProbing.Stop(msg.ID)
		logc.Infof(ctx.Ctx, "[Leader] 已停止拨测任务: %s", msg.Name)
	}
}

// unloadRules 卸载所有规则(失去 Leader 时调用)
func unloadRules() {
	logc.Infof(ctx.Ctx, "本节点失去 Leader 身份，停止所有任务...")

	// 停止消息订阅
	stopMessageSubscribers()

	// 停止所有告警规则评估器
	AlertRule.StopAllEvals()

	// 停止所有故障中心消费者
	ConsumerWork.StopAllConsumers()

	// 停止所有拨测任务
	ProductProbing.StopAllTasks()
	ConsumeProbing.StopAllTasks()
}

// IsLeader 判断节点角色
func IsLeader() bool {
	if !leaderElectionEnabled {
		return true
	}

	return LeaderElector != nil && LeaderElector.IsLeader()
}
