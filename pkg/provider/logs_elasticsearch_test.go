package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
	"watchAlert/internal/models"
)

func TestNewElasticSearchClient(t *testing.T) {
	client, err := NewElasticSearchClient(context.Background(), models.AlertDataSource{})
	if err != nil {
		logrus.Errorf("client -> %s", err.Error())
		return
	}

	client.Query(LogQueryOptions{})
}

func TestElasticsearch_GetIndexName(t *testing.T) {
	var ess = []Elasticsearch{
		{
			Index: "test.2000.10.23",
		},
		{
			Index: "test.YYYYMMdd",
		},
		{
			Index: "test.YYYY-MM-dd",
		},
		{
			Index: "test.YYYY_MM_dd",
		},
	}
	for _, es := range ess {
		/*
			test.2000.10.23
			test.20250223
			test.2025-02-23
			test.2025_02_23
		*/
		fmt.Println(es.GetIndexName())
	}
}

func TestElasticSearch_Query(t *testing.T) {
	client, err := NewElasticSearchClient(context.Background(), models.AlertDataSource{
		HTTP: models.HTTP{
			URL: "http://192.168.1.190:9200",
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	query, _, err := client.Query(LogQueryOptions{ElasticSearch: Elasticsearch{
		Index:     "test-2024-05.20",
		QueryType: "RawJson",
		RawJson:   `{"match_all":{}}`,
	}})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	json, _ := json.Marshal(query[0].Message)
	fmt.Println("query->", string(json))

}
