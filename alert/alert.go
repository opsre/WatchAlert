package alert

import (
	"watchAlert/alert/consumer"
	"watchAlert/alert/eval"
	"watchAlert/alert/probing"
	"watchAlert/internal/ctx"
)

var (
	AlertRule    eval.AlertRuleEval
	ConsumerWork consumer.ConsumeInterface

	ProductProbing probing.ProductProbing
	ConsumeProbing probing.ConsumeProbing
)

func Initialize(ctx *ctx.Context) {
	// 初始化告警规则评估任务
	AlertRule = eval.NewAlertRuleEval(ctx)
	AlertRule.RestartAllEvals()

	ConsumerWork = consumer.NewConsumerWork(ctx)
	ConsumerWork.RestartAllConsumers()

	// 初始化拨测任务
	ConsumeProbing = probing.NewProbingConsumerTask(ctx)
	ProductProbing = probing.NewProbingTask(ctx)
	ProductProbing.RePushRule(&ConsumeProbing)
}
