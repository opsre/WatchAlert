package repo

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
	"watchAlert/internal/models"
	"watchAlert/pkg/client"
	"watchAlert/pkg/tools"
)

type (
	UserRepo struct {
		entryRepo
	}

	InterUserRepo interface {
		List(query, joinDuty string) ([]models.Member, error)
		Get(userId, username, query string) (models.Member, bool, error)
		Create(r models.Member) error
		Update(r models.Member) error
		Delete(userId string) error
		ChangeCache(userId string)
		ChangePass(userId, password string) error
	}
)

func newUserInterface(db *gorm.DB, g InterGormDBCli) InterUserRepo {
	return &UserRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (ur UserRepo) List(query, joinDuty string) ([]models.Member, error) {
	var (
		data []models.Member
		db   = ur.db.Model(&models.Member{})
	)

	if query != "" {
		db.Where("user_name LIKE ? OR email Like ? OR phone LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	}
	if joinDuty == "true" {
		db.Where("join_duty = ?", "true")
	}
	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (ur UserRepo) Get(userId, username, query string) (models.Member, bool, error) {
	var (
		data models.Member
		db   = ur.db.Model(models.Member{})
	)

	if userId != "" {
		db.Where("user_id = ?", userId)
	}
	if username != "" {
		db.Where("user_name = ?", username)
	}
	if query != "" {
		db.Where("user_id LIKE ? or user_name LIKE ? or email LIKE ? or phone LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	err := db.First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return data, false, fmt.Errorf("用户不存在")
		}
		return data, false, err
	}

	return data, true, nil
}

func (ur UserRepo) Create(r models.Member) error {
	err := ur.g.Create(models.Member{}, r)
	if err != nil {
		return err
	}

	if r.UserId == "admin" {
		r.Tenants = append(r.Tenants, "default")
		err = ur.g.Updates(Updates{
			Table: models.Member{},
			Where: map[string]interface{}{
				"user_id = ?": r.UserId,
			},
			Updates: r,
		})
	}

	return nil
}

func (ur UserRepo) Update(r models.Member) error {
	u := Updates{
		Table: models.Member{},
		Where: map[string]interface{}{
			"user_id = ?": r.UserId,
		},
		Updates: r,
	}

	err := ur.g.Updates(u)
	if err != nil {
		return err
	}

	return nil
}

func (ur UserRepo) Delete(userId string) error {
	userInfo, _, err := ur.User().Get(userId, "", "")
	if err != nil {
		return err
	}

	for _, tid := range userInfo.Tenants {
		err = ur.Tenant().RemoveTenantLinkedUsers(tid, userId)
		if err != nil {
			return err
		}
	}

	d := Delete{
		Table: models.Member{},
		Where: map[string]interface{}{
			"user_id = ?": userId,
		},
	}
	err = ur.g.Delete(d)
	if err != nil {
		return err
	}

	return nil
}

func (ur UserRepo) ChangeCache(userId string) {
	var dbUser models.Member
	ur.db.Model(&models.Member{}).Where("user_id = ?", userId).First(&dbUser)

	var cacheUser models.Member
	result, err := client.Redis.Get("uid-" + userId).Result()
	if err != nil {
		logc.Error(context.Background(), err)
	}
	_ = sonic.Unmarshal([]byte(result), &cacheUser)

	duration, _ := client.Redis.TTL("uid-" + userId).Result()
	client.Redis.Set("uid-"+userId, tools.JsonMarshalToString(dbUser), duration)
}

func (ur UserRepo) ChangePass(userId, password string) error {
	u := Update{
		Table: models.Member{},
		Where: map[string]interface{}{
			"user_id = ?": userId,
		},
		Update: []string{"password", password},
	}

	err := ur.g.Update(u)
	if err != nil {
		return err
	}

	return nil
}
