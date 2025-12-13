package probe

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"watchAlert/pkg/provider"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"
)

// MetricsWriterConfig 指标写入器配置
type MetricsWriterConfig struct {
	Endpoint string `json:"endpoint"` // 写入端点
	Username string `json:"username"` // 认证用户名
	Password string `json:"password"` // 认证密码
}

// MetricsWriter 指标写入器接口
type MetricsWriter interface {
	// WriteMetrics 写入指标
	WriteMetrics(ctx context.Context, metrics []provider.ProbeMetric) error
	// Close 关闭写入器
	Close() error
}

// Writer 指标写入器
type Writer struct {
	config     MetricsWriterConfig
	httpClient *http.Client
}

// NewWriter 创建写入器
func NewWriter(config MetricsWriterConfig) *Writer {
	// 验证端点URL
	if config.Endpoint == "" {
		panic("Write endpoint cannot be empty")
	}

	return &Writer{
		config:     config,
		httpClient: http.DefaultClient,
	}
}

// WriteMetrics 写入指标到 Remote Write API
func (w *Writer) WriteMetrics(ctx context.Context, metrics []provider.ProbeMetric) error {
	if len(metrics) == 0 {
		return nil
	}

	// 转换为Prometheus Remote Write格式
	writeRequest, err := w.convertToRemoteWriteFormat(metrics)
	if err != nil {
		return fmt.Errorf("转换为Remote Write格式失败: %w", err)
	}

	// 序列化为Protobuf
	data, err := proto.Marshal(writeRequest)
	if err != nil {
		return fmt.Errorf("序列化Protobuf失败: %w", err)
	}

	// Snappy压缩
	compressed := snappy.Encode(nil, data)

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.config.Endpoint, bytes.NewBuffer(compressed))
	if err != nil {
		return fmt.Errorf("创建Prometheus请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/x-protobuf")
	req.Header.Set("Content-Encoding", "snappy")
	req.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")

	// 设置认证
	if w.config.Username != "" && w.config.Password != "" {
		req.SetBasicAuth(w.config.Username, w.config.Password)
	}

	// 发送请求
	resp, err := w.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("发送Prometheus请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	// 处理错误响应
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("prometheus写入失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
}

// convertToRemoteWriteFormat 转换为 Remote Write格式
func (w *Writer) convertToRemoteWriteFormat(metrics []provider.ProbeMetric) (*prompb.WriteRequest, error) {
	var timeSeries []prompb.TimeSeries

	for _, metric := range metrics {
		// 构建标签
		labels := []prompb.Label{
			{Name: "__name__", Value: metric.Name},
		}

		// 添加其他标签
		for k, v := range metric.Labels {
			labels = append(labels, prompb.Label{
				Name:  k,
				Value: fmt.Sprintf("%v", v),
			})
		}

		// 时间戳转换：Unix毫秒（注意！不是纳秒）
		timestampMs := metric.Timestamp * 1000
		if metric.Timestamp > 1e12 { // 如果已经是毫秒级时间戳
			timestampMs = metric.Timestamp
		}

		// 创建样本
		sample := prompb.Sample{
			Value:     metric.Value,
			Timestamp: timestampMs,
		}

		// 创建时间序列
		ts := prompb.TimeSeries{
			Labels:  labels,
			Samples: []prompb.Sample{sample},
		}

		timeSeries = append(timeSeries, ts)
	}

	// 构造WriteRequest
	writeReq := &prompb.WriteRequest{
		Timeseries: timeSeries,
	}

	return writeReq, nil
}

// Close 关闭写入器
func (w *Writer) Close() error {
	return nil
}
