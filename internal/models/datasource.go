package models

type AlertDataSource struct {
	TenantId         string                 `json:"tenantId"`
	Id               string                 `json:"id"`
	Name             string                 `json:"name"`
	Labels           map[string]interface{} `json:"labels" gorm:"labels;serializer:json"` // 额外标签，会添加到事件Metric中，可用于区分数据来源；
	Type             string                 `json:"type"`
	HTTP             HTTP                   `json:"http" gorm:"http;serializer:json"`
	AliCloudEndpoint string                 `json:"alicloudEndpoint"`
	AliCloudAk       string                 `json:"alicloudAk"`
	AliCloudSk       string                 `json:"alicloudSk"`
	AWSCloudWatch    AWSCloudWatch          `json:"awsCloudwatch" gorm:"awsCloudwatch;serializer:json"`
	Description      string                 `json:"description"`
	KubeConfig       string                 `json:"kubeConfig"`
	ElasticSearch    ElasticSearch          `json:"elasticSearch" gorm:"elasticSearch;serializer:json"`
	Enabled          *bool                  `json:"enabled" `
}

type ElasticSearch struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type HTTP struct {
	URL     string `json:"url"`
	Timeout int64  `json:"timeout"`
}

type DatasourceQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	Id       string `json:"id" form:"id"`
	Type     string `json:"type" form:"type"`
	Query    string `json:"query" form:"query"`
}

type AWSCloudWatch struct {
	//Endpoint  string `json:"endpoint"`
	Region    string `json:"region"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type PromQueryReq struct {
	DatasourceType string `json:"datasourceType"`
	Addr           string `form:"addr"`
	Query          string `form:"query"`
}

type PromQueryRes struct {
	Data data `json:"data"`
}

type data struct {
	Result     []result `json:"result"`
	ResultType string   `json:"resultType"`
}

type result struct {
	Metric map[string]interface{} `json:"metric"`
	Value  []interface{}          `json:"value"`
}

func (d *AlertDataSource) GetEnabled() *bool {
	if d.Enabled == nil {
		isOk := false
		return &isOk
	}
	return d.Enabled
}
