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
	VictoriaDsProviderName      string = "VictoriaLogs"
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
	Victoria      Victoria
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
	// 索引名称
	Index string
	// 过滤条件
	QueryFilter []models.EsQueryFilter
	// filter关系，与或非
	QueryFilterCondition models.EsFilterCondition
	// 查询类型，sql语句查询与条件查询
	QueryType models.EsQueryType
	// wildcard
	QueryWildcard int64
	// 查询sql
	RawJson string
}

// Victoria victoriaMetrics数据源配置
type Victoria struct {
	Query string `json:"query"` // 查询语句
}

func (e Elasticsearch) GetIndexName() string {
	if strings.Contains(e.Index, "YYYY") && strings.Contains(e.Index, "MM") && strings.Contains(e.Index, "dd") {
		indexName := e.Index
		indexName = strings.ReplaceAll(indexName, "YYYY", time.Now().Format("2006"))
		indexName = strings.ReplaceAll(indexName, "MM", time.Now().Format("01"))
		indexName = strings.ReplaceAll(indexName, "dd", time.Now().Format("02"))
		return indexName
	}

	return e.Index
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
			m[0] = strings.ReplaceAll(m[0], ":", "")
			common[m[0]] = m[1]
		}
	}

	return common
}
