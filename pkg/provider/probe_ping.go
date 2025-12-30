package provider

import (
	"fmt"
	"time"

	"github.com/go-ping/ping"
)

type Pinger struct{}

// NewMetricsAwarePinger 创建支持指标的ICMP探测器
func NewMetricsAwarePinger() MetricsAwareProbe {
	return Pinger{}
}

// PilotWithMetrics 执行ICMP探测并直接返回指标
func (p Pinger) PilotWithMetrics(option EndpointOption, ruleInfo ProbeRuleInfo) []ProbeMetric {
	timestamp := time.Now().Unix()

	// 执行ICMP探测
	var detail PingerInformation
	pinger, err := ping.NewPinger(option.Endpoint)
	if err != nil {
		// 返回失败指标
		return p.createFailureMetrics(ruleInfo, timestamp, fmt.Sprintf("New pinger error: %s", err.Error()))
	}
	pinger.SetPrivileged(true)

	// 请求次数
	pinger.Count = option.ICMP.Count
	// 请求间隔
	pinger.Interval = time.Second * time.Duration(option.ICMP.Interval)
	// 超时时间
	pinger.Timeout = time.Second * time.Duration(option.Timeout)

	pinger.OnFinish = func(stats *ping.Statistics) {
		detail = PingerInformation{
			Address:     stats.Addr,
			PacketsSent: stats.PacketsSent,
			PacketsRecv: stats.PacketsRecv,
			PacketLoss:  stats.PacketLoss,
			Addr:        stats.Addr,
			IPAddr:      stats.IPAddr.String(),
			MinRtt:      float64(stats.MinRtt.Milliseconds()),
			MaxRtt:      float64(stats.MaxRtt.Milliseconds()),
			AvgRtt:      float64(stats.AvgRtt.Milliseconds()),
		}
	}

	err = pinger.Run()
	if err != nil {
		// 返回失败指标
		return p.createFailureMetrics(ruleInfo, timestamp, fmt.Sprintf("Ping error: %s", err.Error()))
	}

	// 创建基础标签
	baseLabels := map[string]any{
		"tenant_id":  ruleInfo.TenantID,
		"rule_id":    ruleInfo.RuleID,
		"probe_name": ruleInfo.RuleName,
		"probe_type": ruleInfo.RuleType,
		"endpoint":   ruleInfo.Endpoint,
		"ip_addr":    detail.IPAddr,
	}

	// 创建ICMP指标
	metrics := []ProbeMetric{
		{
			Name:      "probe_icmp_packet_loss_percent",
			Help:      "ICMP packet loss percentage",
			Type:      "gauge",
			Labels:    copyLabelsMap(baseLabels),
			Value:     detail.PacketLoss,
			Timestamp: timestamp,
		},
		{
			Name:      "probe_icmp_rtt_min_ms",
			Help:      "ICMP minimum round trip time in milliseconds",
			Type:      "gauge",
			Labels:    copyLabelsMap(baseLabels),
			Value:     detail.MinRtt,
			Timestamp: timestamp,
		},
		{
			Name:      "probe_icmp_rtt_max_ms",
			Help:      "ICMP maximum round trip time in milliseconds",
			Type:      "gauge",
			Labels:    copyLabelsMap(baseLabels),
			Value:     detail.MaxRtt,
			Timestamp: timestamp,
		},
		{
			Name:      "probe_icmp_rtt_avg_ms",
			Help:      "ICMP average round trip time in milliseconds",
			Type:      "gauge",
			Labels:    copyLabelsMap(baseLabels),
			Value:     detail.AvgRtt,
			Timestamp: timestamp,
		},
		{
			Name:      "probe_icmp_packets_sent_total",
			Help:      "Total ICMP packets sent",
			Type:      "counter",
			Labels:    copyLabelsMap(baseLabels),
			Value:     float64(detail.PacketsSent),
			Timestamp: timestamp,
		},
		{
			Name:      "probe_icmp_packets_received_total",
			Help:      "Total ICMP packets received",
			Type:      "counter",
			Labels:    copyLabelsMap(baseLabels),
			Value:     float64(detail.PacketsRecv),
			Timestamp: timestamp,
		},
	}

	return metrics
}

// createFailureMetrics 创建失败时的指标
func (p Pinger) createFailureMetrics(ruleInfo ProbeRuleInfo, timestamp int64, errorMsg string) []ProbeMetric {
	baseLabels := map[string]any{
		"tenant_id":  ruleInfo.TenantID,
		"probe_id":   ruleInfo.RuleID,
		"probe_name": ruleInfo.RuleName,
		"probe_type": ruleInfo.RuleType,
		"endpoint":   ruleInfo.Endpoint,
		"error":      errorMsg,
	}

	return []ProbeMetric{
		{
			Name:      "probe_icmp_packet_loss_percent",
			Help:      "ICMP packet loss percentage",
			Type:      "gauge",
			Labels:    copyLabelsMap(baseLabels),
			Value:     100.0, // 完全失败
			Timestamp: timestamp,
		},
		{
			Name:      "probe_icmp_success",
			Help:      "ICMP probe success (1 for success, 0 for failure)",
			Type:      "gauge",
			Labels:    copyLabelsMap(baseLabels),
			Value:     0.0, // 失败
			Timestamp: timestamp,
		},
	}
}
