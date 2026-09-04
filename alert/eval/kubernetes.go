package eval

import (
	"fmt"
	"strings"
	"watchAlert/alert/process"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/provider"

	"github.com/zeromicro/go-zero/core/logc"
)

func kubernetesEvent(ctx *ctx.Context, datasourceId, datasourceType string, rule models.AlertRule) []string {
	// 获取数据源实例信息
	datasourceObj, err := ctx.DB.Datasource().GetInstance(datasourceId)
	if err != nil {
		logc.Errorf(ctx.Ctx, "获取数据源实例失败, 规则ID: %s, 规则名称: %s, 数据源ID: %s, 错误: %v", rule.RuleId, rule.RuleName, datasourceId, err)
		return []string{}
	}

	pools := ctx.Redis.ProviderPools()
	cli, err := pools.GetClient(datasourceId)
	if err != nil {
		logc.Errorf(ctx.Ctx, "获取Kubernetes数据源客户端失败, 规则ID: %s, 规则名称: %s, 数据源ID: %s, 错误: %v", rule.RuleId, rule.RuleName, datasourceId, err)
		return []string{}
	}

	k8sClient := cli.(provider.KubernetesClient)
	externalLabels := k8sClient.GetExternalLabels()

	// 查询 Kubernetes 事件
	k8sEventMap, err := k8sClient.GetWarningEvent(rule.KubernetesConfig.Reason, rule.KubernetesConfig.Scope, rule.KubernetesConfig.Filter)
	if err != nil {
		logc.Errorf(ctx.Ctx, "获取Kubernetes警告事件失败, 规则ID: %s, 规则名称: %s, 数据源ID: %s, 原因: %s, 错误: %v", rule.RuleId, rule.RuleName, datasourceId, rule.KubernetesConfig.Reason, err)
		return []string{}
	}

	// 无事件返回
	if len(k8sEventMap) == 0 {
		return []string{}
	}

	// 遍历事件组，评估并生成告警
	curFingerprints := make([]string, 0, len(k8sEventMap))
	for _, eventItems := range k8sEventMap {
		for _, k8sEvent := range eventItems {
			fingerprint := k8sEvent.GetFingerprint()

			// 构建告警事件
			event := process.BuildEvent(rule, func() map[string]interface{} {
				metric := k8sEvent.GetMetrics()
				metric["rule_name"] = rule.RuleName
				metric["severity"] = rule.Severity
				metric["fingerprint"] = fingerprint
				for k, v := range externalLabels {
					metric[k] = v
				}
				for k, v := range rule.ExternalLabels {
					metric[k] = v
				}
				return metric
			})

			// 设置事件基本信息
			event.DatasourceId = datasourceId
			event.Fingerprint = fingerprint
			event.SearchQL = rule.KubernetesConfig.Resource

			// 构建注释信息
			var msgList []string
			for _, e := range eventItems {
				msg := strings.ReplaceAll(e.Message, "\"", "'")
				msgList = append(msgList, msg)
			}

			event.Annotations = fmt.Sprintf(
				"- 数据源: %s\n- 命名空间: %s\n- 资源类型: %s\n- 资源名称: %s\n- 事件类型: %s\n- 事件详情:\n%s",
				datasourceObj.Name,
				k8sEvent.Namespace,
				k8sEvent.InvolvedObject.Kind,
				k8sEvent.InvolvedObject.Name,
				k8sEvent.Reason,
				strings.Join(msgList, "\n"),
			)

			// 推送到故障中心
			process.PushEventToFaultCenter(ctx, &event)
			curFingerprints = append(curFingerprints, fingerprint)
		}

	}

	return curFingerprints
}
