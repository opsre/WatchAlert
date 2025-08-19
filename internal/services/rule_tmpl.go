package services

import (
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
)

type ruleTmplService struct {
	ctx *ctx.Context
}

type InterRuleTmplService interface {
	List(req interface{}) (interface{}, interface{})
	Create(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Delete(req interface{}) (interface{}, interface{})
}

func newInterRuleTmplService(ctx *ctx.Context) InterRuleTmplService {
	return &ruleTmplService{
		ctx: ctx,
	}
}

func (rt ruleTmplService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRuleTemplateQuery)
	data, count, err := rt.ctx.DB.RuleTmpl().List(r.RuleGroupName, r.Type, r.Query, r.Page)
	if err != nil {
		return nil, err
	}

	return types.ResponseRuleTemplateList{
		List: data,
		Page: models.Page{
			Total: count,
			Index: r.Page.Index,
			Size:  r.Page.Size,
		},
	}, nil
}

func (rt ruleTmplService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRuleTemplateCreate)
	err := rt.ctx.DB.RuleTmpl().Create(models.RuleTemplate{
		Type:                 r.Type,
		RuleGroupName:        r.RuleGroupName,
		RuleName:             r.RuleName,
		DatasourceType:       r.DatasourceType,
		EvalInterval:         r.EvalInterval,
		ForDuration:          r.ForDuration,
		RepeatNoticeInterval: r.RepeatNoticeInterval,
		Description:          r.Description,
		PrometheusConfig:     r.PrometheusConfig,
		AliCloudSLSConfig:    r.AliCloudSLSConfig,
		LokiConfig:           r.LokiConfig,
		JaegerConfig:         r.JaegerConfig,
		KubernetesConfig:     r.KubernetesConfig,
		ElasticSearchConfig:  r.ElasticSearchConfig,
		VictoriaLogsConfig:   r.VictoriaLogsConfig,
		ClickHouseConfig:     r.ClickHouseConfig,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rt ruleTmplService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRuleTemplateUpdate)
	err := rt.ctx.DB.RuleTmpl().Update(models.RuleTemplate{
		Type:                 r.Type,
		RuleGroupName:        r.RuleGroupName,
		RuleName:             r.RuleName,
		DatasourceType:       r.DatasourceType,
		EvalInterval:         r.EvalInterval,
		ForDuration:          r.ForDuration,
		RepeatNoticeInterval: r.RepeatNoticeInterval,
		Description:          r.Description,
		PrometheusConfig:     r.PrometheusConfig,
		AliCloudSLSConfig:    r.AliCloudSLSConfig,
		LokiConfig:           r.LokiConfig,
		JaegerConfig:         r.JaegerConfig,
		KubernetesConfig:     r.KubernetesConfig,
		ElasticSearchConfig:  r.ElasticSearchConfig,
		VictoriaLogsConfig:   r.VictoriaLogsConfig,
		ClickHouseConfig:     r.ClickHouseConfig,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rt ruleTmplService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRuleTemplateQuery)
	err := rt.ctx.DB.RuleTmpl().Delete(r.RuleGroupName, r.RuleName)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
