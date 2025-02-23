package provider

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"watchAlert/internal/models"
	utilsHttp "watchAlert/pkg/tools"
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
		elastic.SetURL(ds.ElasticSearch.Url),
		elastic.SetBasicAuth(ds.ElasticSearch.Username, ds.ElasticSearch.Password),
		elastic.SetSniff(false),
	)
	if err != nil {
		return ElasticSearchDsProvider{}, err
	}

	return ElasticSearchDsProvider{
		cli:            client,
		url:            ds.ElasticSearch.Url,
		username:       ds.ElasticSearch.Username,
		password:       ds.ElasticSearch.Password,
		ExternalLabels: ds.Labels,
	}, nil
}

type esQueryResponse struct {
	Source map[string]interface{} `json:"_source"`
}

func (e ElasticSearchDsProvider) Query(options LogQueryOptions) ([]Logs, int, error) {
	indexName := options.ElasticSearch.GetIndexName()
	var query elastic.Query

	switch options.ElasticSearch.QueryType {
	case models.EsQueryTypeRawJson:
		if options.ElasticSearch.RawJson == "" {
			return nil, 0, errors.New("RawJson 为空")
		}
		query = elastic.NewRawStringQuery(options.ElasticSearch.RawJson)
	case models.EsQueryTypeCondition:
		conditionQuery := elastic.NewBoolQuery()
		if len(options.ElasticSearch.QueryFilter) > 0 {
			subQueries := make([]elastic.Query, 0, len(options.ElasticSearch.QueryFilter))
			for _, filter := range options.ElasticSearch.QueryFilter {
				var q elastic.Query
				if options.ElasticSearch.QueryWildcard {
					q = elastic.NewWildcardQuery(filter.Field, fmt.Sprintf("*%v*", filter.Value))
				} else {
					q = elastic.NewMatchQuery(filter.Field, filter.Value)
				}
				subQueries = append(subQueries, q)
			}
			switch options.ElasticSearch.QueryFilterCondition {
			case models.EsFilterConditionOr:
				conditionQuery = conditionQuery.Should(subQueries...).MinimumNumberShouldMatch(1)
			case models.EsFilterConditionAnd:
				conditionQuery = conditionQuery.Must(subQueries...)
			case models.EsFilterConditionNot:
				conditionQuery = conditionQuery.MustNot(subQueries...)
			}
		}
		conditionQuery.Must(elastic.NewRangeQuery("@timestamp").Gte(options.StartAt.(string)).Lte(options.EndAt.(string)))
		query = conditionQuery
	default:
		return nil, 0, errors.New("undefined QueryType")
	}

	res, err := e.cli.Search().
		Index(indexName).
		Query(query).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		return nil, 0, err
	}

	var response []esQueryResponse
	marshalHits, err := json.Marshal(res.Hits.Hits)
	if err != nil {
		return nil, 0, err
	}
	err = json.Unmarshal(marshalHits, &response)
	if err != nil {
		return nil, 0, err
	}

	var (
		data      []Logs
		msg       []interface{}
		kvMapList []map[string]interface{}
	)
	for _, v := range response {
		kvMapList = append(kvMapList, v.Source)
	}

	for _, m := range kvMapList {
		msg = append(msg, m["message"])
	}

	data = append(data, Logs{
		ProviderName: ElasticSearchDsProviderName,
		Metric:       commonKeyValuePairs(kvMapList),
		Message:      msg,
	})

	return data, len(msg), nil
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
	res, err := utilsHttp.Get(header, url, 10)
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
