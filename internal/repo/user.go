package repo

import (
	"context"
	"encoding/json"
	"fmt"
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
		Search(r models.MemberQuery) ([]models.Member, error)
		List() ([]models.Member, error)
		Get(r models.MemberQuery) (models.Member, bool, error)
		Create(r models.Member) error
		Update(r models.Member) error
		Delete(r models.MemberQuery) error
		ChangeCache(userId string)
		ChangePass(r models.Member) error
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

func (ur UserRepo) Search(r models.MemberQuery) ([]models.Member, error) {
	var data []models.Member
	var db = ur.db.Model(&models.Member{})
	if r.Query != "" {
		db.Where("user_name LIKE ? OR email Like ? OR phone LIKE ?", "%"+r.Query+"%", "%"+r.Query+"%", "%"+r.Query+"%")
	}
	if r.JoinDuty == "true" {
		db.Where("join_duty = ?", "true")
	}
	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (ur UserRepo) List() ([]models.Member, error) {
	var data []models.Member
	var db = ur.db.Model(&models.Member{})
	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (ur UserRepo) Get(r models.MemberQuery) (models.Member, bool, error) {
	var data models.Member
	db := ur.db.Model(models.Member{})
	if r.UserId != "" {
		db.Where("user_id = ?", r.UserId)
	}

	if r.UserName != "" {
		db.Where("user_name = ?", r.UserName)
	}

	if r.Query != "" {
		db.Where("user_id LIKE ? or user_name LIKE ? or email LIKE ? or phone LIKE ?", "%"+r.Query+"%", "%"+r.Query+"%", "%"+r.Query+"%", "%"+r.Query+"%")
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

func (ur UserRepo) Delete(r models.MemberQuery) error {
	userInfo, _, err := ur.User().Get(models.MemberQuery{UserId: r.UserId})
	if err != nil {
		return err
	}

	for _, tid := range userInfo.Tenants {
		err = ur.Tenant().RemoveTenantLinkedUsers(models.TenantQuery{ID: tid, UserID: r.UserId})
		if err != nil {
			return err
		}
	}

	d := Delete{
		Table: models.Member{},
		Where: map[string]interface{}{
			"user_id = ?": r.UserId,
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
	_ = json.Unmarshal([]byte(result), &cacheUser)

	duration, _ := client.Redis.TTL("uid-" + userId).Result()
	client.Redis.Set("uid-"+userId, tools.JsonMarshal(dbUser), duration)
}

func (ur UserRepo) ChangePass(r models.Member) error {
	u := Update{
		Table: models.Member{},
		Where: map[string]interface{}{
			"user_id = ?": r.UserId,
		},
		Update: []string{"password", r.Password},
	}

	err := ur.g.Update(u)
	if err != nil {
		return err
	}

	return nil
}
