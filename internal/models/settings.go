package models

type Settings struct {
	IsInit          int             `json:"isInit"`
	EmailConfig     emailConfig     `json:"emailConfig" gorm:"emailConfig;serializer:json"`
	AppVersion      string          `json:"appVersion" gorm:"-"`
	PhoneCallConfig phoneCallConfig `json:"phoneCallConfig" gorm:"phoneCallConfig;serializer:json"`
	AiConfig        AiConfig        `json:"aiConfig" gorm:"aiConfig;serializer:json"`
}

type emailConfig struct {
	ServerAddress string `json:"serverAddress"`
	Port          int    `json:"port"`
	Email         string `json:"email"`
	Token         string `json:"token"`
}

type phoneCallConfig struct {
	Provider        string `json:"provider"`
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	TtsCode         string `json:"ttsCode"`
}

// AiConfig ai config
type AiConfig struct {
	Enable    *bool  `json:"enable"`
	Type      string `json:"type"` // OpenAi, DeepSeek
	Url       string `json:"url"`
	AppKey    string `json:"appKey"`
	Model     string `json:"model"`
	Timeout   int    `json:"timeout"`
	MaxTokens int    `json:"maxTokens"`
	Prompt    string `json:"prompt"`
}

func (a AiConfig) GetEnable() bool {
	if a.Enable == nil {
		return false
	}

	return *a.Enable
}
