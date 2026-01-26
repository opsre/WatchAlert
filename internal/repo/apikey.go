package repo

import (
	"fmt"
	"watchAlert/internal/models"

	"gorm.io/gorm"
)

type (
	ApiKeyRepo struct {
		entryRepo
	}

	InterApiKeyRepo interface {
		List(userId string) ([]models.ApiKey, error)
		Get(id int, userId string) (models.ApiKey, bool, error)
		Create(key models.ApiKey) error
		Update(key models.ApiKey) error
		Delete(id int, userId string) error
		GetByKey(key string) (models.ApiKey, bool, error)
	}
)

func newApiKeyInterface(db *gorm.DB, g InterGormDBCli) InterApiKeyRepo {
	return &ApiKeyRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (ar ApiKeyRepo) List(userId string) ([]models.ApiKey, error) {
	var data []models.ApiKey
	db := ar.db.Model(&models.ApiKey{})

	if userId != "" {
		db = db.Where("user_id = ?", userId)
	}

	err := db.Order("created_at DESC").Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (ar ApiKeyRepo) Get(id int, userId string) (models.ApiKey, bool, error) {
	var data models.ApiKey
	db := ar.db.Model(&models.ApiKey{})

	db = db.Where("id = ?", id)
	if userId != "" {
		db = db.Where("user_id = ?", userId)
	}

	err := db.First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return data, false, fmt.Errorf("API密钥不存在")
		}
		return data, false, err
	}

	return data, true, nil
}

func (ar ApiKeyRepo) GetByKey(key string) (models.ApiKey, bool, error) {
	var data models.ApiKey
	db := ar.db.Model(&models.ApiKey{})

	db = db.Where("`key` = ?", key)

	err := db.First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return data, false, fmt.Errorf("API密钥不存在")
		}
		return data, false, err
	}

	return data, true, nil
}

func (ar ApiKeyRepo) Create(key models.ApiKey) error {
	err := ar.g.Create(&models.ApiKey{}, &key)
	if err != nil {
		return err
	}

	return nil
}

func (ar ApiKeyRepo) Update(key models.ApiKey) error {
	u := Updates{
		Table: &models.ApiKey{},
		Where: map[string]interface{}{
			"id = ?": key.ID,
		},
		Updates: key,
	}

	err := ar.g.Updates(u)
	if err != nil {
		return err
	}

	return nil
}

func (ar ApiKeyRepo) Delete(id int, userId string) error {
	key, _, err := ar.Get(id, userId)
	if err != nil {
		return err
	}

	db := ar.db.Model(&models.ApiKey{})
	db = db.Where("id = ?", id)
	if userId != "" {
		db = db.Where("user_id = ?", userId)
	}

	err = db.Delete(&key).Error
	if err != nil {
		return err
	}

	return nil
}
