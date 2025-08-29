package process

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"watchAlert/internal/models"
)

type ConditionEvaluator func(condition models.EvalCondition) bool

var EvalOperators = map[string]ConditionEvaluator{
	">": func(condition models.EvalCondition) bool {
		return condition.QueryValue > condition.ExpectedValue
	},
	">=": func(condition models.EvalCondition) bool {
		return condition.QueryValue >= condition.ExpectedValue
	},
	"<": func(condition models.EvalCondition) bool {
		return condition.QueryValue < condition.ExpectedValue
	},
	"<=": func(condition models.EvalCondition) bool {
		return condition.QueryValue <= condition.ExpectedValue
	},
	"==": func(condition models.EvalCondition) bool {
		return condition.QueryValue == condition.ExpectedValue
	},
	"!=": func(condition models.EvalCondition) bool {
		return condition.QueryValue != condition.ExpectedValue
	},
}

// EvalCondition 评估告警条件
func EvalCondition(ec models.EvalCondition) bool {
	evaluator, ok := EvalOperators[ec.Operator]
	if !ok {
		logc.Error(context.Background(), fmt.Sprintf("无效的评估条件, Operator: %s, ExpectedValue: %v", ec.Operator, ec.ExpectedValue))
		return false
	}

	return evaluator(ec)
}
