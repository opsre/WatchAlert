package models

type AlertDataSource struct {
	TenantId         string                 `json:"tenantId"`
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	Labels           map[string]interface{} `json:"labels" gorm:"labels;serializer:json"` // 额外标签，会添加到事件Metric中，可用于区分数据来源；
	Type             string                 `json:"type"`
	HTTP             HTTP                   `json:"http" gorm:"http;serializer:json"`
	Auth             Auth                   `json:"Auth" gorm:"auth;serializer:json"`
	DsAliCloudConfig DsAliCloudConfig       `json:"dsAliCloudConfig" gorm:"dsAliCloudConfig;serializer:json"`
	AWSCloudWatch    AWSCloudWatch          `json:"awsCloudwatch" gorm:"awsCloudwatch;serializer:json"`
	ClickHouseConfig DsClickHouseConfig     `json:"clickhouseConfig" gorm:"clickhouseConfig;serializer:json"`
	Description      string                 `json:"description"`
	KubeConfig       string                 `json:"kubeConfig"`
	Enabled          *bool                  `json:"enabled" `
}

type HTTP struct {
	URL     string `json:"url"`
	Timeout int64  `json:"timeout"`
}

type Auth struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type DsClickHouseConfig struct {
	Addr    string
	Timeout int64
}

type DsAliCloudConfig struct {
	AliCloudEndpoint   string `json:"alicloudEndpoint"`
	AliCloudAk         string `json:"alicloudAk"`
	AliCloudSk         string `json:"alicloudSk"`
	AliCloudSlsProject string `json:"alicloudSlsProject"`
}

type AWSCloudWatch struct {
	//Endpoint  string `json:"endpoint"`
	Region    string `json:"region"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

//type PromQueryRes struct {
//	Data data `json:"data"`
//}
//
//type data struct {
//	Result     []result `json:"result"`
//	ResultType string   `json:"resultType"`
//}
//
//type result struct {
//	Metric map[string]interface{} `json:"metric"`
//	Value  []interface{}          `json:"value"`
//}

func (d *AlertDataSource) GetEnabled() *bool {
	if d.Enabled == nil {
		isOk := false
		return &isOk
	}
	return d.Enabled
}
