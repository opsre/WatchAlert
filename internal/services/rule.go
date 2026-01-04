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

	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logc"
	"gopkg.in/yaml.v3"
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
	Change(req interface{}) (interface{}, interface{})
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
		UpdateAt:             time.Now().Unix(),
		UpdateBy:             r.UpdateBy,
		Enabled:              r.Enabled,
	}

	err := rs.ctx.DB.Rule().Create(data)
	if err != nil {
		return nil, err
	}

	// 判断当前节点角色
	if *r.GetEnabled() {
		if alert.IsLeader() {
			// Leader: 直接启动评估协程
			alert.AlertRule.Submit(data)
		} else {
			// Follower: 发布 Redis 消息通知 Leader
			tools.PublishReloadMessage(rs.ctx.Ctx, client.Redis, tools.ChannelRuleReload, tools.ReloadMessage{
				Action:   tools.ActionCreate,
				ID:       data.RuleId,
				TenantID: data.TenantId,
				Name:     data.RuleName,
			})
		}
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
	var action string
	if *oldRule.Enabled == true && *r.Enabled == false {
		action = tools.ActionDisable
	} else if *oldRule.Enabled == false && *r.Enabled == true {
		action = tools.ActionEnable
	} else if *oldRule.Enabled == true && *r.Enabled == true {
		action = tools.ActionUpdate
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
		UpdateAt:             time.Now().Unix(),
		UpdateBy:             r.UpdateBy,
		Enabled:              r.Enabled,
	}

	// 更新数据
	err := rs.ctx.DB.Rule().Update(data)
	if err != nil {
		return nil, err
	}

	// 判断当前节点角色并处理
	if action != "" {
		if alert.IsLeader() {
			// Leader: 直接操作协程
			if action == tools.ActionDisable || action == tools.ActionUpdate {
				alert.AlertRule.Stop(r.RuleId)
			}
			if (action == tools.ActionEnable || action == tools.ActionUpdate) && *r.GetEnabled() {
				alert.AlertRule.Submit(data)
			}
		} else {
			// Follower: 发布消息通知 Leader
			tools.PublishReloadMessage(rs.ctx.Ctx, client.Redis, tools.ChannelRuleReload, tools.ReloadMessage{
				Action:   action,
				ID:       r.RuleId,
				TenantID: r.TenantId,
				Name:     r.RuleName,
			})
		}
	}

	// 如果禁用，删除缓存
	if !*r.GetEnabled() {
		fingerprints := rs.ctx.Redis.Alert().GetFingerprintsByRuleId(r.TenantId, r.FaultCenterId, r.RuleId)
		for _, fingerprint := range fingerprints {
			rs.ctx.Redis.Alert().RemoveAlertEvent(r.TenantId, r.FaultCenterId, fingerprint)
		}
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

	// 判断当前节点角色
	if *info.GetEnabled() {
		if alert.IsLeader() {
			// Leader: 直接停止协程
			alert.AlertRule.Stop(r.RuleId)
		} else {
			// Follower: 发布消息通知 Leader
			tools.PublishReloadMessage(rs.ctx.Ctx, client.Redis, tools.ChannelRuleReload, tools.ReloadMessage{
				Action:   tools.ActionDelete,
				ID:       r.RuleId,
				TenantID: r.TenantId,
				Name:     info.RuleName,
			})
		}
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
	var action string
	switch *r.GetEnabled() {
	case true:
		action = tools.ActionEnable
	case false:
		action = tools.ActionDisable
		// 删除缓存
		fingerprints := rs.ctx.Redis.Alert().GetFingerprintsByRuleId(r.TenantId, r.FaultCenterId, r.RuleId)
		for _, fingerprint := range fingerprints {
			rs.ctx.Redis.Alert().RemoveAlertEvent(r.TenantId, r.FaultCenterId, fingerprint)
		}
	}

	// 更新数据库
	err := rs.ctx.DB.Rule().ChangeStatus(r.TenantId, r.RuleGroupId, r.RuleId, r.GetEnabled())
	if err != nil {
		return nil, err
	}

	// 判断当前节点角色
	rule := rs.ctx.DB.Rule().GetRuleObject(r.RuleId)
	if alert.LeaderElector != nil && alert.LeaderElector.IsLeader() {
		// Leader: 直接操作协程
		switch *r.GetEnabled() {
		case true:
			var enable = true
			newRule := rule
			newRule.Enabled = &enable
			alert.AlertRule.Submit(newRule)
		case false:
			alert.AlertRule.Stop(r.RuleId)
		}
	} else {
		// Follower: 发布消息通知 Leader
		tools.PublishReloadMessage(rs.ctx.Ctx, client.Redis, tools.ChannelRuleReload, tools.ReloadMessage{
			Action:   action,
			ID:       r.RuleId,
			TenantID: r.TenantId,
			Name:     rule.RuleName,
		})
	}

	return nil, nil
}

func (rs ruleService) Import(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRuleImport)
	var (
		rules []types.RequestRuleCreate
		// 导入的规则默认为关闭状态
		disable bool
	)

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
		err := sonic.Unmarshal([]byte(r.Rules), &rules)
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
			RuleGroupId:          r.RuleGroupId,
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
			Enabled:              &disable,
		})
		if err != nil {
			logc.Errorf(rs.ctx.Ctx, err.Error())
			continue
		}
	}

	return nil, nil
}

