package probing

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"golang.org/x/sync/errgroup"
	"time"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/provider"
	"watchAlert/pkg/tools"
)

type ProductProbing struct {
	ctx         *ctx.Context
	WatchCtxMap map[string]context.CancelFunc
	Timing      map[string]int
}

func NewProbingTask(ctx *ctx.Context) ProductProbing {
	return ProductProbing{
		ctx:         ctx,
		Timing:      make(map[string]int),
		WatchCtxMap: make(map[string]context.CancelFunc),
	}
}

func (t *ProductProbing) Submit(rule models.ProbingRule) {
	t.ctx.Mux.Lock()
	defer t.ctx.Mux.Unlock()

	c, cancel := context.WithCancel(t.ctx.Ctx)
	t.WatchCtxMap[rule.RuleId] = cancel
	go t.Eval(c, rule)
}

func (t *ProductProbing) Stop(id string) {
	t.ctx.Mux.Lock()
	defer t.ctx.Mux.Unlock()

	if cancel, exists := t.WatchCtxMap[id]; exists {
		cancel()
		delete(t.WatchCtxMap, id)
	}
}

func (t *ProductProbing) Eval(ctx context.Context, rule models.ProbingRule) {
	timer := time.NewTicker(time.Second * time.Duration(rule.ProbingEndpointConfig.Strategy.EvalInterval))
	defer timer.Stop()
	t.worker(rule)

	for {
		select {
		case <-timer.C:
			logc.Infof(t.ctx.Ctx, fmt.Sprintf("网络监控: %s", tools.JsonMarshal(rule)))
			t.worker(rule)
		case <-ctx.Done():
			return
		}
	}
}

func (t *ProductProbing) worker(rule models.ProbingRule) {
	var (
		eValue     provider.EndpointValue
		err        error
		ruleConfig = rule.ProbingEndpointConfig
	)

	eValue, err = t.runEvaluation(rule)
	if err != nil {
		logc.Errorf(t.ctx.Ctx, err.Error())
		return
	}

	event := t.processDefaultEvent(rule)
	event.Fingerprint = eValue.GetFingerprint()
	event.Metric = eValue.GetLabels()
	var isValue float64
	if rule.RuleType != provider.TCPEndpointProvider {
		event.Metric["value"] = eValue[ruleConfig.Strategy.Field].(float64)
	} else {
		if eValue["IsSuccessful"] == true {
			isValue = 1
		}
		event.Metric["value"] = isValue
	}
	event.Annotations = tools.ParserVariables(rule.Annotations, event.Metric)

	var option EvalStrategy
	if rule.RuleType != provider.TCPEndpointProvider {
		option = EvalStrategy{
			Operator:      ruleConfig.Strategy.Operator,
			QueryValue:    eValue[ruleConfig.Strategy.Field].(float64),
			ExpectedValue: ruleConfig.Strategy.ExpectedValue,
		}
	} else {
		option = EvalStrategy{
			Operator:      "==",
			QueryValue:    isValue,
			ExpectedValue: 1,
		}
	}

	err = SetProbingValueMap(event.GetProbingMappingKey(), eValue)
	if err != nil {
		return
	}

	t.Evaluation(event, option)
	return
}

func (t *ProductProbing) runEvaluation(rule models.ProbingRule) (provider.EndpointValue, error) {
	var ruleConfig = rule.ProbingEndpointConfig
	switch rule.RuleType {
	case provider.ICMPEndpointProvider:
		return provider.NewEndpointPinger().Pilot(provider.EndpointOption{
			Endpoint: ruleConfig.Endpoint,
			Timeout:  ruleConfig.Strategy.Timeout,
			ICMP: provider.Eicmp{
				Interval: ruleConfig.ICMP.Interval,
				Count:    ruleConfig.ICMP.Count,
			},
		})
	case provider.HTTPEndpointProvider:
		return provider.NewEndpointHTTPer().Pilot(provider.EndpointOption{
			Endpoint: ruleConfig.Endpoint,
			Timeout:  ruleConfig.Strategy.Timeout,
			HTTP: provider.Ehttp{
				Method: ruleConfig.HTTP.Method,
				Header: ruleConfig.HTTP.Header,
				Body:   ruleConfig.HTTP.Body,
			},
		})
	case provider.TCPEndpointProvider:
		return provider.NewEndpointTcper().Pilot(provider.EndpointOption{
			Endpoint: ruleConfig.Endpoint,
			Timeout:  ruleConfig.Strategy.Timeout,
		})
	case provider.SSLEndpointProvider:
		return provider.NewEndpointSSLer().Pilot(provider.EndpointOption{
			Endpoint: ruleConfig.Endpoint,
			Timeout:  ruleConfig.Strategy.Timeout,
		})
	}
	return provider.EndpointValue{}, fmt.Errorf("unsupported rule type: %s", rule.RuleType)
}

func (t *ProductProbing) Evaluation(event models.ProbingEvent, option EvalStrategy) {
	// Retry mechanism for the condition evaluation
	if EvalCondition(option) {
		t.setTime(event.RuleId)
		if t.getTime(event.RuleId) >= event.ProbingEndpointConfig.Strategy.Failure {
			SaveProbingEndpointEvent(event)
			t.cleanTime(event.RuleId)
		}
	} else {
		c := ctx.Redis.Event()
		neCache, err := c.GetPECache(event.GetFiringAlertCacheKey())
		if err != nil {
			return
		}
		neCache.FirstTriggerTime = c.GetPEFirstTime(event.GetFiringAlertCacheKey())
		neCache.IsRecovered = true
		neCache.RecoverTime = time.Now().Unix()
		neCache.LastSendTime = 0
		c.SetPECache(neCache, 0)
		delete(t.Timing, neCache.RuleId)
	}
}

func (t *ProductProbing) processDefaultEvent(rule models.ProbingRule) models.ProbingEvent {
	return models.ProbingEvent{
		TenantId:              rule.TenantId,
		RuleId:                rule.RuleId,
		RuleType:              rule.RuleType,
		NoticeId:              rule.NoticeId,
		Severity:              rule.Severity,
		IsRecovered:           false,
		RepeatNoticeInterval:  rule.RepeatNoticeInterval,
		RecoverNotify:         rule.GetRecoverNotify(),
		ProbingEndpointConfig: rule.ProbingEndpointConfig,
	}
}

func (t *ProductProbing) RePushRule(consumer *ConsumeProbing) {
	var ruleList []models.ProbingRule
	if err := t.ctx.DB.DB().Where("enabled = ?", true).Find(&ruleList).Error; err != nil {
		logc.Errorf(t.ctx.Ctx, err.Error())
		return
	}

	g := new(errgroup.Group)
	for _, rule := range ruleList {
		rule := rule
		g.Go(func() error {
			t.Submit(rule)
			consumer.Add(rule)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		logc.Errorf(t.ctx.Ctx, err.Error())
	}
}

func (t *ProductProbing) setTime(ruleId string) {
	t.ctx.Mux.Lock()
	defer t.ctx.Mux.Unlock()

	t.Timing[ruleId]++
}

func (t *ProductProbing) getTime(ruleId string) int {
	return t.Timing[ruleId]
}

func (t *ProductProbing) cleanTime(ruleId string) {
	t.Timing[ruleId] = 0
}
