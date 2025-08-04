package services

import (
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"gopkg.in/yaml.v3"
	"time"
	"watchAlert/alert"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
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
	Import(req interface{}) (interface{}, interface{})
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
	data, count, err := rs.ctx.DB.Rule().List(r.TenantId, r.RuleGroupId, r.DatasourceType, r.Query, r.Status, r.Page)
	if err != nil {
		return nil, err
	}

	return types.ResponseRuleList{
		List: data,
		Page: models.Page{
			Total: count,
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

func (rs ruleService) Import(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRuleImport)
	var rules []types.RequestRuleCreate

	switch r.ImportType {
	case types.WithPrometheusRuleImport:
		var alerts types.PrometheusAlerts
		err := yaml.Unmarshal([]byte(r.Rules), &alerts)
		if err != nil {
			return nil, err
		}

		for _, alert := range alerts.Rules {
			var forDuration int64
			d, err := time.ParseDuration(alert.For)
			if err != nil {
				forDuration = 1
				logc.Error(rs.ctx.Ctx, err.Error())
			}
			forDuration = int64(d.Seconds())

			rule := types.RequestRuleCreate{
				TenantId:         r.TenantId,
				RuleGroupId:      r.RuleGroupId,
				ExternalLabels:   alert.Labels,
				DatasourceType:   r.DatasourceType,
				DatasourceIdList: r.DatasourceIdList,
				RuleName:         alert.Alert,
				EvalInterval:     15,
				EvalTimeType:     "second",
				PrometheusConfig: models.PrometheusConfig{
					PromQL:      alert.Expr,
					Annotations: alert.Annotations.Description,
					Rules: []models.Rules{
						{
							ForDuration: forDuration,
							Severity:    "P1",
							Expr:        ">= 0",
						},
					},
				},
				FaultCenterId: r.FaultCenterId,
				Enabled:       alert.GetEnable(),
			}
			rules = append(rules, rule)
		}

	case types.WithWatchAlertJsonImport:
		err := json.Unmarshal([]byte(r.Rules), &rules)
		if err != nil {
			return nil, err
		}
	}

	if len(rules) == 0 {
		return nil, fmt.Errorf("导入失败, 识别到 0 条规则")
	}

	for _, rule := range rules {
		if len(rule.RuleGroupId) == 0 {
			continue
		}

		err := rs.ctx.DB.Rule().Create(models.AlertRule{
			TenantId:             r.TenantId,
			RuleId:               "a-" + tools.RandId(),
			RuleGroupId:          rule.RuleGroupId,
			ExternalLabels:       rule.ExternalLabels,
			DatasourceType:       rule.DatasourceType,
			DatasourceIdList:     rule.DatasourceIdList,
			RuleName:             rule.RuleName,
			EvalInterval:         rule.EvalInterval,
			EvalTimeType:         rule.EvalTimeType,
			RepeatNoticeInterval: rule.RepeatNoticeInterval,
			Description:          rule.Description,
			EffectiveTime:        rule.EffectiveTime,
			Severity:             rule.Severity,
			PrometheusConfig:     rule.PrometheusConfig,
			AliCloudSLSConfig:    rule.AliCloudSLSConfig,
			LokiConfig:           rule.LokiConfig,
			VictoriaLogsConfig:   rule.VictoriaLogsConfig,
			ClickHouseConfig:     rule.ClickHouseConfig,
			JaegerConfig:         rule.JaegerConfig,
			CloudWatchConfig:     rule.CloudWatchConfig,
			KubernetesConfig:     rule.KubernetesConfig,
			ElasticSearchConfig:  rule.ElasticSearchConfig,
			LogEvalCondition:     rule.LogEvalCondition,
			FaultCenterId:        rule.FaultCenterId,
			Enabled:              rule.Enabled,
		})
		if err != nil {
			logc.Errorf(rs.ctx.Ctx, err.Error())
			continue
		}
	}

	return nil, nil
}
