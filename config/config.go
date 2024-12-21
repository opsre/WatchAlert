package config

import (
	"log"

	"github.com/spf13/viper"
)

type App struct {
	Server   Server   `json:"Server"`
	DataBase DataBase `json:"DataBase"`
	MySQL    MySQL    `json:"MySQL"`
	Redis    Redis    `json:"Redis"`
	Jwt      Jwt      `json:"Jwt"`
	Jaeger   Jaeger   `json:"Jaeger"`
	Ldap     Ldap     `json:"ldap"`
}

type Server struct {
	Mode        string `json:"mode"`
	Port        string `json:"port"`
	AlarmConfig AlarmConfig
}

type AlarmConfig struct {
	GroupWait     int `json:"groupWait"`
	GroupInterval int `json:"groupInterval"`
	RecoverWait   int `json:"recoverWait"`
}

type DataBase struct {
	Driver string `json:"driver"`
	Source string `json:"source"`
}

type MySQL struct {
	Host    string `json:"host"`
	Port    string `json:"port"`
	User    string `json:"user"`
	Pass    string `json:"pass"`
	DBName  string `json:"dbName"`
	Timeout string `json:"timeout"`
}

type Redis struct {
	Host string `json:"host"`
	Port string `json:"port"`
	Pass string `json:"pass"`
}

type Jwt struct {
	Expire int64 `json:"expire"`
}

type Jaeger struct {
	URL string `json:"url"`
}

type Ldap struct {
	Enabled         bool   `json:"enabled"`
	Address         string `json:"address"`
	BaseDN          string `json:"baseDN"`
	UserDN          string `json:"userDN"`
	AdminUser       string `json:"adminUser"`
	AdminPass       string `json:"adminPass"`
	UserPrefix      string `json:"userPrefix"`
	DefaultUserRole string `json:"defaultUserRole"`
	Cronjob         string `json:"cronjob"`
}

var configFile = "config/config.yaml"

func InitConfig() App {
	v := viper.New()
	v.SetConfigFile(configFile)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatal("配置读取失败:", err)
	}
	var config App
	if err := v.Unmarshal(&config); err != nil {
		log.Fatal("配置解析失败:", err)
	}
	return config
}
