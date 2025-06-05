package test

import (
	"context"
	"fmt"
	"testing"
	"watchAlert/internal/models"
	"watchAlert/pkg/provider"
)

func TestNewClickHouseClient(t *testing.T) {
	cli, err := provider.NewClickHouseClient(context.Background(), models.AlertDataSource{
		Auth: models.Auth{
			User: "root",
			Pass: "shimo",
		},
		ClickHouseConfig: models.DsClickHouseConfig{
			Addr:    "172.17.83.177:19000",
			Timeout: 10,
		},
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	cli.Check()

	cli.Query(provider.LogQueryOptions{ClickHouse: provider.ClickHouse{
		Query: "SELECT * FROM zprod.express WHERE `_time_second_`='2025-05-29 00:00:00';",
	}})
}
