package provider

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

const (
	LokiDsProviderName          string = "Loki"
	AliCloudSLSDsProviderName   string = "AliCloudSLS"
	ElasticSearchDsProviderName string = "ElasticSearch"
)

type LogsFactoryProvider interface {
	Query(options LogQueryOptions) ([]Logs, int, error)
	Check() (bool, error)
	GetExternalLabels() map[string]interface{}
}

type LogQueryOptions struct {
	AliCloudSLS   AliCloudSLS
	Loki          Loki
	ElasticSearch Elasticsearch
	StartAt       interface{} // 查询的开始时间。
	EndAt         interface{} // 查询的结束时间。
}

type Loki struct {
	Query     string // 查询语句
	Direction string // 日志排序顺序，支持的值为forward或backward，默认为backward
	Limit     int64  // 要返回的最大条目数
}

type AliCloudSLS struct {
	Query    string // 查询语句
	Project  string // AliCloud SLS Project
	LogStore string // AliCloud SLS LogStore
}

type Elasticsearch struct {
	Index string // 索引名称
	// 索引选项
	IndexOption models.EsIndexOption
	QueryFilter []models.EsQueryFilter // 过滤条件
	// filter关系，与或非
	QueryFilterCondition models.EsFilterCondition
	// 查询类型，sql语句查询与条件查询
	QueryType models.EsQueryType
	// wildcard
	QueryWildcard bool
	// 查询sql
	RawJson string
}

func (e Elasticsearch) GetIndexName() string {
	indexName := e.Index
	if e.IndexOption.Index != "" {
		indexName = e.IndexOption.Index
		if e.IndexOption.WithDate {
			if e.IndexOption.Separator != "" {
				indexName += e.IndexOption.Separator
			}
			layout := e.IndexOption.ConvertDatePattern()
			indexName += time.Now().Format(layout)
		}
	}

	return indexName
}

type Logs struct {
	ProviderName string
	Metric       map[string]interface{}
	Message      []interface{}
}

func (l Logs) GetFingerprint() string {
	h := md5.New()
	streamString := tools.JsonMarshal(l.Metric)
	h.Write([]byte(streamString))
	fingerprint := hex.EncodeToString(h.Sum(nil))
	return fingerprint
}

func (l Logs) GetMetric() map[string]interface{} {
	return l.Metric
}

func (l Logs) GetAnnotations() []interface{} {
	return l.Message
}

func commonKeyValuePairs(maps []map[string]interface{}) map[string]interface{} {
	// 初始化一个map，用于记录每个key-value对的出现次数
	counts := make(map[string]int)

	// 获取map的数量
	mapCount := len(maps)

	// 遍历每个map并记录每个key-value对的出现次数
	for _, m := range maps {
		for k, v := range m {
			keyValue := fmt.Sprintf("%s:%v", k, v)
			counts[keyValue]++
		}
	}

	// 初始化结果map
	common := make(map[string]interface{})

	// 过滤只出现在所有map中的key-value对
	for keyValue, count := range counts {
		if count == mapCount {
			// 提取出key和value
			m := strings.SplitAfterN(keyValue, ":", 2)
			common[m[0]] = m[1]
		}
	}

	return common
}
