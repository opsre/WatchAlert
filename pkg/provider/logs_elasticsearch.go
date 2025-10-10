package provider

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/olivere/elastic/v7"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

type ElasticSearchDsProvider struct {
	cli            *elastic.Client
	url            string
	username       string
	password       string
	ExternalLabels map[string]interface{}
}

func NewElasticSearchClient(ctx context.Context, ds models.AlertDataSource) (LogsFactoryProvider, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(ds.HTTP.URL),
		elastic.SetBasicAuth(ds.Auth.User, ds.Auth.Pass),
		elastic.SetSniff(false),
	)
	if err != nil {
		return ElasticSearchDsProvider{}, err
	}

	return ElasticSearchDsProvider{
		cli:            client,
		url:            ds.HTTP.URL,
		username:       ds.Auth.User,
		password:       ds.Auth.Pass,
		ExternalLabels: ds.Labels,
	}, nil
}

type esQueryResponse struct {
	Source map[string]interface{} `json:"_source"`
}

func (e ElasticSearchDsProvider) Query(options LogQueryOptions) (Logs, int, error) {
	indexName := options.ElasticSearch.GetIndexName()
	var query elastic.Query

	switch options.ElasticSearch.QueryType {
	case models.EsQueryTypeRawJson:
		if options.ElasticSearch.RawJson == "" {
			return Logs{}, 0, errors.New("RawJson 为空")
		}
		query = elastic.NewRawStringQuery(options.ElasticSearch.RawJson)
	case models.EsQueryTypeField:
		conditionQuery := elastic.NewBoolQuery()
		if len(options.ElasticSearch.QueryFilter) > 0 {
			subQueries := make([]elastic.Query, 0, len(options.ElasticSearch.QueryFilter))
			for _, filter := range options.ElasticSearch.QueryFilter {
				var q elastic.Query
				switch options.ElasticSearch.QueryWildcard {
				case 0:
					// 精准匹配
					q = elastic.NewMatchQuery(filter.Field, filter.Value)
				case 1:
					// 模糊匹配
					q = elastic.NewWildcardQuery(filter.Field, fmt.Sprintf("*%v*", filter.Value))
				default:
					return Logs{}, 0, errors.New("undefined QueryWildcard")
				}
				subQueries = append(subQueries, q)
			}
			switch options.ElasticSearch.QueryFilterCondition {
			case models.EsFilterConditionOr:
				// 表示"或"关系，至少有一个子查询需要匹配
				conditionQuery = conditionQuery.Should(subQueries...).MinimumNumberShouldMatch(1)
			case models.EsFilterConditionAnd:
				// 表示"与"关系，所有子查询都必须匹配
				conditionQuery = conditionQuery.Must(subQueries...)
			case models.EsFilterConditionNot:
				// 表示"非"关系，所有子查询都不能匹配
				conditionQuery = conditionQuery.MustNot(subQueries...)
			default:
				return Logs{}, 0, errors.New("undefined QueryFilterCondition")
			}
		}
		conditionQuery.Must(elastic.NewRangeQuery("@timestamp").Gte(options.StartAt.(string)).Lte(options.EndAt.(string)))
		query = conditionQuery
	default:
		return Logs{}, 0, fmt.Errorf("undefined QueryType, type: %s", options.ElasticSearch.QueryType)
	}

	res, err := e.cli.Search().
		Index(indexName).
		Query(query).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		return Logs{}, 0, err
	}

	var response []esQueryResponse
	marshalHits, err := sonic.Marshal(res.Hits.Hits)
	if err != nil {
		return Logs{}, 0, err
	}
	err = sonic.Unmarshal(marshalHits, &response)
	if err != nil {
		return Logs{}, 0, err
	}

	var message []map[string]interface{}

	for _, v := range response {
		message = append(message, v.Source)
	}

	return Logs{
		ProviderName: ElasticSearchDsProviderName,
		Message:      message,
	}, len(response), nil
}

func (e ElasticSearchDsProvider) Check() (bool, error) {
	header := make(map[string]string)
	url := fmt.Sprintf("%s/_cat/health", e.url)
	if e.username != "" {
		auth := e.username + ":" + e.password
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
		header["Authorization"] = basicAuth
		url = fmt.Sprintf("%s/_cat/health", e.url)
	}
	res, err := tools.Get(header, url, 10)
	if err != nil {
		return false, err
	}

	if res.StatusCode != 200 {
		return false, fmt.Errorf("状态码非200, 当前: %d", res.StatusCode)
	}
	return true, nil
}

func (e ElasticSearchDsProvider) GetExternalLabels() map[string]interface{} {
	return e.ExternalLabels
}
