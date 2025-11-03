package api

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
	ctx2 "watchAlert/internal/ctx"
	"watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
	"watchAlert/pkg/provider"
	"watchAlert/pkg/tools"
	jwtUtils "watchAlert/pkg/tools"

	"github.com/gin-gonic/gin"
)

type datasourceController struct{}

var DatasourceController = new(datasourceController)

/*
数据源 API
/api/w8t/datasource
*/
func (datasourceController datasourceController) API(gin *gin.RouterGroup) {
	a := gin.Group("datasource")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		a.POST("dataSourceCreate", datasourceController.Create)
		a.POST("dataSourceUpdate", datasourceController.Update)
		a.POST("dataSourceDelete", datasourceController.Delete)
	}

	b := gin.Group("datasource")
	b.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		b.GET("dataSourceList", datasourceController.List)
		b.GET("dataSourceGet", datasourceController.Get)
	}

	c := gin.Group("datasource")
	c.Use(
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		c.GET("promQuery", datasourceController.PromQuery)
		c.GET("promQueryRange", datasourceController.PromQueryRange)
		c.POST("dataSourcePing", datasourceController.Ping)
		c.POST("searchViewLogsContent", datasourceController.SearchViewLogsContent)
	}

}

func (datasourceController datasourceController) Create(ctx *gin.Context) {
	r := new(types.RequestDatasourceCreate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		userName := jwtUtils.GetUser(ctx.Request.Header.Get("Authorization"))
		r.UpdateBy = userName

		tid, _ := ctx.Get("TenantID")
		r.TenantId = tid.(string)

		return services.DatasourceService.Create(r)
	})
}

func (datasourceController datasourceController) List(ctx *gin.Context) {
	r := new(types.RequestDatasourceQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DatasourceService.List(r)
	})
}

func (datasourceController datasourceController) Get(ctx *gin.Context) {
	r := new(types.RequestDatasourceQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DatasourceService.Get(r)
	})
}

func (datasourceController datasourceController) Update(ctx *gin.Context) {
	r := new(types.RequestDatasourceUpdate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		userName := jwtUtils.GetUser(ctx.Request.Header.Get("Authorization"))
		r.UpdateBy = userName

		tid, _ := ctx.Get("TenantID")
		r.TenantId = tid.(string)

		return services.DatasourceService.Update(r)
	})
}

func (datasourceController datasourceController) Delete(ctx *gin.Context) {
	r := new(types.RequestDatasourceQuery)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DatasourceService.Delete(r)
	})
}

func (datasourceController datasourceController) PromQuery(ctx *gin.Context) {
	r := new(types.RequestQueryMetricsValue)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		var ress []provider.QueryResponse
		path := "/api/v1/query"
		params := url.Values{}
		params.Add("query", r.Query)
		params.Add("time", strconv.FormatInt(time.Now().Unix(), 10))

		var ids = []string{}
		ids = strings.Split(r.DatasourceIds, ",")
		for _, id := range ids {
			var res provider.QueryResponse
			source, err := ctx2.DO().DB.Datasource().Get(id)
			if err != nil {
				return nil, err
			}
			fullURL := fmt.Sprintf("%s%s?%s", source.HTTP.URL, path, params.Encode())
			get, err := tools.Get(tools.CreateBasicAuthHeader(source.Auth.User, source.Auth.Pass), fullURL, 10)
			if err != nil {
				return nil, err
			}

			if err := tools.ParseReaderBody(get.Body, &res); err != nil {
				return nil, err
			}

			ress = append(ress, res)
		}

		return ress, nil
	})
}

func (datasourceController datasourceController) PromQueryRange(ctx *gin.Context) {
	r := new(types.RequestQueryMetricsValue)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		err := r.Validate()
		if err != nil {
			return nil, err
		}

		var ress []provider.QueryResponse
		path := "/api/v1/query_range"
		params := url.Values{}
		params.Add("query", r.Query)
		params.Add("start", strconv.FormatInt(r.GetStartTime().Unix(), 10))
		params.Add("end", strconv.FormatInt(r.GetEndTime().Unix(), 10))
		params.Add("step", fmt.Sprintf("%.0fs", r.GetStep().Seconds()))

		var ids = []string{}
		ids = strings.Split(r.DatasourceIds, ",")

		for _, id := range ids {
			var res provider.QueryResponse
			source, err := ctx2.DO().DB.Datasource().Get(id)
			if err != nil {
				return nil, err
			}
			fullURL := fmt.Sprintf("%s%s?%s", source.HTTP.URL, path, params.Encode())
			get, err := tools.Get(tools.CreateBasicAuthHeader(source.Auth.User, source.Auth.Pass), fullURL, 10)
			if err != nil {
				return nil, err
			}

			if err := tools.ParseReaderBody(get.Body, &res); err != nil {
				return nil, err
			}

			ress = append(ress, res)
		}

		return ress, nil
	})
}

func (datasourceController datasourceController) Ping(ctx *gin.Context) {
	r := new(types.RequestDatasourceCreate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		ok, err := provider.CheckDatasourceHealth(models.AlertDataSource{
			TenantId:         r.TenantId,
			Name:             r.Name,
			Labels:           r.Labels,
			Type:             r.Type,
			HTTP:             r.HTTP,
			Auth:             r.Auth,
			DsAliCloudConfig: r.DsAliCloudConfig,
			AWSCloudWatch:    r.AWSCloudWatch,
			ClickHouseConfig: r.ClickHouseConfig,
			Description:      r.Description,
			KubeConfig:       r.KubeConfig,
			Enabled:          r.Enabled,
		})
		if !ok {
			return "", fmt.Errorf("数据源不可达, err: %s", err.Error())
		}
		return "", nil
	})
}

// SearchViewLogsContent Logs 数据预览
func (datasourceController datasourceController) SearchViewLogsContent(ctx *gin.Context) {
	r := new(types.RequestSearchLogsContent)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		data, err := services.DatasourceService.Get(&types.RequestDatasourceQuery{ID: r.DatasourceId})
		if err != nil {
			return nil, err
		}

		datasource := data.(models.AlertDataSource)

		var (
			client  provider.LogsFactoryProvider
			options provider.LogQueryOptions
		)

		// 使用 base64.StdEncoding 进行解码
		decodedBytes, err := base64.StdEncoding.DecodeString(r.Query)
		if err != nil {
			return nil, fmt.Errorf("base64 解码失败: %s", err)
		}
		// 将解码后的字节转换为字符串
		QueryStr := string(decodedBytes)

		switch r.Type {
		case provider.VictoriaLogsDsProviderName:
			client, err = provider.NewVictoriaLogsClient(ctx, datasource)
			if err != nil {
				return nil, err
			}

			options = provider.LogQueryOptions{
				VictoriaLogs: provider.VictoriaLogs{
					Query: QueryStr,
				},
			}
		case provider.ElasticSearchDsProviderName:
			client, err = provider.NewElasticSearchClient(ctx, datasource)
			if err != nil {
				return nil, err
			}

			options = provider.LogQueryOptions{
				ElasticSearch: provider.Elasticsearch{
					Index:     r.GetElasticSearchIndexName(),
					QueryType: "RawJson",
					RawJson:   QueryStr,
				},
			}
		case provider.ClickHouseDsProviderName:
			client, err = provider.NewClickHouseClient(ctx, datasource)
			if err != nil {
				return nil, err
			}

			options = provider.LogQueryOptions{
				ClickHouse: provider.ClickHouse{
					Query: QueryStr,
				},
			}
		}

		query, _, err := client.Query(options)
		if err != nil {
			return nil, err
		}

		return query, nil
	})
}
