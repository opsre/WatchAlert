package client

import (
	"context"
	"fmt"

	"watchAlert/internal/global"
	"watchAlert/internal/models"

	"github.com/glebarez/sqlite"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 全局数据库对象
var db *gorm.DB

func InitDB() *gorm.DB {
	// 初始化本地 test.db 数据库文件
	// db, err := gorm.Open(sqlite.Open("data/sql.db"), &gorm.Config{})
	switch global.Config.DataBase.Driver {
	case "mysql":
		db = ConnMysql()
	case "sqlite3":
		db = ConnSqlite()
	}
	dbAutoMigrate()
	return db
}

func dbAutoMigrate() {
	// 检查 Product 结构是否变化，变化则进行迁移
	_ = db.AutoMigrate(
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
		&models.ServiceResource{},
		&models.Tenant{},
		&models.Dashboard{},
		&models.AuditLog{},
		&models.Settings{},
		&models.TenantLinkedUsers{},
		&models.DashboardFolders{},
		&models.AlertSubscribe{},
		&models.NoticeRecord{},
		&models.ProbingRule{},
	)
}

func ConnSqlite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(global.Config.DataBase.Source), &gorm.Config{
		// 禁用外键(指定外键时不会在mysql创建真实的外键约束)
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logc.Errorf(context.Background(), "failed to connect database: %s", err.Error())
		return nil
	}
	dbObj, err := db.DB()
	if err != nil {
		logc.Error(context.Background(), err.Error())
		return nil
	}
	// 参见： https://github.com/glebarez/sqlite/issues/52
	dbObj.SetMaxOpenConns(1)
	return db
}

func ConnMysql() *gorm.DB {
	sql := global.Config.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4,utf8&parseTime=True&loc=Local&timeout=%s",
		sql.User,
		sql.Pass,
		sql.Host,
		sql.Port,
		sql.DBName,
		sql.Timeout)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logc.Errorf(context.Background(), "failed to connect database: %s", err.Error())
		return nil
	}

	if global.Config.Server.Mode == "debug" {
		db.Debug()
	} else {
		db.Logger = logger.Default.LogMode(logger.Silent)
	}

	return db
}
