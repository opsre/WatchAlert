package client

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"watchAlert/internal/global"
	"watchAlert/internal/models"

	"github.com/glebarez/sqlite"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Type    string // 数据库类型: mysql 或 sqlite
	Host    string // MySQL 主机地址
	Port    string // MySQL 端口
	User    string // MySQL 用户名
	Pass    string // MySQL 密码
	DBName  string // MySQL 数据库名
	Timeout string // MySQL 连接超时
	Path    string // SQLite 数据库文件路径
}

func NewDBClient(config DBConfig) *gorm.DB {
	var db *gorm.DB
	var err error

	// 设置默认数据库类型为 mysql
	if config.Type == "" {
		config.Type = "mysql"
	}

	switch config.Type {
	case "sqlite":
		db, err = initSQLiteDB(config)
	case "mysql":
		db, err = initMySQLDB(config)
	default:
		logc.Errorf(context.Background(), "unsupported database type: %s", config.Type)
		return nil
	}

	if err != nil {
		logc.Errorf(context.Background(), "failed to connect database: %s", err.Error())
		return nil
	}

	// 检查 Product 结构是否变化，变化则进行迁移
	err = db.AutoMigrate(
		&models.DutySchedule{},
		&models.DutyManagement{},
		&models.AlertNotice{},
		&models.AlertDataSource{},
		&models.AlertRule{},
		&models.AlertCurEvent{},
		&models.AlertHisEvent{},
		&models.AlertSilences{},
		&models.Member{},
		&models.UserRole{},
		&models.UserPermissions{},
		&models.NoticeTemplateExample{},
		&models.RuleGroups{},
		&models.RuleTemplateGroup{},
		&models.RuleTemplate{},
		&models.Tenant{},
		&models.Dashboard{},
		&models.AuditLog{},
		&models.Settings{},
		&models.TenantLinkedUsers{},
		&models.DashboardFolders{},
		&models.AlertSubscribe{},
		&models.NoticeRecord{},
		&models.ProbeRule{},
		&models.FaultCenter{},
		&models.AiContentRecord{},
		&models.Comment{},
		&models.Topology{},
		&models.ApiKey{},
	)
	if err != nil {
		logc.Error(context.Background(), err.Error())
		return nil
	}

	if global.Config.Server.Mode == "debug" {
		db.Debug()
	} else {
		db.Logger = logger.Default.LogMode(logger.Silent)
	}

	return db
}

// initMySQLDB 初始化 MySQL 数据库连接
func initMySQLDB(config DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4,utf8&parseTime=True&loc=Local&timeout=%s",
		config.User,
		config.Pass,
		config.Host,
		config.Port,
		config.DBName,
		config.Timeout)

	logc.Infof(context.Background(), "connecting to MySQL database: %s:%s/%s", config.Host, config.Port, config.DBName)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

// initSQLiteDB 初始化 SQLite 数据库连接
func initSQLiteDB(config DBConfig) (*gorm.DB, error) {
	// 设置默认 SQLite 文件路径
	if config.Path == "" {
		config.Path = "data/watchalert.db"
	}

	// 确保目录存在
	dir := filepath.Dir(config.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	logc.Infof(context.Background(), "connecting to SQLite database: %s", config.Path)
	return gorm.Open(sqlite.Open(config.Path), &gorm.Config{})
}
