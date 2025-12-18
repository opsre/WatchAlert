package probe

import (
	"context"
	"fmt"
	"sync"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/provider"

	"github.com/zeromicro/go-zero/core/logc"
	"golang.org/x/sync/errgroup"
)

// ProbeService 拨测服务
type ProbeService struct {
	ctx           *ctx.Context
	watchCtxMap   map[string]context.CancelFunc
	mu            sync.RWMutex
	writerCache   map[string]MetricsWriter // 数据源写入器缓存，key为datasourceId
	writerCacheMu sync.RWMutex             // 写入器缓存锁
}

// NewProbeService 创建新的拨测服务
func NewProbeService(ctx *ctx.Context) *ProbeService {
	return &ProbeService{
		ctx:         ctx,
		watchCtxMap: make(map[string]context.CancelFunc),
		writerCache: make(map[string]MetricsWriter),
	}
}

// Add 添加拨测规则
func (s *ProbeService) Add(rule models.ProbeRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查规则是否已存在
	if _, exists := s.watchCtxMap[rule.RuleId]; exists {
		return fmt.Errorf("rule %s already exists", rule.RuleId)
	}

	c, cancel := context.WithCancel(s.ctx.Ctx)
	s.watchCtxMap[rule.RuleId] = cancel

	// 启动拨测协程
	go s.runProbing(c, rule)

	logc.Infof(s.ctx.Ctx, "Added probing rule: %s (%s)", rule.RuleName, rule.RuleType)
	return nil
}

// Stop 停止指定规则的拨测
func (s *ProbeService) Stop(ruleID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cancel, exists := s.watchCtxMap[ruleID]
	if !exists {
		return fmt.Errorf("rule %s not found", ruleID)
	}

	cancel()
	delete(s.watchCtxMap, ruleID)
	return nil
}

// StopAll 停止所有拨测任务
func (s *ProbeService) StopAll() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	count := len(s.watchCtxMap)
	if count == 0 {
		return nil
	}

	logc.Infof(s.ctx.Ctx, "Stopping %d probing tasks...", count)

	// 取消所有拨测任务
	for ruleID, cancel := range s.watchCtxMap {
		cancel()
		delete(s.watchCtxMap, ruleID)
	}

	// 关闭所有缓存的指标写入器
	s.writerCacheMu.Lock()
	for datasourceId, writer := range s.writerCache {
		if err := writer.Close(); err != nil {
			logc.Errorf(s.ctx.Ctx, "关闭数据源%s的指标写入器失败: %v", datasourceId, err)
		}
	}
	s.writerCache = make(map[string]MetricsWriter) // 清空缓存
	s.writerCacheMu.Unlock()

	logc.Infof(s.ctx.Ctx, "All probing tasks stopped")
	return nil
}

// GetActiveRules 获取活跃规则数量
func (s *ProbeService) GetActiveRules() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.watchCtxMap)
}

// runProbing 运行拨测
func (s *ProbeService) runProbing(ctx context.Context, rule models.ProbeRule) {
	timer := time.NewTicker(time.Second * time.Duration(rule.ProbingEndpointConfig.Strategy.EvalInterval))
	defer timer.Stop()

	// 立即执行一次
	s.executeProbing(rule)

	for {
		select {
		case <-timer.C:
			s.executeProbing(rule)
		case <-ctx.Done():
			logc.Infof(s.ctx.Ctx, "Probing stopped for rule: %s", rule.RuleId)
			return
		}
	}
}

// executeProbing 执行拨测
func (s *ProbeService) executeProbing(rule models.ProbeRule) {
	// 执行拨测并获取指标
	metrics, err := s.executeProbeWithMetrics(rule)
	if err != nil {
		logc.Errorf(s.ctx.Ctx, "Probing failed for rule %s: %v", rule.RuleId, err)
		return
	}

	// 写入指标到数据源
	if len(metrics) > 0 && rule.DatasourceId != "" {
		if err := s.writeMetricsToDataSource(rule.DatasourceId, metrics); err != nil {
			logc.Errorf(s.ctx.Ctx, "Failed to write metrics for rule %s to datasource %s: %v",
				rule.RuleId, rule.DatasourceId, err)
		}
	}
}

