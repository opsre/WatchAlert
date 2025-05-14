package probing

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"golang.org/x/sync/errgroup"
	"time"
	"watchAlert/alert/process"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/provider"
	"watchAlert/pkg/tools"
)

type ProductProbing struct {
	ctx           *ctx.Context
	WatchCtxMap   map[string]context.CancelFunc
	FailFrequency map[string]int
	OkFrequency   map[string]int
}

func NewProbingTask(ctx *ctx.Context) ProductProbing {
	return ProductProbing{
		ctx:           ctx,
		FailFrequency: make(map[string]int),
		OkFrequency:   make(map[string]int),
		WatchCtxMap:   make(map[string]context.CancelFunc),
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

	eValue, err = t.runProbing(rule)
	if err != nil {
		logc.Errorf(t.ctx.Ctx, err.Error())
		return
	}
	err = t.ctx.DB.Probing().AddRecord(models.ProbingHistory{
		Timestamp: time.Now().Unix(),
		RuleId:    rule.RuleId,
		Value:     eValue,
	})
	if err != nil {
		logc.Errorf(t.ctx.Ctx, err.Error())
		return
	}

	event := t.buildEvent(rule)
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

	var option models.EvalCondition
	switch rule.RuleType {
	// 如果拨测类型是 TCP ，直接定义好计算条件 == 0 则表示异常
	case provider.TCPEndpointProvider:
		option = models.EvalCondition{
			Operator:      "==",
			QueryValue:    isValue,
			ExpectedValue: 0,
		}
	default:
		option = models.EvalCondition{
			Operator:      ruleConfig.Strategy.Operator,
			QueryValue:    eValue[ruleConfig.Strategy.Field].(float64),
			ExpectedValue: ruleConfig.Strategy.ExpectedValue,
		}
	}

	err = SetProbingValueMap(models.BuildProbingValueCacheKey(event.TenantId, event.RuleId), eValue)
	if err != nil {
		return
	}

	t.Evaluation(event, option)
	return
}

func (t *ProductProbing) runProbing(rule models.ProbingRule) (provider.EndpointValue, error) {
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

func (t *ProductProbing) Evaluation(event models.ProbingEvent, option models.EvalCondition) {
	if process.EvalCondition(option) {
		// 控制失败频次
		t.setFrequency(t.FailFrequency, event.RuleId)
		// 如果失败频次达到设定次数后记录事件
		if t.getFrequency(t.FailFrequency, event.RuleId) >= event.ProbingEndpointConfig.Strategy.Failure {
			defer func() {
				t.cleanFrequency(t.FailFrequency, event.RuleId)
			}()

			SaveProbingEndpointEvent(t.ctx, event)
		}
	} else {
		// 控制成功频次
		t.setFrequency(t.OkFrequency, event.RuleId)
		if t.getFrequency(t.OkFrequency, event.RuleId) >= 3 {
			defer func() {
				t.cleanFrequency(t.OkFrequency, event.RuleId)
			}()

			key := models.BuildProbingEventCacheKey(event.TenantId, event.RuleId)
			c := ctx.Redis.Probing()
			neCache, err := c.GetProbingEventCache(key)
			if err != nil {
				logc.Error(ctx.Ctx, err.Error())
				return
			}
			neCache.FirstTriggerTime = c.GetProbingEventFirstTime(key)
			neCache.IsRecovered = true
			neCache.RecoverTime = time.Now().Unix()
			neCache.LastSendTime = 0
			c.SetProbingEventCache(neCache, 0)
		}
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

func (t *ProductProbing) setFrequency(frequencyStorage map[string]int, ruleId string) {
	t.ctx.Mux.Lock()
	defer t.ctx.Mux.Unlock()

	frequencyStorage[ruleId]++
}

func (t *ProductProbing) getFrequency(frequencyStorage map[string]int, ruleId string) int {
	t.ctx.Mux.RLock()
	defer t.ctx.Mux.RUnlock()

	return frequencyStorage[ruleId]
}

func (t *ProductProbing) cleanFrequency(frequencyStorage map[string]int, ruleId string) {
	delete(frequencyStorage, ruleId)
}
