package services

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"watchAlert/alert"
	"watchAlert/internal/ctx"
	models "watchAlert/internal/models"
)

type ruleService struct {
	ctx *ctx.Context
}

type InterRuleService interface {
	Create(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Delete(req interface{}) (interface{}, interface{})
	List(req interface{}) (interface{}, interface{})
	Search(req interface{}) (interface{}, interface{})
	ChangeStatus(req interface{}) (interface{}, interface{})
}

func newInterRuleService(ctx *ctx.Context) InterRuleService {
	return &ruleService{
		ctx: ctx,
	}
}

func (rs ruleService) Create(req interface{}) (interface{}, interface{}) {
	rule := req.(*models.AlertRule)
	ok := rs.ctx.DB.Rule().GetQuota(rule.TenantId)
	if !ok {
		return nil, fmt.Errorf("创建失败, 配额不足")
	}

	alert.AlertRule.Submit(*rule)

	err := rs.ctx.DB.Rule().Create(*rule)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rs ruleService) Update(req interface{}) (interface{}, interface{}) {
	rule := req.(*models.AlertRule)
	oldRule := models.AlertRule{}
	rs.ctx.DB.DB().Model(&models.AlertRule{}).
		Where("tenant_id = ? AND rule_id = ?", rule.TenantId, rule.RuleId).
		First(&oldRule)

	if oldRule.FaultCenterId != rule.FaultCenterId {
		fingerprints := rs.ctx.Redis.Alert().GetFingerprintsByRuleId(oldRule.TenantId, oldRule.FaultCenterId, oldRule.RuleId)
		for _, fingerprint := range fingerprints {
			rs.ctx.Redis.Alert().RemoveAlertEvent(oldRule.TenantId, oldRule.FaultCenterId, fingerprint)
		}
	}

	/*
		重启协程
		判断当前状态是否是false 并且 历史状态是否为true
	*/
	if *oldRule.Enabled == true && *rule.Enabled == false {
		alert.AlertRule.Stop(rule.RuleId)
	}
	if *oldRule.Enabled == true && *rule.Enabled == true {
		alert.AlertRule.Stop(rule.RuleId)
	}

	// 启动协程
	if *rule.GetEnabled() {
		alert.AlertRule.Submit(*rule)
		logc.Infof(rs.ctx.Ctx, fmt.Sprintf("重启 RuleId 为 %s 的 Worker 进程", rule.RuleId))
	} else {
		// 删除缓存
		fingerprints := rs.ctx.Redis.Alert().GetFingerprintsByRuleId(rule.TenantId, rule.FaultCenterId, rule.RuleId)
		for _, fingerprint := range fingerprints {
			rs.ctx.Redis.Alert().RemoveAlertEvent(rule.TenantId, rule.FaultCenterId, fingerprint)
		}
	}

	// 更新数据
	err := rs.ctx.DB.Rule().Update(*rule)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rs ruleService) Delete(req interface{}) (interface{}, interface{}) {
	rule := req.(*models.AlertRuleQuery)

	info, err := rs.ctx.DB.Rule().Search(*rule)
	if err != nil {
		return nil, err
	}

	err = rs.ctx.DB.Rule().Delete(*rule)
	if err != nil {
		return nil, err
	}

	// 退出该规则的协程
	if *info.GetEnabled() {
		logc.Infof(rs.ctx.Ctx, fmt.Sprintf("停止 RuleId 为 %s 的 Worker 进程", rule.RuleId))
		alert.AlertRule.Stop(rule.RuleId)
	}

	// 删除缓存
	fingerprints := rs.ctx.Redis.Alert().GetFingerprintsByRuleId(rule.TenantId, info.FaultCenterId, rule.RuleId)
	for _, fingerprint := range fingerprints {
		rs.ctx.Redis.Alert().RemoveAlertEvent(rule.TenantId, info.FaultCenterId, fingerprint)
	}

	return nil, nil
}

func (rs ruleService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*models.AlertRuleQuery)
	data, err := rs.ctx.DB.Rule().List(*r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (rs ruleService) Search(req interface{}) (interface{}, interface{}) {
	r := req.(*models.AlertRuleQuery)
	data, err := rs.ctx.DB.Rule().Search(*r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (rs ruleService) ChangeStatus(req interface{}) (interface{}, interface{}) {
	r := req.(*models.RequestRuleChangeStatus)
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

	return nil, rs.ctx.DB.Rule().ChangeStatus(*r)
}
