package probe

// ProbeMetric 探测指标
type ProbeMetric struct {
	Name      string         `json:"name"`
	Help      string         `json:"help"`
	Type      string         `json:"type"`
	Labels    map[string]any `json:"labels"`
	Value     float64        `json:"value"`
	Timestamp int64          `json:"timestamp"`
}
