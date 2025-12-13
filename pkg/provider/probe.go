package provider

const (
	ICMPEndpointProvider string = "ICMP"
	HTTPEndpointProvider string = "HTTP"
	TCPEndpointProvider  string = "TCP"
	SSLEndpointProvider  string = "SSL"
)

// ProbeMetric 探测指标结构
type ProbeMetric struct {
	Name      string         `json:"name"`
	Help      string         `json:"help"`
	Type      string         `json:"type"`
	Labels    map[string]any `json:"labels"`
	Value     float64        `json:"value"`
	Timestamp int64          `json:"timestamp"`
}

// MetricsAwareProbe 支持直接返回指标的探测接口（通用）
type MetricsAwareProbe interface {
	PilotWithMetrics(option EndpointOption, ruleInfo ProbeRuleInfo) []ProbeMetric
}

// ProbeRuleInfo 探测规则信息，用于生成指标标签（通用）
type ProbeRuleInfo struct {
	TenantID string `json:"tenant_id"`
	RuleID   string `json:"rule_id"`
	RuleName string `json:"rule_name"`
	RuleType string `json:"rule_type"`
	Endpoint string `json:"endpoint"`
}

// BoolToFloat 将布尔值转换为浮点数
func BoolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

type EndpointOption struct {
	Endpoint string `json:"endpoint"`
	Timeout  int    `json:"timeout"`
	HTTP     Ehttp  `json:"http"`
	ICMP     Eicmp  `json:"icmp"`
}

type Ehttp struct {
	Method string            `json:"method"`
	Header map[string]string `json:"header"`
	Body   string            `json:"body"`
}

type Eicmp struct {
	Interval int `json:"interval"`
	Count    int `json:"count"`
}

type PingerInformation struct {
	Address string `json:"address"`
	// 发送的数据包数量
	PacketsSent int `json:"packetsSent"`
	// 成功接收到的数据包数量
	PacketsRecv int `json:"packetsRecv"`
	// 丢包率的百分比
	PacketLoss float64 `json:"packetLoss"`
	// 目标主机的地址（例如域名或 IP 地址）
	Addr string `json:"addr"`
	// 目标主机的 IP 地址
	IPAddr string `json:"IPAddr"`
	// 最短的 RTT 时间, ms
	MinRtt float64 `json:"minRtt"`
	// 最长的 RTT 时间, ms
	MaxRtt float64 `json:"maxRtt"`
	// 平均 RTT 时间, ms
	AvgRtt float64 `json:"avgRtt"`
}

// copyLabelsMap 复制标签映射
func copyLabelsMap(labels map[string]any) map[string]any {
	copied := make(map[string]any)
	for k, v := range labels {
		copied[k] = v
	}
	return copied
}

// boolToFloat 将布尔值转换为浮点数
func boolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}
