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
	VictoriaLogsDsProviderName  string = "VictoriaLogs"
	ClickHouseDsProviderName    string = "ClickHouse"
)

type LogsFactoryProvider interface {
	Query(options LogQueryOptions) (Logs, int, error)
	Check() (bool, error)
	GetExternalLabels() map[string]interface{}
}

type LogQueryOptions struct {
	AliCloudSLS   AliCloudSLS
	Loki          Loki
	ElasticSearch Elasticsearch
	VictoriaLogs  VictoriaLogs
	ClickHouse    ClickHouse
	StartAt       interface{} // 查询的开始时间。
	EndAt         interface{} // 查询的结束时间。
}

type Loki struct {
	Query     string // 查询语句
	Direction string // 日志排序顺序，支持的值为forward或backward，默认为backward
	Limit     int64  // 要返回的最大条目数
}

type AliCloudSLS struct {
	Query    string   // 查询语句
	Project  string   // AliCloud SLS Project
	LogStore []string // AliCloud SLS LogStore
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

// VictoriaLogs victoriaMetrics数据源配置
type VictoriaLogs struct {
	Query string `json:"query"` // 查询语句
	Limit int    // 要返回的最大条目数
}

type ClickHouse struct {
	// 查询语句
	Query string
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
	Message      []map[string]interface{}
}

func (l Logs) GenerateFingerprint(ruleId string) string {
	h := md5.New()
	streamString := tools.JsonMarshalToString(map[string]string{
		"ruleId": ruleId,
	})
	h.Write([]byte(streamString))
	fingerprint := hex.EncodeToString(h.Sum(nil))
	return fingerprint
}

func (l Logs) GetAnnotations() map[string]interface{} {
	msg := make(map[string]interface{})
	if len(l.Message) == 0 {
		return msg
	}

	for k, v := range l.Message[0] {
		if v == nil {
			continue
		}

		switch value := v.(type) {
		case string:
			// 如果是字符串类型，处理长度限制
			content := value
			length := len(content)
			if length > 1000 {
				msg[k] = fmt.Sprintf("%s... 内容过长省略其中 ...%s", content[:500], content[length-500:])
			} else {
				msg[k] = content
			}
		case map[string]interface{}:
			// 如果是嵌套的 map，递归调用处理
			msg[k] = processNestedMap(value)
		default:
			// 对于其他类型，直接保留原值
			msg[k] = value
		}
	}
	return msg
}

// 辅助函数：递归处理嵌套的 map[string]interface{}
func processNestedMap(nestedMap map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range nestedMap {
		if v == nil {
			continue
		}

		switch value := v.(type) {
		case string:
			// 如果是字符串类型，处理长度限制
			content := value
			length := len(content)
			if length > 1000 {
				result[k] = fmt.Sprintf("%s... 内容过长省略其中 ...%s", content[:500], content[length-500:])
			} else {
				result[k] = content
			}
		case map[string]interface{}:
			// 如果是嵌套的 map，继续递归处理
			result[k] = processNestedMap(value)
		default:
			// 对于其他类型，直接保留原值
			result[k] = value
		}
	}
	return result
}
