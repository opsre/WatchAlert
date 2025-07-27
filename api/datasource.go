package api

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"strconv"
	"strings"
	"time"
	ctx2 "watchAlert/internal/ctx"
	"watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
	"watchAlert/pkg/provider"
	"watchAlert/pkg/tools"
)

type DatasourceController struct{}

/*
数据源 API
/api/w8t/datasource
*/
func (dc DatasourceController) API(gin *gin.RouterGroup) {
	datasourceA := gin.Group("datasource")
	datasourceA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		datasourceA.POST("dataSourceCreate", dc.Create)
		datasourceA.POST("dataSourceUpdate", dc.Update)
		datasourceA.POST("dataSourceDelete", dc.Delete)
	}

	datasourceB := gin.Group("datasource")
	datasourceB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		datasourceB.GET("dataSourceList", dc.List)
		datasourceB.GET("dataSourceGet", dc.Get)
		datasourceB.GET("dataSourceSearch", dc.Search)
	}

	c := gin.Group("datasource")
	c.Use(
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		c.GET("promQuery", dc.PromQuery)
		c.POST("dataSourcePing", dc.Ping)
		c.POST("searchViewLogsContent", dc.SearchViewLogsContent)
	}

}

func (dc DatasourceController) Create(ctx *gin.Context) {
	d := new(models.AlertDataSource)
	BindJson(ctx, d)

	tid, _ := ctx.Get("TenantID")
	d.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DatasourceService.Create(d)
	})
}

func (dc DatasourceController) List(ctx *gin.Context) {
	r := new(models.DatasourceQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DatasourceService.List(r)
	})
}

func (dc DatasourceController) Get(ctx *gin.Context) {
	r := new(models.DatasourceQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DatasourceService.Get(r)
	})
}

func (dc DatasourceController) Search(ctx *gin.Context) {
	r := new(models.DatasourceQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DatasourceService.Search(r)
	})
}

func (dc DatasourceController) Update(ctx *gin.Context) {
	r := new(models.AlertDataSource)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DatasourceService.Update(r)
	})
}

func (dc DatasourceController) Delete(ctx *gin.Context) {
	r := new(models.DatasourceQuery)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DatasourceService.Delete(r)
	})
}

func (dc DatasourceController) PromQuery(ctx *gin.Context) {
	r := new(models.PromQueryReq)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		var ress []provider.QueryResponse
		path := "/api/v1/query"
		params := url.Values{}
		params.Add("query", r.Query)
		params.Add("time", strconv.FormatInt(time.Now().Unix(), 10))

		var ids []string
		if len(r.DatasourceIds) < 2 {
			ids = []string{r.DatasourceIds}
		}

		ids = strings.Split(r.DatasourceIds, ",")
		for _, id := range ids {
			var res provider.QueryResponse
			source, err := ctx2.DO().DB.Datasource().Get(models.DatasourceQuery{
				Id: id,
			})
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

func (dc DatasourceController) Ping(ctx *gin.Context) {
	r := new(models.AlertDataSource)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		ok, err := provider.CheckDatasourceHealth(*r)
		if !ok {
			return "", fmt.Errorf("数据源不可达, err: %s", err.Error())
		}
		return "", nil
	})
}

// SearchViewLogsContent Logs 数据预览
func (dc DatasourceController) SearchViewLogsContent(ctx *gin.Context) {
	r := new(models.SearchLogsContentReq)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		data, err := services.DatasourceService.Get(&models.DatasourceQuery{
			Id: r.DatasourceId,
		})
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
