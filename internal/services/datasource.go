package services

import (
	"fmt"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/provider"
	"watchAlert/pkg/tools"
)

type datasourceService struct {
	ctx *ctx.Context
}

type InterDatasourceService interface {
	Create(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Delete(req interface{}) (interface{}, interface{})
	List(req interface{}) (interface{}, interface{})
	Get(req interface{}) (interface{}, interface{})
	WithAddClientToProviderPools(datasource models.AlertDataSource) error
	WithRemoveClientForProviderPools(datasourceId string)
}

func newInterDatasourceService(ctx *ctx.Context) InterDatasourceService {
	return &datasourceService{
		ctx: ctx,
	}
}

func (ds datasourceService) Create(req interface{}) (interface{}, interface{}) {
	dataSource := req.(*types.RequestDatasourceCreate)

	data := models.AlertDataSource{
		TenantId:         dataSource.TenantId,
		ID:               "ds-" + tools.RandId(),
		Name:             dataSource.Name,
		Labels:           dataSource.Labels,
		Type:             dataSource.Type,
		HTTP:             dataSource.HTTP,
		Auth:             dataSource.Auth,
		DsAliCloudConfig: dataSource.DsAliCloudConfig,
		AWSCloudWatch:    dataSource.AWSCloudWatch,
		ClickHouseConfig: dataSource.ClickHouseConfig,
		Description:      dataSource.Description,
		KubeConfig:       dataSource.KubeConfig,
		UpdateBy:         dataSource.UpdateBy,
		UpdateAt:         time.Now().Unix(),
		Enabled:          dataSource.Enabled,
	}

	err := ds.ctx.DB.Datasource().Create(data)
	if err != nil {
		return nil, err
	}

	err = ds.WithAddClientToProviderPools(data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ds datasourceService) Update(req interface{}) (interface{}, interface{}) {
	dataSource := req.(*types.RequestDatasourceUpdate)

	data := models.AlertDataSource{
		TenantId:         dataSource.TenantId,
		ID:               dataSource.ID,
		Name:             dataSource.Name,
		Labels:           dataSource.Labels,
		Type:             dataSource.Type,
		HTTP:             dataSource.HTTP,
		Auth:             dataSource.Auth,
		DsAliCloudConfig: dataSource.DsAliCloudConfig,
		AWSCloudWatch:    dataSource.AWSCloudWatch,
		ClickHouseConfig: dataSource.ClickHouseConfig,
		Description:      dataSource.Description,
		KubeConfig:       dataSource.KubeConfig,
		UpdateBy:         dataSource.UpdateBy,
		UpdateAt:         time.Now().Unix(),
		Enabled:          dataSource.Enabled,
	}

	err := ds.ctx.DB.Datasource().Update(data)
	if err != nil {
		return nil, err
	}

	err = ds.WithAddClientToProviderPools(data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ds datasourceService) Delete(req interface{}) (interface{}, interface{}) {
	dataSource := req.(*types.RequestDatasourceQuery)
	err := ds.ctx.DB.Datasource().Delete(dataSource.TenantId, dataSource.ID)
	if err != nil {
		return nil, err
	}

	ds.WithRemoveClientForProviderPools(dataSource.ID)

	return nil, nil
}

func (ds datasourceService) Get(req interface{}) (interface{}, interface{}) {
	dataSource := req.(*types.RequestDatasourceQuery)
	data, err := ds.ctx.DB.Datasource().Get(dataSource.ID)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (ds datasourceService) List(req interface{}) (interface{}, interface{}) {
	var newData []models.AlertDataSource
	dataSource := req.(*types.RequestDatasourceQuery)
	data, err := ds.ctx.DB.Datasource().List(dataSource.TenantId, dataSource.ID, dataSource.Type, dataSource.Query)
	if err != nil {
		return nil, err
	}
	newData = data

	return newData, nil
}

func (ds datasourceService) WithAddClientToProviderPools(datasource models.AlertDataSource) error {
	var (
		cli interface{}
		err error
	)
	pools := ds.ctx.Redis.ProviderPools()
	switch datasource.Type {
	case provider.PrometheusDsProvider:
		cli, err = provider.NewPrometheusClient(datasource)
	case provider.VictoriaMetricsDsProvider:
		cli, err = provider.NewVictoriaMetricsClient(datasource)
	case provider.LokiDsProviderName:
		cli, err = provider.NewLokiClient(datasource)
	case provider.AliCloudSLSDsProviderName:
		cli, err = provider.NewAliCloudSlsClient(datasource)
	case provider.ElasticSearchDsProviderName:
		cli, err = provider.NewElasticSearchClient(ctx.Ctx, datasource)
	case provider.VictoriaLogsDsProviderName:
		cli, err = provider.NewVictoriaLogsClient(ctx.Ctx, datasource)
	case provider.JaegerDsProviderName:
		cli, err = provider.NewJaegerClient(datasource)
	case "Kubernetes":
		cli, err = provider.NewKubernetesClient(ds.ctx.Ctx, datasource.KubeConfig, datasource.Labels)
	case "CloudWatch":
		cli, err = provider.NewAWSCredentialCfg(datasource.AWSCloudWatch.Region, datasource.AWSCloudWatch.AccessKey, datasource.AWSCloudWatch.SecretKey, datasource.Labels)
	case "ClickHouse":
		cli, err = provider.NewClickHouseClient(ctx.Ctx, datasource)
	}

	if err != nil {
		return fmt.Errorf("New %s client failed, err: %s", datasource.Type, err.Error())
	}

	pools.SetClient(datasource.ID, cli)
	return nil
}

func (ds datasourceService) WithRemoveClientForProviderPools(datasourceId string) {
	pools := ds.ctx.Redis.ProviderPools()
	pools.RemoveClient(datasourceId)
}
