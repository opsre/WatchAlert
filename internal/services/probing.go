package services

import (
	"watchAlert/alert/probing"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/provider"
	"watchAlert/pkg/tools"
)

type (
	probingService struct {
		ctx          *ctx.Context
		ProductTask  *probing.ProductProbing
		ConsumerTask *probing.ConsumeProbing
	}

	InterProbingService interface {
		Create(req interface{}) (interface{}, interface{})
		Update(req interface{}) (interface{}, interface{})
		Delete(req interface{}) (interface{}, interface{})
		List(req interface{}) (interface{}, interface{})
		Search(req interface{}) (interface{}, interface{})
		Once(req interface{}) (interface{}, interface{})
	}
)

func newInterProbingService(ctx *ctx.Context, NetworkMonProduct *probing.ProductProbing, NetworkMonConsumer *probing.ConsumeProbing) InterProbingService {
	return &probingService{
		ctx:          ctx,
		ProductTask:  NetworkMonProduct,
		ConsumerTask: NetworkMonConsumer,
	}
}

func (m probingService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*models.ProbingRule)
	r.RuleId = "r-" + tools.RandId()

	err := m.ctx.DB.Probing().Create(*r)
	if err != nil {
		return nil, err
	}

	if *r.GetEnabled() {
		m.ProductTask.Submit(*r)
	}
	m.ConsumerTask.Add(*r)

	return nil, nil
}

func (m probingService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*models.ProbingRule)
	_, err := m.ctx.DB.Probing().Search(models.ProbingRuleQuery{RuleId: r.RuleId})
	if err != nil {
		return nil, err
	}

	err = m.ctx.DB.Probing().Update(*r)
	if err != nil {
		return nil, err
	}

	m.ProductTask.Stop(r.RuleId)
	m.ConsumerTask.Stop(r.RuleId)
	if *r.GetEnabled() {
		m.ProductTask.Submit(*r)
		m.ConsumerTask.Add(*r)
	}

	return nil, nil
}

func (m probingService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*models.ProbingRuleQuery)
	res, err := m.ctx.DB.Probing().Search(models.ProbingRuleQuery{RuleId: r.RuleId})
	if err != nil {
		return nil, err
	}

	err = m.ctx.DB.Probing().Delete(*r)
	if err != nil {
		return nil, err
	}

	m.ProductTask.Stop(r.RuleId)
	m.ConsumerTask.Stop(r.RuleId)
	err = m.ctx.Redis.Redis().Del(res.GetFiringAlertCacheKey(), res.GetProbingMappingKey()).Err()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (m probingService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*models.ProbingRuleQuery)
	data, err := m.ctx.DB.Probing().List(*r)
	if err != nil {
		return nil, err
	}

	for k, v := range data {
		value := &data[k].ProbingEndpointValues
		nv := probing.GetProbingValueMap(v.GetProbingMappingKey())
		switch r.RuleType {
		case provider.HTTPEndpointProvider:
			value.PHTTP.Latency = nv["Latency"]
			value.PHTTP.StatusCode = nv["StatusCode"]
		case provider.ICMPEndpointProvider:
			value.PICMP.PacketLoss = nv["PacketLoss"]
			value.PICMP.MinRtt = nv["MinRtt"]
			value.PICMP.MaxRtt = nv["MaxRtt"]
			value.PICMP.AvgRtt = nv["AvgRtt"]
		case provider.TCPEndpointProvider:
			value.PTCP.ErrorMessage = nv["ErrorMessage"]
			value.PTCP.IsSuccessful = nv["IsSuccessful"]
		case provider.SSLEndpointProvider:
			value.PSSL.ExpireTime = nv["ExpireTime"]
			value.PSSL.StartTime = nv["StartTime"]
			value.PSSL.ResponseTime = nv["ResponseTime"]
			value.PSSL.TimeRemaining = nv["TimeRemaining"]
		}
	}

	return data, nil
}

func (m probingService) Search(req interface{}) (interface{}, interface{}) {
	r := req.(*models.ProbingRuleQuery)
	data, err := m.ctx.DB.Probing().Search(*r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (m probingService) Once(req interface{}) (interface{}, interface{}) {
	r := req.(*models.OnceProbing)
	var ruleConfig = r.ProbingEndpointConfig
	switch r.RuleType {
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
	return nil, nil
}
