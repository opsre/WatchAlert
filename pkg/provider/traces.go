package provider

import (
	"fmt"
	"watchAlert/internal/models"
)

const (
	JaegerDsProviderName string = "Jaeger"
)

type TracesFactoryProvider interface {
	Query(options TraceQueryOptions) ([]Traces, error)
	Check() (bool, error)
	GetJaegerService() (JaegerServiceData, error)
	GetExternalLabels() map[string]interface{}
}

type TraceQueryOptions struct {
	Tags    string `json:"tags,omitempty"`    // 查询标签
	Service string `json:"service,omitempty"` // 服务名称
	Limit   int64  `json:"limit,omitempty"`   // 要返回的最大条目数
	StartAt int64  `json:"startAt,omitempty"` // 查询的开始时间，以微秒 Unix 表示。
	EndAt   int64  `json:"endAt,omitempty"`   // 查询的结束时间，以微秒 Unix 表示。
}

type Traces struct {
	Service string
	TraceId string
}

func (t Traces) GetFingerprint() string {
	return t.TraceId
}

func (t Traces) GetMetric() map[string]interface{} {
	return map[string]interface{}{
		"service": t.Service,
		"trace":   t.TraceId,
	}
}

func (t Traces) GetAnnotations(rule models.AlertRule, ds models.AlertDataSource) string {
	return fmt.Sprintf("服务: %s 链路中存在异常状态码接口\nJaeger URL: %s/trace/%s\n\n详情查看 Jaeger Trace ⬆️", rule.JaegerConfig.Service, ds.HTTP.URL, t.TraceId)
}
