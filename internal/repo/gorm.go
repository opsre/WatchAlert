package repo

import (
	"fmt"
	"gorm.io/gorm"
)

type GormDBCli struct {
	db *gorm.DB
}

type InterGormDBCli interface {
	Create(table, value interface{}) error
	Update(value Update) error
	Updates(value Updates) error
	Delete(value Delete) error
}

func NewInterGormDBCli(db *gorm.DB) InterGormDBCli {
	return &GormDBCli{
		db: db,
	}
}

// Create 插入数据
func (g GormDBCli) Create(table, value interface{}) error {
	return g.executeTransaction(func(tx *gorm.DB) error {
		return tx.Model(table).Create(value).Error
	}, "数据写入失败")
}

// Update 更新单条数据
func (g GormDBCli) Update(value Update) error {
	return g.executeTransaction(func(tx *gorm.DB) error {
		tx = tx.Model(value.Table)
		for column, val := range value.Where {
			tx = tx.Where(column, val)
		}
		return tx.Update(value.Update[0], value.Update[1:]).Error
	}, "数据更新失败")
}

// Updates 更新多条数据
func (g GormDBCli) Updates(value Updates) error {
	return g.executeTransaction(func(tx *gorm.DB) error {
		tx = tx.Model(value.Table)
		for column, val := range value.Where {
			tx = tx.Where(column, val)
		}
		return tx.Updates(value.Updates).Error
	}, "数据更新失败")
}

// Delete 删除数据
func (g GormDBCli) Delete(value Delete) error {
	return g.executeTransaction(func(tx *gorm.DB) error {
		tx = tx.Model(value.Table)
		for column, val := range value.Where {
			tx = tx.Where(column, val)
		}
		return tx.Delete(value.Table).Error
	}, "数据删除失败")
}

// executeTransaction 执行事务并处理错误
func (g GormDBCli) executeTransaction(operation func(tx *gorm.DB) error, errorMessage string) error {
	tx := g.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("事务启动失败, err: %s", tx.Error)
	}

	if err := operation(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("%s -> %s", errorMessage, err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("事务提交失败, err: %s", err)
	}

	return nil
}

// Update 定义更新单条数据的结构
type Update struct {
	Table  interface{}
	Where  map[string]interface{}
	Update []string
}

// Updates 定义更新多条数据的结构
type Updates struct {
	Table   interface{}
	Where   map[string]interface{}
	Updates interface{}
}

// Delete 定义删除数据的结构
type Delete struct {
	Table interface{}
	Where map[string]interface{}
}
