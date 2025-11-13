package provider

import (
	"fmt"
	"strconv"
	"time"
	"watchAlert/pkg/tools"
)

const (
	PrometheusDsProvider      string = "Prometheus"
	VictoriaMetricsDsProvider string = "VictoriaMetrics"
)

type MetricsFactoryProvider interface {
	Query(promQL string) ([]Metrics, error)
	QueryRange(promQL string, start, end time.Time, step time.Duration) ([]Metrics, error)
	Check() (bool, error)
	GetExternalLabels() map[string]interface{}
}

type Metrics struct {
	Metric    map[string]interface{}
	Value     float64
	Timestamp float64
}

func (m Metrics) GetFingerprint() string {
	var labels = m.Metric

	if len(labels) == 0 {
		return strconv.FormatUint(tools.HashNew(), 10)
	}

	var result uint64
	for labelName, labelValue := range labels {
		sum := tools.HashNew()
		sum = tools.HashAdd(sum, labelName)
		sum = tools.HashAdd(sum, fmt.Sprintf("%v", labelValue))
		result ^= sum
	}

	return strconv.FormatUint(result, 10)
}

func (m Metrics) GetMetric() map[string]interface{} {
	return m.Metric
}

func (m Metrics) GetValue() float64 {
	return m.Value
}