// executeProbeWithMetrics 执行拨测并获取指标
func (s *ProbeService) executeProbeWithMetrics(rule models.ProbeRule) ([]provider.ProbeMetric, error) {
	config := rule.ProbingEndpointConfig

	// 准备规则信息
	ruleInfo := provider.ProbeRuleInfo{
		TenantID: rule.TenantId,
		RuleID:   rule.RuleId,
		RuleName: rule.RuleName,
		RuleType: rule.RuleType,
		Endpoint: config.Endpoint,
	}

	var metrics []provider.ProbeMetric

	// 根据协议类型选择相应的指标感知探测器
	switch rule.RuleType {
	case provider.HTTPEndpointProvider:
		httper := provider.NewMetricsAwareHTTPer()
		metrics = httper.PilotWithMetrics(provider.EndpointOption{
			Endpoint: config.Endpoint,
			Timeout:  config.Strategy.Timeout,
			HTTP: provider.Ehttp{
				Method: config.HTTP.Method,
				Header: config.HTTP.Header,
				Body:   config.HTTP.Body,
			},
		}, ruleInfo)

	case provider.ICMPEndpointProvider:
		pinger := provider.NewMetricsAwarePinger()
		metrics = pinger.PilotWithMetrics(provider.EndpointOption{
			Endpoint: config.Endpoint,
			Timeout:  config.Strategy.Timeout,
			ICMP: provider.Eicmp{
				Interval: config.ICMP.Interval,
				Count:    config.ICMP.Count,
			},
		}, ruleInfo)

	case provider.TCPEndpointProvider:
		tcper := provider.NewMetricsAwareTcper()
		metrics = tcper.PilotWithMetrics(provider.EndpointOption{
			Endpoint: config.Endpoint,
			Timeout:  config.Strategy.Timeout,
		}, ruleInfo)

	case provider.SSLEndpointProvider:
		ssler := provider.NewMetricsAwareSSLer()
		metrics = ssler.PilotWithMetrics(provider.EndpointOption{
			Endpoint: config.Endpoint,
			Timeout:  config.Strategy.Timeout,
		}, ruleInfo)

	default:
		return nil, fmt.Errorf("unsupported rule type: %s", rule.RuleType)
	}

	return metrics, nil
}

// writeMetricsToDataSource 根据datasourceId写入指标到对应数据源
func (s *ProbeService) writeMetricsToDataSource(datasourceId string, metrics []provider.ProbeMetric) error {
	// 获取或创建写入器
	writer, err := s.getOrCreateWriter(datasourceId)
	if err != nil {
		return fmt.Errorf("获取数据源写入器失败: %w", err)
	}

	if writer == nil {
		// 如果没有写入器，只记录日志
		for _, metric := range metrics {
			logc.Debugf(s.ctx.Ctx, "Metric (no writer for datasource %s): %s{%v} %f @%d",
				datasourceId, metric.Name, metric.Labels, metric.Value, metric.Timestamp)
		}
		return nil
	}

	// 使用写入器写入指标
	ctx, cancel := context.WithTimeout(s.ctx.Ctx, 30*time.Second)
	defer cancel()

	if err := writer.WriteMetrics(ctx, metrics); err != nil {
		return err
	}

	return nil
}

