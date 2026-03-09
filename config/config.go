package config

import (
	"log"

	"github.com/spf13/viper"
)

type App struct {
	Server   Server   `json:"Server"`
	Database Database `json:"Database"`
	Redis    Redis    `json:"Redis"`
	Jwt      Jwt      `json:"Jwt"`
	Jaeger   Jaeger   `json:"Jaeger"`
}

type Server struct {
	Mode           string `json:"mode"`
	Port           string `json:"port"`
	EnableElection bool   `json:"enableElection"`
}

type Database struct {
	Type    string `json:"type"`    // mysql 或 sqlite
	Host    string `json:"host"`    // MySQL 主机地址
	Port    string `json:"port"`    // MySQL 端口
	User    string `json:"user"`    // MySQL 用户名
	Pass    string `json:"pass"`    // MySQL 密码
	DBName  string `json:"dbName"`  // MySQL 数据库名
	Timeout string `json:"timeout"` // MySQL 连接超时
	Path    string `json:"path"`    // SQLite 数据库文件路径
}

type Redis struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Pass     string `json:"pass"`
	Database int    `json:"database"`
}

type Jwt struct {
	Expire int64 `json:"expire"`
}

type Jaeger struct {
	URL string `json:"url"`
}

var (
	Application App
	Version     string
	configFile  = "config/config.yaml"
)

func InitConfig(version string) {
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

	Version = version
	Application = config
}
