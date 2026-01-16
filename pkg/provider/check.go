package provider

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"watchAlert/internal/models"
)

// HealthChecker 统一健康检查接口
type HealthChecker interface {
	Check() (bool, error)
}

// ClientFactory 客户端工厂函数类型
type ClientFactory func(models.AlertDataSource) (HealthChecker, error)

// 注册所有数据源类型的工厂方法
var datasourceFactories = map[string]ClientFactory{
	"Prometheus": func(ds models.AlertDataSource) (HealthChecker, error) {
		return NewPrometheusClient(ds)
	},
	"Kubernetes": func(ds models.AlertDataSource) (HealthChecker, error) {
		return NewKubernetesClient(context.Background(), ds.KubeConfig, ds.Labels)
	},
	"ElasticSearch": func(ds models.AlertDataSource) (HealthChecker, error) {
		return NewElasticSearchClient(context.Background(), ds)
	},
	"AliCloudSLS": func(ds models.AlertDataSource) (HealthChecker, error) {
		return NewAliCloudSlsClient(ds)
	},
	"Loki": func(ds models.AlertDataSource) (HealthChecker, error) {
		return NewLokiClient(ds)
	},
	"Jaeger": func(ds models.AlertDataSource) (HealthChecker, error) {
		return NewJaegerClient(ds)
	},
	"CloudWatch": func(ds models.AlertDataSource) (HealthChecker, error) {
		return &CloudWatchDummyChecker{}, nil
	},
	"VictoriaLogs": func(ds models.AlertDataSource) (HealthChecker, error) {
		return NewVictoriaLogsClient(context.Background(), ds)
	},
	"ClickHouse": func(ds models.AlertDataSource) (HealthChecker, error) {
		return NewClickHouseClient(context.Background(), ds)
	},
}

// CloudWatchDummyChecker 云监控哑检查器
type CloudWatchDummyChecker struct{}

func (c *CloudWatchDummyChecker) Check() (bool, error) {
	return true, nil
}

// CheckDatasourceHealth 统一健康检查入口
func CheckDatasourceHealth(datasource models.AlertDataSource) (bool, error) {
	// 获取对应的工厂方法
	factory, ok := datasourceFactories[datasource.Type]
	if !ok {
		err := fmt.Errorf("unsupported datasource type: %s", datasource.Type)
		logDatasourceError(datasource, err)
		return false, err
	}

	// 创建客户端
	client, err := factory(datasource)
	if err != nil {
		logDatasourceError(datasource, fmt.Errorf("client creation failed: %w", err))
		return false, err
	}

	// 执行健康检查
	healthy, err := client.Check()
	if err != nil || !healthy {
		logDatasourceError(datasource, fmt.Errorf("health check failed: %w", err))
		return false, err
	}

	return true, nil
}

// 统一日志记录方法
func logDatasourceError(ds models.AlertDataSource, err error) {
	logc.Errorf(context.Background(), "Datasource error",
		map[string]interface{}{
			"id":   ds.ID,
			"name": ds.Name,
			"type": ds.Type,
			"err":  err.Error(),
		})
}
