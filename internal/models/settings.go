package models

import "watchAlert/config"

type Settings struct {
	IsInit          int                `json:"isInit"`
	AlarmConfig     config.AlarmConfig `json:"alarmConfig" gorm:"alarmConfig;serializer:json"`
	EmailConfig     emailConfig        `json:"emailConfig" gorm:"emailConfig;serializer:json"`
	AppVersion      string             `json:"appVersion" gorm:"-"`
	PhoneCallConfig phoneCallConfig    `json:"phoneCallConfig" gorm:"phoneCallConfig;serializer:json"`
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
