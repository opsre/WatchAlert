package services

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"watchAlert/alert"
	"watchAlert/internal/ctx"
	models "watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"
)

type ruleService struct {
	ctx *ctx.Context
}

type InterRuleService interface {
	Create(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Delete(req interface{}) (interface{}, interface{})
	List(req interface{}) (interface{}, interface{})
	Get(req interface{}) (interface{}, interface{})
	ChangeStatus(req interface{}) (interface{}, interface{})
}

func newInterRuleService(ctx *ctx.Context) InterRuleService {
	return &ruleService{
		ctx: ctx,
	}
}

func (rs ruleService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRuleCreate)
	ok := rs.ctx.DB.Rule().GetQuota(r.TenantId)
	if !ok {
		return nil, fmt.Errorf("创建失败, 配额不足")
	}

	data := models.AlertRule{
		TenantId:             r.TenantId,
		RuleId:               "a-" + tools.RandId(),
		RuleGroupId:          r.RuleGroupId,
		ExternalLabels:       r.ExternalLabels,
		DatasourceType:       r.DatasourceType,
		DatasourceIdList:     r.DatasourceIdList,
		RuleName:             r.RuleName,
		EvalInterval:         r.EvalInterval,
		EvalTimeType:         r.EvalTimeType,
		RepeatNoticeInterval: r.RepeatNoticeInterval,
		Description:          r.Description,
		EffectiveTime:        r.EffectiveTime,
		Severity:             r.Severity,
		PrometheusConfig:     r.PrometheusConfig,
		AliCloudSLSConfig:    r.AliCloudSLSConfig,
		LokiConfig:           r.LokiConfig,
		VictoriaLogsConfig:   r.VictoriaLogsConfig,
		ClickHouseConfig:     r.ClickHouseConfig,
		JaegerConfig:         r.JaegerConfig,
		CloudWatchConfig:     r.CloudWatchConfig,
		KubernetesConfig:     r.KubernetesConfig,
		ElasticSearchConfig:  r.ElasticSearchConfig,
		LogEvalCondition:     r.LogEvalCondition,
		FaultCenterId:        r.FaultCenterId,
		Enabled:              r.Enabled,
	}

	if *r.GetEnabled() {
		alert.AlertRule.Submit(data)
	}

	err := rs.ctx.DB.Rule().Create(data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rs ruleService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRuleUpdate)
	oldRule := models.AlertRule{}
	rs.ctx.DB.DB().Model(&models.AlertRule{}).
		Where("tenant_id = ? AND rule_id = ?", r.TenantId, r.RuleId).
		First(&oldRule)

	if oldRule.FaultCenterId != r.FaultCenterId {
		fingerprints := rs.ctx.Redis.Alert().GetFingerprintsByRuleId(oldRule.TenantId, oldRule.FaultCenterId, oldRule.RuleId)
		for _, fingerprint := range fingerprints {
			rs.ctx.Redis.Alert().RemoveAlertEvent(oldRule.TenantId, oldRule.FaultCenterId, fingerprint)
		}
	}

	/*
		重启协程
		判断当前状态是否是false 并且 历史状态是否为true
	*/
	if *oldRule.Enabled == true && *r.Enabled == false {
		alert.AlertRule.Stop(r.RuleId)
	}
	if *oldRule.Enabled == true && *r.Enabled == true {
		alert.AlertRule.Stop(r.RuleId)
	}

	data := models.AlertRule{
		TenantId:             r.TenantId,
		RuleId:               r.RuleId,
		RuleGroupId:          r.RuleGroupId,
		ExternalLabels:       r.ExternalLabels,
		DatasourceType:       r.DatasourceType,
		DatasourceIdList:     r.DatasourceIdList,
		RuleName:             r.RuleName,
		EvalInterval:         r.EvalInterval,
		EvalTimeType:         r.EvalTimeType,
		RepeatNoticeInterval: r.RepeatNoticeInterval,
		Description:          r.Description,
		EffectiveTime:        r.EffectiveTime,
		Severity:             r.Severity,
		PrometheusConfig:     r.PrometheusConfig,
		AliCloudSLSConfig:    r.AliCloudSLSConfig,
		LokiConfig:           r.LokiConfig,
		VictoriaLogsConfig:   r.VictoriaLogsConfig,
		ClickHouseConfig:     r.ClickHouseConfig,
		JaegerConfig:         r.JaegerConfig,
		CloudWatchConfig:     r.CloudWatchConfig,
		KubernetesConfig:     r.KubernetesConfig,
		ElasticSearchConfig:  r.ElasticSearchConfig,
		LogEvalCondition:     r.LogEvalCondition,
		FaultCenterId:        r.FaultCenterId,
		Enabled:              r.Enabled,
	}

	// 启动协程
	if *r.GetEnabled() {
		alert.AlertRule.Submit(data)
		logc.Infof(rs.ctx.Ctx, fmt.Sprintf("重启 RuleId 为 %s 的 Worker 进程", r.RuleId))
	} else {
		// 删除缓存
		fingerprints := rs.ctx.Redis.Alert().GetFingerprintsByRuleId(r.TenantId, r.FaultCenterId, r.RuleId)
		for _, fingerprint := range fingerprints {
			rs.ctx.Redis.Alert().RemoveAlertEvent(r.TenantId, r.FaultCenterId, fingerprint)
		}
	}

	// 更新数据
	err := rs.ctx.DB.Rule().Update(data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rs ruleService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRuleQuery)
	info, err := rs.ctx.DB.Rule().Get(r.TenantId, r.RuleGroupId, r.RuleId)
	if err != nil {
		return nil, err
	}

	err = rs.ctx.DB.Rule().Delete(r.TenantId, r.RuleId)
	if err != nil {
		return nil, err
	}

	// 退出该规则的协程
	if *info.GetEnabled() {
		logc.Infof(rs.ctx.Ctx, fmt.Sprintf("停止 RuleId 为 %s 的 Worker 进程", r.RuleId))
		alert.AlertRule.Stop(r.RuleId)
	}

	// 删除缓存
	fingerprints := rs.ctx.Redis.Alert().GetFingerprintsByRuleId(r.TenantId, info.FaultCenterId, r.RuleId)
	for _, fingerprint := range fingerprints {
		rs.ctx.Redis.Alert().RemoveAlertEvent(r.TenantId, info.FaultCenterId, fingerprint)
	}

	return nil, nil
}

