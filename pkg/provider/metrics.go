package provider

import (
	"fmt"
	"strconv"
	"watchAlert/pkg/tools"
)

const (
	PrometheusDsProvider      string = "Prometheus"
	VictoriaMetricsDsProvider string = "VictoriaMetrics"
)

type MetricsFactoryProvider interface {
	Query(promQL string) ([]Metrics, error)
	Check() (bool, error)
	GetExternalLabels() map[string]interface{}
}

type Metrics struct {
	Metric    map[string]interface{}
	Value     float64
	Timestamp float64
}

func (m Metrics) GetFingerprint() string {
	if len(m.Metric) == 0 {
		return strconv.FormatUint(tools.HashNew(), 10)
	}

	var result uint64
	for labelName, labelValue := range m.Metric {
		sum := tools.HashNew()
		sum = tools.HashAdd(sum, labelName)
		sum = tools.HashAdd(sum, fmt.Sprintf("%v", labelValue))
		result ^= sum
	}

	return strconv.FormatUint(result, 10)
}

func (m Metrics) GetMetric() *map[string]interface{} {
	m.Metric["value"] = m.Value
	return &m.Metric
}

func (m Metrics) GetValue() *float64 {
	return &m.Value
}
