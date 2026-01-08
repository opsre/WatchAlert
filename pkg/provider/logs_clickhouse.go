package provider

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"watchAlert/internal/models"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/zeromicro/go-zero/core/logc"
)

type ClickHouseProvider struct {
	client         *sql.DB
	ExternalLabels map[string]interface{}
}

func NewClickHouseClient(ctx context.Context, ds models.AlertDataSource) (LogsFactoryProvider, error) {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{ds.ClickHouseConfig.Addr},
		Auth: clickhouse.Auth{
			Username: ds.Auth.User,
			Password: ds.Auth.Pass,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: time.Second * time.Duration(ds.ClickHouseConfig.Timeout),
	})
	if conn == nil {
		return nil, errors.New("clickhouse connection failed")
	}

	return ClickHouseProvider{
		client:         conn,
		ExternalLabels: ds.Labels,
	}, nil
}

func (c ClickHouseProvider) Query(options LogQueryOptions) (Logs, int, error) {
	rows, err := c.client.Query(options.ClickHouse.Query)
	if err != nil {
		return Logs{}, 0, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return Logs{}, 0, err
	}

	var messages []map[string]interface{}

	for rows.Next() {
		// 每次循环都重新绑定指针,因为 Scan 是通过指针写入数据的.
		var values = make([]interface{}, len(columns))
		for i := range columns {
			values[i] = new(interface{})
		}

		// 扫描数据到 values
		if err := rows.Scan(values...); err != nil {
			logc.Error(context.Background(), "clickhouse scan error:", err)
			return Logs{}, 0, err
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			// 取出指针指向的实际数据
			val := *(values[i].(*interface{}))
			if val == nil {
				entry[col] = ""
				continue
			}
			if b, ok := val.([]byte); ok {
				entry[col] = string(b)
			} else {
				entry[col] = val
			}
		}
		messages = append(messages, entry)
	}

	if err := rows.Err(); err != nil {
		return Logs{}, 0, err
	}

	return Logs{
		ProviderName: ClickHouseDsProviderName,
		Message:      messages,
	}, len(messages), nil
}

func (c ClickHouseProvider) Check() (bool, error) {
	err := c.client.Ping()
	if err != nil {
		return false, errors.New("check clickhouse datasource is unhealthy")
	}

	return true, nil
}

func (c ClickHouseProvider) GetExternalLabels() map[string]interface{} {
	return c.ExternalLabels
}
