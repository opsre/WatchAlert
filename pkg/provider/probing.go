package provider

import (
	"watchAlert/pkg/tools"
)

const (
	ICMPEndpointProvider string = "ICMP"
	HTTPEndpointProvider string = "HTTP"
	TCPEndpointProvider  string = "TCP"
	SSLEndpointProvider  string = "SSL"
)

type EndpointFactoryProvider interface {
	Pilot(option EndpointOption) (EndpointValue, error)
}

type EndpointValue map[string]any

func (e EndpointValue) GetLabels() map[string]interface{} {
	return e
}

func (e EndpointValue) GetFingerprint() string {
	return tools.Md5Hash([]byte(e["address"].(string)))
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

type HttperInformation struct {
	Address string `json:"address"`
	// 状态码
	StatusCode float64 `json:"statusCode"`
	// 响应时间, ms
	Latency float64 `json:"latency"`
	// Headers
	Headers map[string][]string
}

type SslInformation struct {
	Address string `json:"address"`
	// 证书开始时间
	StartTime string
	// 证书过期时间
	ExpireTime string
	// 剩余有效天数
	TimeRemaining float64
	// 响应时间（毫秒）
	ResponseTime float64
}

type TcperInformation struct {
	// 目标地址
	Address string
	// 连接响应时间
	ResponseTime float64
	// 是否拨测成功
	IsSuccessful bool
	// 错误信息（拨测失败时）
	ErrorMessage string
}
