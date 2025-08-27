package provider

import (
	"github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sls20201230 "github.com/alibabacloud-go/sls-20201230/v6/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"watchAlert/internal/models"
)

type AliCloudSlsDsProvider struct {
	client         *sls20201230.Client
	project        string
	ExternalLabels map[string]interface{}
}

func NewAliCloudSlsClient(source models.AlertDataSource) (LogsFactoryProvider, error) {
	config := &openapi.Config{
		AccessKeyId:     &source.DsAliCloudConfig.AliCloudAk,
		AccessKeySecret: &source.DsAliCloudConfig.AliCloudSk,
	}
	config.Endpoint = tea.String(source.DsAliCloudConfig.AliCloudEndpoint)
	result, err := sls20201230.NewClient(config)
	if err != nil {
		return AliCloudSlsDsProvider{}, err
	}

	return AliCloudSlsDsProvider{
		client:         result,
		project:        source.DsAliCloudConfig.AliCloudSlsProject,
		ExternalLabels: source.Labels,
	}, nil
}

func (a AliCloudSlsDsProvider) Query(query LogQueryOptions) (Logs, int, error) {
	var err error
	getLogsRequest := &sls20201230.GetLogsRequest{
		To:    tea.Int32(query.EndAt.(int32)),
		From:  tea.Int32(query.StartAt.(int32)),
		Query: tea.String(query.AliCloudSLS.Query),
	}
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	defer func() {
		if r := tea.Recover(recover()); r != nil {
			err = r
		}
	}()

	res, err := a.client.GetLogsWithOptions(tea.String(a.project), tea.String(query.AliCloudSLS.LogStore), getLogsRequest, headers, runtime)
	if err != nil {
		return Logs{}, 0, err
	}

	return Logs{
		ProviderName: AliCloudSLSDsProviderName,
		Message:      res.Body,
	}, len(res.Body), nil
}

func (a AliCloudSlsDsProvider) Check() (bool, error) {
	err := a.client.CheckConfig(&client.Config{})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (a AliCloudSlsDsProvider) GetExternalLabels() map[string]interface{} {
	return a.ExternalLabels
}
