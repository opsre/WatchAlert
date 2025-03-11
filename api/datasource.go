package api

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"strconv"
	"strings"
	"time"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
	ctx2 "watchAlert/pkg/ctx"
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
		datasourceB.GET("promQuery", dc.PromQuery)
		datasourceB.POST("dataSourcePing", dc.Ping)
		datasourceB.POST("esSearch", dc.EsSearch)
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

// EsSearch es 内容搜索
func (dc DatasourceController) EsSearch(ctx *gin.Context) {
	r := new(models.EsSearchReq)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		data, err := services.DatasourceService.Get(&models.DatasourceQuery{
			Id: r.DatasourceId,
		})
		if err != nil {
			return nil, err
		}
		datasource := data.(models.AlertDataSource)
		client, err := provider.NewElasticSearchClient(ctx, datasource)
		if err != nil {
			return nil, err
		}

		// 使用 base64.StdEncoding 进行解码
		decodedBytes, err := base64.StdEncoding.DecodeString(r.Query)
		if err != nil {
			return nil, fmt.Errorf("base64 解码失败: %s", err)
		}

		// 将解码后的字节转换为字符串
		QueryStr := string(decodedBytes)

		query, _, err := client.Query(provider.LogQueryOptions{ElasticSearch: provider.Elasticsearch{
			Index:     r.GetIndexName(),
			QueryType: "RawJson",
			RawJson:   QueryStr,
		}})
		if err != nil {
			return nil, err
		}

		type newLogStruct struct {
			Metric  map[string]interface{}
			Message interface{}
		}

		var newData []newLogStruct
		for _, v := range query {
			for _, message := range v.Message {
				newData = append(newData, newLogStruct{
					Metric:  v.Metric,
					Message: message,
				})
			}
		}

		return tools.JsonMarshal(newData), nil
	})
}
