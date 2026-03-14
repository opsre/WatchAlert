package services

import (
	"fmt"
	"time"
	"watchAlert/alert"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/client"
	"watchAlert/pkg/tools"
)

type recordingRuleService struct {
	ctx *ctx.Context
}

type InterRecordingRuleService interface {
	Create(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Delete(req interface{}) (interface{}, interface{})
	List(req interface{}) (interface{}, interface{})
	Get(req interface{}) (interface{}, interface{})
	ChangeStatus(req interface{}) (interface{}, interface{})
}

func newInterRecordingRuleService(ctx *ctx.Context) InterRecordingRuleService {
	return &recordingRuleService{
		ctx: ctx,
	}
}

func (rs recordingRuleService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRecordingRuleCreate)

	// 如果指定了ruleGroupId，验证规则组是否存在
	if r.RuleGroupId > 0 {
		var group models.RecordingRuleGroup
		err := rs.ctx.DB.DB().Model(&models.RecordingRuleGroup{}).
			Where("tenant_id = ? AND id = ?", r.TenantId, r.RuleGroupId).
			First(&group).Error
		if err != nil {
			return nil, fmt.Errorf("规则组不存在")
		}
	}

	data := models.RecordingRule{
		TenantId:       r.TenantId,
		RuleId:         "rr-" + tools.RandId(),
		DatasourceType: r.DatasourceType,
		DatasourceId:   r.DatasourceId,
		MetricName:     r.MetricName,
		PromQL:         r.PromQL,
		Labels:         r.Labels,
		EvalInterval:   r.EvalInterval,
		UpdateAt:       time.Now().Unix(),
		UpdateBy:       r.UpdateBy,
		CreateAt:       time.Now().Unix(),
		CreateBy:       r.UpdateBy,
		Enabled:        r.Enabled,
		RuleGroupId:    r.RuleGroupId,
	}

	// Validate the rule
	err := data.Validate()
	if err != nil {
		return nil, err
	}

	err = rs.ctx.DB.RecordingRule().Create(data)
	if err != nil {
		return nil, err
	}

	// 判断当前节点角色
	if *r.GetEnabled() {
		if alert.IsLeader() {
			// Leader: 直接启动评估协程
			alert.RecordingRule.Submit(data)
		} else {
			// Follower: 发布 Redis 消息通知 Leader
			tools.PublishReloadMessage(rs.ctx.Ctx, client.Redis, tools.ChannelRecordingRuleReload, tools.ReloadMessage{
				Action:   tools.ActionCreate,
				ID:       data.RuleId,
				TenantID: data.TenantId,
			})
		}
	}

	return nil, nil
}

func (rs recordingRuleService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRecordingRuleUpdate)

	// Get the old rule
	oldRule := models.RecordingRule{}
	rs.ctx.DB.DB().Model(&models.RecordingRule{}).
		Where("tenant_id = ? AND rule_id = ?", r.TenantId, r.RuleId).
		First(&oldRule)

	// 判断当前状态变化
	var action string

	if *oldRule.GetEnabled() == true && *r.GetEnabled() == false {
		action = tools.ActionDisable
	} else if *oldRule.GetEnabled() == false && *r.GetEnabled() == true {
		action = tools.ActionEnable
	} else if *oldRule.GetEnabled() == true && *r.GetEnabled() == true {
		action = tools.ActionUpdate
	}

	// 如果指定了ruleGroupId，验证规则组是否存在
	if r.RuleGroupId > 0 {
		var group models.RecordingRuleGroup
		err := rs.ctx.DB.DB().Model(&models.RecordingRuleGroup{}).
			Where("tenant_id = ? AND id = ?", r.TenantId, r.RuleGroupId).
			First(&group).Error
		if err != nil {
			return nil, fmt.Errorf("规则组不存在")
		}
	}

	data := models.RecordingRule{
		TenantId:       r.TenantId,
		RuleId:         r.RuleId,
		DatasourceType: r.DatasourceType,
		DatasourceId:   r.DatasourceId,
		MetricName:     r.MetricName,
		PromQL:         r.PromQL,
		Labels:         r.Labels,
		EvalInterval:   r.EvalInterval,
		UpdateAt:       time.Now().Unix(),
		UpdateBy:       r.UpdateBy,
		Enabled:        r.Enabled,
		RuleGroupId:    r.RuleGroupId,
	}

	// Validate the rule
	err := data.Validate()
	if err != nil {
		return nil, err
	}

	// 更新数据
	err = rs.ctx.DB.RecordingRule().Update(data)
	if err != nil {
		return nil, err
	}

	// 判断当前节点角色并处理
	if action != "" {
		if alert.IsLeader() {
			// Leader: 直接操作协程
			if action == tools.ActionDisable || action == tools.ActionUpdate {
				alert.RecordingRule.Stop(r.RuleId)
			}
			if (action == tools.ActionEnable || action == tools.ActionUpdate) && *r.GetEnabled() {
				alert.RecordingRule.Submit(data)
			}
		} else {
			// Follower: 发布消息通知 Leader
			tools.PublishReloadMessage(rs.ctx.Ctx, client.Redis, tools.ChannelRecordingRuleReload, tools.ReloadMessage{
				Action:   action,
				ID:       r.RuleId,
				TenantID: r.TenantId,
			})
		}
	}

	// 如果禁用，删除缓存
	if !*r.GetEnabled() {
		// TODO: 删除缓存
	}

	return nil, nil
}