func (rs ruleService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRuleQuery)
	data, err := rs.ctx.DB.Rule().List(r.TenantId, r.RuleGroupId, r.DatasourceType, r.Query, r.Status, r.Page)
	if err != nil {
		return nil, err
	}

	return types.ResponseRuleList{
		List: data,
		Page: models.Page{
			Total: int64(len(data)),
			Index: r.Page.Index,
			Size:  r.Page.Size,
		},
	}, nil
}

func (rs ruleService) Get(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRuleQuery)
	data, err := rs.ctx.DB.Rule().Get(r.TenantId, r.RuleGroupId, r.RuleId)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (rs ruleService) ChangeStatus(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRuleChangeStatus)
	switch *r.GetEnabled() {
	case true:
		logc.Infof(rs.ctx.Ctx, fmt.Sprintf("重启 RuleId 为 %s 的 Worker 进程", r.RuleId))
		var enable = true
		rule := rs.ctx.DB.Rule().GetRuleObject(r.RuleId)
		newRule := rule
		newRule.Enabled = &enable
		alert.AlertRule.Submit(newRule)

	case false:
		logc.Infof(rs.ctx.Ctx, fmt.Sprintf("停止 RuleId 为 %s 的 Worker 进程", r.RuleId))
		alert.AlertRule.Stop(r.RuleId)
		// 删除缓存
		fingerprints := rs.ctx.Redis.Alert().GetFingerprintsByRuleId(r.TenantId, r.FaultCenterId, r.RuleId)
		for _, fingerprint := range fingerprints {
			rs.ctx.Redis.Alert().RemoveAlertEvent(r.TenantId, r.FaultCenterId, fingerprint)
		}
	}

	return nil, rs.ctx.DB.Rule().ChangeStatus(r.TenantId, r.RuleGroupId, r.RuleId, r.GetEnabled())
}