// getOrCreateWriter 获取或创建指定数据源的写入器
func (s *ProbeService) getOrCreateWriter(datasourceId string) (MetricsWriter, error) {
	// 先从缓存中查找
	s.writerCacheMu.RLock()
	if writer, exists := s.writerCache[datasourceId]; exists {
		s.writerCacheMu.RUnlock()
		return writer, nil
	}
	s.writerCacheMu.RUnlock()

	// 缓存中没有，需要创建新的写入器
	s.writerCacheMu.Lock()
	defer s.writerCacheMu.Unlock()

	// 双重检查，防止并发创建
	if writer, exists := s.writerCache[datasourceId]; exists {
		return writer, nil
	}

	// 从数据库获取数据源配置
	datasource, err := s.ctx.DB.Datasource().Get(datasourceId)
	if err != nil {
		return nil, fmt.Errorf("获取数据源配置失败: %w", err)
	}

	// 检查数据源是否启用
	if !datasource.GetEnabled() || datasource.Write.Enabled == "Off" {
		logc.Infof(s.ctx.Ctx, "数据源 %s 已禁用或未开启写入，跳过指标写入", datasourceId)
		return nil, nil
	}

	// 根据数据源类型创建写入器
	writer, err := s.createWriterFromDataSource(datasource)
	if err != nil {
		return nil, fmt.Errorf("创建数据源写入器失败: %w", err)
	}

	// 缓存写入器
	if writer != nil {
		s.writerCache[datasourceId] = writer
		logc.Infof(s.ctx.Ctx, "为数据源 %s (%s) 创建并缓存写入器",
			datasourceId, datasource.Name)
	}

	return writer, nil
}

// createWriterFromDataSource 根据数据源配置创建写入器
func (s *ProbeService) createWriterFromDataSource(datasource models.AlertDataSource) (MetricsWriter, error) {
	config := MetricsWriterConfig{
		Endpoint: datasource.Write.URL,
		Username: datasource.Auth.User,
		Password: datasource.Auth.Pass,
	}

	// 验证配置
	if err := ValidateWriteConfig(config); err != nil {
		return nil, fmt.Errorf("Write 配置验证失败: %w", err)
	}

	logc.Infof(s.ctx.Ctx, "创建 Write 写入器，端点: %s",
		config.Endpoint)

	return NewWriter(config), nil
}

// RemoveWriterFromCache 从缓存中移除指定数据源的写入器
func (s *ProbeService) RemoveWriterFromCache(datasourceId string) {
	s.writerCacheMu.Lock()
	defer s.writerCacheMu.Unlock()

	if writer, exists := s.writerCache[datasourceId]; exists {
		if err := writer.Close(); err != nil {
			logc.Errorf(s.ctx.Ctx, "关闭数据源 %s 的写入器失败: %v", datasourceId, err)
		}
		delete(s.writerCache, datasourceId)
		logc.Infof(s.ctx.Ctx, "从缓存中移除数据源 %s 的写入器", datasourceId)
	}
}

// ClearWriterCache 清空写入器缓存
func (s *ProbeService) ClearWriterCache() {
	s.writerCacheMu.Lock()
	defer s.writerCacheMu.Unlock()

	for datasourceId, writer := range s.writerCache {
		if err := writer.Close(); err != nil {
			logc.Errorf(s.ctx.Ctx, "关闭数据源%s的写入器失败: %v", datasourceId, err)
		}
	}
	s.writerCache = make(map[string]MetricsWriter)
	logc.Infof(s.ctx.Ctx, "已清空所有写入器缓存")
}

// RePushRule 重新推送规则
func (s *ProbeService) RePushRule() error {
	var ruleList []models.ProbeRule
	if err := s.ctx.DB.DB().Where("enabled = ?", true).Find(&ruleList).Error; err != nil {
		return fmt.Errorf("failed to fetch rules: %w", err)
	}

	g := new(errgroup.Group)
	for _, rule := range ruleList {
		rule := rule
		g.Go(func() error {
			if err := s.Add(rule); err != nil {
				return fmt.Errorf("failed to add rule %s: %w", rule.RuleId, err)
			}
			return nil
		})
	}

	return g.Wait()
}

// GetCachedWritersCount 获取缓存的写入器数量
func (s *ProbeService) GetCachedWritersCount() int {
	s.writerCacheMu.RLock()
	defer s.writerCacheMu.RUnlock()
	return len(s.writerCache)
}
