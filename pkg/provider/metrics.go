package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"watchAlert/pkg/tools"
)

const (
	PrometheusDsProvider string = "Prometheus"
)

type MetricsFactoryProvider interface {
	Query(promQL string) ([]Metrics, error)
	QueryRange(promQL string, start, end time.Time, step time.Duration) ([]Metrics, error)
	Check() (bool, error)
	GetExternalLabels() map[string]interface{}
	Write(ctx context.Context, result []Metrics, labels map[string]string) error
}

type Metrics struct {
	Name      string         `json:"name"`
	Help      string         `json:"help"`
	Labels    map[string]any `json:"labels"`
	Value     float64        `json:"value"`
	Timestamp int64          `json:"timestamp"`
}

func (m Metrics) GetFingerprint() string {
	var labels = m.Labels

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
	return m.Labels
}

func (m Metrics) GetValue() float64 {
	return m.Value
}
