package services

import (
	"fmt"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
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
	Search(req interface{}) (interface{}, interface{})
	WithAddClientToProviderPools(datasource models.AlertDataSource) error
	WithRemoveClientForProviderPools(datasourceId string)
}

func newInterDatasourceService(ctx *ctx.Context) InterDatasourceService {
	return &datasourceService{
		ctx: ctx,
	}
}

func (ds datasourceService) Create(req interface{}) (interface{}, interface{}) {
	dataSource := req.(*models.AlertDataSource)

	id := "ds-" + tools.RandId()
	data := dataSource
	data.Id = id

	err := ds.ctx.DB.Datasource().Create(*dataSource)
	if err != nil {
		return nil, err
	}

	err = ds.WithAddClientToProviderPools(*dataSource)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ds datasourceService) Update(req interface{}) (interface{}, interface{}) {
	dataSource := req.(*models.AlertDataSource)

	err := ds.ctx.DB.Datasource().Update(*dataSource)
	if err != nil {
		return nil, err
	}

	err = ds.WithAddClientToProviderPools(*dataSource)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ds datasourceService) Delete(req interface{}) (interface{}, interface{}) {
	dataSource := req.(*models.DatasourceQuery)
	err := ds.ctx.DB.Datasource().Delete(*dataSource)
	if err != nil {
		return nil, err
	}

	ds.WithRemoveClientForProviderPools(dataSource.Id)

	return nil, nil
}

func (ds datasourceService) List(req interface{}) (interface{}, interface{}) {
	dataSource := req.(*models.DatasourceQuery)
	data, err := ds.ctx.DB.Datasource().List(*dataSource)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (ds datasourceService) Get(req interface{}) (interface{}, interface{}) {
	dataSource := req.(*models.DatasourceQuery)
	data, err := ds.ctx.DB.Datasource().Get(*dataSource)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (ds datasourceService) Search(req interface{}) (interface{}, interface{}) {
	var newData []models.AlertDataSource
	r := req.(*models.DatasourceQuery)
	data, err := ds.ctx.DB.Datasource().Search(*r)
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
	case provider.VictoriaDsProviderName:
		cli, err = provider.NewVictoriaLogsProvider(ctx.Ctx, datasource)
	case provider.JaegerDsProviderName:
		cli, err = provider.NewJaegerClient(datasource)
	case "Kubernetes":
		cli, err = provider.NewKubernetesClient(ds.ctx.Ctx, datasource.KubeConfig, datasource.Labels)
	case "CloudWatch":
		cli, err = provider.NewAWSCredentialCfg(datasource.AWSCloudWatch.Region, datasource.AWSCloudWatch.AccessKey, datasource.AWSCloudWatch.SecretKey, datasource.Labels)
	}

	if err != nil {
		return fmt.Errorf("New %s client failed, err: %s", datasource.Type, err.Error())
	}

	pools.SetClient(datasource.Id, cli)
	return nil
}

func (ds datasourceService) WithRemoveClientForProviderPools(datasourceId string) {
	pools := ds.ctx.Redis.ProviderPools()
	pools.RemoveClient(datasourceId)
}
