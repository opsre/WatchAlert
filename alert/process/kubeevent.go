package process

import (
	"crypto/md5"
	"encoding/hex"
	v1 "k8s.io/api/core/v1"
	"strings"
	"watchAlert/internal/ctx"
	"watchAlert/pkg/tools"
)

type KubernetesEvent struct {
	ctx   *ctx.Context
	event v1.Event
}

func KubernetesAlertEvent(ctx *ctx.Context, event v1.Event) KubernetesEvent {
	return KubernetesEvent{ctx: ctx, event: event}
}

func (a KubernetesEvent) GetFingerprint() string {
	h := md5.New()
	s := map[string]interface{}{
		"namespace": a.event.Namespace,
		"resource":  a.event.Reason,
		"podName":   a.event.InvolvedObject.Name,
	}

	h.Write(tools.JsonMarshalToByte(s))
	fingerprint := hex.EncodeToString(h.Sum(nil))
	return fingerprint
}

func (a KubernetesEvent) GetMetrics() map[string]interface{} {
	return map[string]interface{}{
		"namespace": a.event.Namespace,
		"resource":  a.event.Reason,
		"podName":   a.event.InvolvedObject.Name,
	}
}

// EvalKubeEvent 评估 Kubernetes 事件
type EvalKubeEvent struct {
	Reason string
	Filter []string
}

// FilterKubeEvent 过滤资源
func FilterKubeEvent(event *v1.EventList, filter []string) *v1.EventList {
	if filter == nil {
		return event
	}

	warningEvents := &v1.EventList{}
	for _, event := range event.Items {
		var found bool
		for _, f := range filter {
			if strings.Contains(event.InvolvedObject.Name, f) {
				found = true
				break
			}
		}

		if !found {
			warningEvents.Items = append(warningEvents.Items, event)
		}
	}

	return warningEvents
}