func (rs recordingRuleService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRecordingRuleQuery)
	info, err := rs.ctx.DB.RecordingRule().Get(r.TenantId, r.RuleId)
	if err != nil {
		return nil, err
	}

	err = rs.ctx.DB.RecordingRule().Delete(r.TenantId, r.RuleId)
	if err != nil {
		return nil, err
	}

	// 判断当前节点角色
	if *info.GetEnabled() {
		if alert.IsLeader() {
			// Leader: 直接停止协程
			alert.RecordingRule.Stop(r.RuleId)
		} else {
			// Follower: 发布消息通知 Leader
			tools.PublishReloadMessage(rs.ctx.Ctx, client.Redis, tools.ChannelRecordingRuleReload, tools.ReloadMessage{
				Action:   tools.ActionDelete,
				ID:       r.RuleId,
				TenantID: r.TenantId,
			})
		}
	}

	// 删除缓存
	// TODO: 删除缓存

	return nil, nil
}

func (rs recordingRuleService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRecordingRuleQuery)
	data, count, err := rs.ctx.DB.RecordingRule().List(r.TenantId, r.DatasourceType, r.Query, r.Status, r.Page)
	if err != nil {
		return nil, err
	}

	// 如果指定了ruleGroupId，过滤结果
	if r.RuleGroupId > 0 {
		var filtered []models.RecordingRule
		for _, item := range data {
			if item.RuleGroupId == r.RuleGroupId {
				filtered = append(filtered, item)
			}
		}
		data = filtered
		count = int64(len(data))
	}

	return types.ResponseRecordingRuleList{
		List: data,
		Page: models.Page{
			Total: count,
			Index: r.Page.Index,
			Size:  r.Page.Size,
		},
	}, nil
}

func (rs recordingRuleService) Get(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRecordingRuleQuery)
	data, err := rs.ctx.DB.RecordingRule().Get(r.TenantId, r.RuleId)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (rs recordingRuleService) ChangeStatus(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRecordingRuleChangeStatus)
	var action string
	switch *r.GetEnabled() {
	case true:
		action = tools.ActionEnable
	case false:
		action = tools.ActionDisable
	}

	// 更新数据库
	err := rs.ctx.DB.RecordingRule().ChangeStatus(r.TenantId, r.RuleId, r.GetEnabled())
	if err != nil {
		return nil, err
	}

	// 判断当前节点角色
	rule := rs.ctx.DB.RecordingRule().GetRuleObject(r.RuleId)
	if alert.LeaderElector != nil && alert.LeaderElector.IsLeader() {
		// Leader: 直接操作协程
		switch *r.GetEnabled() {
		case true:
			var enable = true
			newRule := rule
			newRule.Enabled = &enable
			alert.RecordingRule.Submit(newRule)
		case false:
			alert.RecordingRule.Stop(r.RuleId)
		}
	} else {
		// Follower: 发布消息通知 Leader
		tools.PublishReloadMessage(rs.ctx.Ctx, client.Redis, tools.ChannelRecordingRuleReload, tools.ReloadMessage{
			Action:   action,
			ID:       r.RuleId,
			TenantID: r.TenantId,
		})
	}

	return nil, nil
}