func (rs ruleService) Change(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestRuleChange)

	if len(r.RuleIds) == 0 {
		return nil, fmt.Errorf("rule_ids 不能为空")
	}

	if len(r.Change) == 0 {
		return nil, fmt.Errorf("change 参数不能为空")
	}

	// 遍历所有要更新的规则ID
	for _, ruleId := range r.RuleIds {
		// 获取当前规则
		rule := models.AlertRule{}
		err := rs.ctx.DB.DB().Model(&models.AlertRule{}).
			Where("tenant_id = ? AND rule_id = ?", r.TenantId, ruleId).
			First(&rule).Error
		if err != nil {
			return nil, fmt.Errorf("获取规则失败, ruleId: %s, error: %v", ruleId, err)
		}

		rule.UpdateAt = time.Now().Unix()
		rule.UpdateBy = r.UpdateBy

		// 根据 change 参数动态更新字段
		for field, value := range r.Change {
			switch field {
			case "rule_group_id":
				if v, ok := value.(string); ok {
					rule.RuleGroupId = v
				} else {
					return nil, fmt.Errorf("字段 %s 的值类型错误", field)
				}
				break
			case "datasource_ids":
				if v, ok := value.([]interface{}); ok {
					var datasourceIds []string
					for _, val := range v {
						datasourceIds = append(datasourceIds, fmt.Sprintf("%v", val))
					}
					rule.DatasourceIdList = datasourceIds
				} else if v, ok := value.([]string); ok {
					rule.DatasourceIdList = v
				} else {
					return nil, fmt.Errorf("字段 %s 的值类型错误", field)
				}
				break
			case "fault_center_id":
				if v, ok := value.(string); ok {
					rule.FaultCenterId = v
				} else {
					return nil, fmt.Errorf("필드 %s 의 값 타입이 잘못되었습니다", field)
				}
				break
			case "enabled":
				if v, ok := value.(bool); ok {
					isEnabled := v
					rule.Enabled = &isEnabled
				} else {
					return nil, fmt.Errorf("필드 %s 의 값 타입이 잘못되었습니다", field)
				}
				break
			default:
				return nil, fmt.Errorf("지원되지 않는 필드: %s", field)
			}
		}

		// 데이터베이스 업데이트
		err = rs.ctx.DB.Rule().Update(rule)
		if err != nil {
			return nil, fmt.Errorf("규칙 업데이트 실패, ruleId: %s, 오류: %v", ruleId, err)
		}

		// 규칙 활성화 상태에 따라 적절한 처리
		if *rule.Enabled {
			if alert.IsLeader() {
				// 리더: 기존 평가 고루틴을 중지하고 새로운 것을 시작합니다
				alert.AlertRule.Stop(ruleId)
				alert.AlertRule.Submit(rule)
			} else {
				tools.PublishReloadMessage(rs.ctx.Ctx, client.Redis, tools.ChannelRuleReload, tools.ReloadMessage{
					Action:   tools.ActionUpdate,
					ID:       rule.RuleId,
					TenantID: rule.TenantId,
					Name:     rule.RuleName,
				})
			}
		} else {
			// 删除缓存
			fingerprints := rs.ctx.Redis.Alert().GetFingerprintsByRuleId(r.TenantId, rule.FaultCenterId, rule.RuleId)
			for _, fingerprint := range fingerprints {
				rs.ctx.Redis.Alert().RemoveAlertEvent(r.TenantId, rule.FaultCenterId, fingerprint)
			}

			if alert.IsLeader() {
				// 리더: 평가 고루틴 중지
				alert.AlertRule.Stop(ruleId)
			} else {
				// 팔로워: 비활성화 메시지를 발행하여 리더에게 알림
				tools.PublishReloadMessage(rs.ctx.Ctx, client.Redis, tools.ChannelRuleReload, tools.ReloadMessage{
					Action:   tools.ActionDisable,
					ID:       rule.RuleId,
					TenantID: rule.TenantId,
					Name:     rule.RuleName,
				})
			}
		}
	}

	return nil, nil
}
