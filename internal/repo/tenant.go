package repo

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	TenantRepo struct {
		entryRepo
	}

	InterTenantRepo interface {
		Create(t models.Tenant) error
		Update(t models.Tenant) error
		Delete(t models.TenantQuery) error
		List(t models.TenantQuery) (data []models.Tenant, err error)
		Get(t models.TenantQuery) (data models.Tenant, err error)
		CreateTenantLinkedUserRecord(t models.TenantLinkedUsers) error
		AddTenantLinkedUsers(t models.TenantLinkedUsers) error
		RemoveTenantLinkedUsers(t models.TenantQuery) error
		GetTenantLinkedUsers(t models.TenantQuery) (models.TenantLinkedUsers, error)
		DelTenantLinkedUserRecord(t models.TenantQuery) error
		GetTenantLinkedUserInfo(t models.GetTenantLinkedUserInfo) (models.TenantUser, error)
		ChangeTenantUserRole(t models.ChangeTenantUserRole) error
	}
)

func newTenantInterface(db *gorm.DB, g InterGormDBCli) InterTenantRepo {
	return &TenantRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (tr TenantRepo) Create(t models.Tenant) error {
	err := tr.g.Create(&models.Tenant{}, t)
	if err != nil {
		logc.Error(context.Background(), err)
		return err
	}

	var users = []models.TenantUser{
		{
			UserID:   "admin",
			UserName: "admin",
		},
	}
	if t.UserId != "admin" {
		users = append(users, models.TenantUser{
			UserID:   t.UserId,
			UserName: t.CreateBy,
		})
	}

	for _, u := range users {
		err = tr.Tenant().CreateTenantLinkedUserRecord(
			models.TenantLinkedUsers{
				ID: t.ID,
				Users: []models.TenantUser{
					{
						UserID:   u.UserID,
						UserName: u.UserName,
						UserRole: "admin",
					},
				}})
		if err != nil {
			return err
		}

		userData, _, err := tr.User().Get(models.MemberQuery{UserId: u.UserID})
		if err != nil {
			return err
		}

		userData.Tenants = append(userData.Tenants, t.ID)
		err = tr.g.Updates(Updates{
			Table: models.Member{},
			Where: map[string]interface{}{
				"user_id = ?": u.UserID,
			},
			Updates: userData,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (tr TenantRepo) Update(t models.Tenant) error {
	u := Updates{
		Table: &models.Tenant{},
		Where: map[string]interface{}{
			"id = ?": t.ID,
		},
		Updates: t,
	}
	err := tr.g.Updates(u)
	if err != nil {
		logc.Error(context.Background(), err)
		return err
	}
	return nil
}

func (tr TenantRepo) Delete(t models.TenantQuery) error {
	getTenant, err := tr.Tenant().GetTenantLinkedUsers(models.TenantQuery{ID: t.ID})
	if err != nil {
		return err
	}

	for _, u := range getTenant.Users {
		err := tr.Tenant().RemoveTenantLinkedUsers(models.TenantQuery{ID: t.ID, UserID: u.UserID})
		if err != nil {
			return err
		}
	}

	err = tr.Tenant().DelTenantLinkedUserRecord(models.TenantQuery{ID: t.ID})
	if err != nil {
		return err
	}

	err = tr.g.Delete(Delete{
		Table: &models.Tenant{},
		Where: map[string]interface{}{
			"id = ?": t.ID,
		},
	})
	if err != nil {
		logc.Error(context.Background(), err)
		return err
	}
	return nil
}

func (tr TenantRepo) List(t models.TenantQuery) (data []models.Tenant, err error) {
	getUser, _, err := tr.User().Get(models.MemberQuery{UserId: t.UserID})
	if err != nil {
		return nil, err
	}

	var ts = &[]models.Tenant{}
	for _, tid := range getUser.Tenants {
		getT, err := tr.Tenant().Get(models.TenantQuery{ID: tid})
		if err != nil {
			return nil, err
		}
		*ts = append(*ts, getT)
	}

	return *ts, nil
}

func (tr TenantRepo) Get(t models.TenantQuery) (data models.Tenant, err error) {
	var d models.Tenant
	err = tr.db.Model(&models.Tenant{}).Where("id = ?", t.ID).First(&d).Error
	if err != nil {
		return d, err
	}
	return d, nil
}

// CreateTenantLinkedUserRecord 创建租户关联的用户记录
func (tr TenantRepo) CreateTenantLinkedUserRecord(t models.TenantLinkedUsers) error {
	err := tr.g.Create(&models.TenantLinkedUsers{}, t)
	if err != nil {
		logc.Error(context.Background(), err)
		return err
	}
	return nil
}

// AddTenantLinkedUsers 新增租户用户数据
func (tr TenantRepo) AddTenantLinkedUsers(t models.TenantLinkedUsers) error {
	oldTenantUsers, err := tr.Tenant().GetTenantLinkedUsers(models.TenantQuery{ID: t.ID})
	if err != nil {
		return err
	}

	// 在新增成员时不会一并将角色写入，需要找到新增的成员，并且修改它的角色。
	for _, nUser := range t.Users {
		found := false
		for _, oUser := range oldTenantUsers.Users {
			if oUser.UserID == nUser.UserID {
				found = true
				break
			}
		}
		if !found {
			oldTenantUsers.Users = append(oldTenantUsers.Users, models.TenantUser{
				UserID:   nUser.UserID,
				UserName: nUser.UserName,
				UserRole: t.UserRole,
			})
		}
	}

	err = tr.g.Updates(Updates{
		Table: models.TenantLinkedUsers{},
		Where: map[string]interface{}{
			"id = ?": t.ID,
		},
		Updates: oldTenantUsers,
	})
	if err != nil {
		return err
	}

	for _, u := range t.Users {
		userData, _, err := tr.User().Get(models.MemberQuery{UserId: u.UserID})
		if err != nil {
			return err
		}

		var exist bool
		for _, tid := range userData.Tenants {
			if tid == t.ID {
				exist = true
			}
		}

		if !exist {
			userData.Tenants = append(userData.Tenants, t.ID)
		}
		err = tr.g.Updates(Updates{
			Table: models.Member{},
			Where: map[string]interface{}{
				"user_id = ?": u.UserID,
			},
			Updates: userData,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// RemoveTenantLinkedUsers 移除租户关联的用户数据
func (tr TenantRepo) RemoveTenantLinkedUsers(t models.TenantQuery) error {
	record, err := tr.GetTenantLinkedUsers(models.TenantQuery{ID: t.ID})
	if err != nil {
		return err
	}

	var newRecord []models.TenantUser
	for _, u := range record.Users {
		if u.UserID == t.UserID {
			continue
		}
		newRecord = append(newRecord, u)
	}
	record.Users = newRecord

	err = tr.g.Updates(Updates{
		Table: models.TenantLinkedUsers{},
		Where: map[string]interface{}{
			"id = ?": t.ID,
		},
		Updates: record,
	})
	if err != nil {
		return err
	}

	userData, _, err := tr.User().Get(models.MemberQuery{UserId: t.UserID})
	if err != nil {
		return err
	}

	var newTenants = &[]string{}
	for _, tid := range userData.Tenants {
		if tid == t.ID {
			continue
		}
		*newTenants = append(*newTenants, tid)
	}

	userData.Tenants = *newTenants
	err = tr.g.Updates(Updates{
		Table: models.Member{},
		Where: map[string]interface{}{
			"user_id = ?": t.UserID,
		},
		Updates: userData,
	})
	if err != nil {
		return err
	}

	return nil
}

// GetTenantLinkedUsers 获取租户关联的用户数据
func (tr TenantRepo) GetTenantLinkedUsers(t models.TenantQuery) (models.TenantLinkedUsers, error) {
	var d models.TenantLinkedUsers
	err := tr.db.Model(&models.TenantLinkedUsers{}).Where("id = ?", t.ID).First(&d).Error
	if err != nil {
		return d, err
	}
	return d, nil
}

// DelTenantLinkedUserRecord 删除租户关联表记录
func (tr TenantRepo) DelTenantLinkedUserRecord(t models.TenantQuery) error {
	err := tr.g.Delete(Delete{
		Table: &models.TenantLinkedUsers{},
		Where: map[string]interface{}{
			"id = ?": t.ID,
		},
	})
	if err != nil {
		logc.Error(context.Background(), err)
		return err
	}

	return nil
}

// GetTenantLinkedUserInfo 获取租户关联用户的详细信息
func (tr TenantRepo) GetTenantLinkedUserInfo(t models.GetTenantLinkedUserInfo) (models.TenantUser, error) {
	var (
		tlu models.TenantLinkedUsers
		tu  models.TenantUser
	)
	err := tr.db.Model(&models.TenantLinkedUsers{}).Where("id = ?", t.ID).First(&tlu).Error
	if err != nil {
		return tu, err
	}

	for _, u := range tlu.Users {
		if u.UserID == t.UserID {
			tu = u
			break
		}
	}

	return tu, nil
}

// ChangeTenantUserRole 修改用户角色
func (tr TenantRepo) ChangeTenantUserRole(t models.ChangeTenantUserRole) error {
	tenant, err := tr.GetTenantLinkedUsers(models.TenantQuery{ID: t.ID})
	if err != nil {
		return err
	}

	fmt.Println(tenant)
	var users []models.TenantUser
	for _, u := range tenant.Users {
		if u.UserID != t.UserID {
			users = append(users, u)
		} else {
			u.UserRole = t.UserRole
			users = append(users, u)
		}
	}

	tenant.Users = users
	err = tr.g.Updates(Updates{
		Table: models.TenantLinkedUsers{},
		Where: map[string]interface{}{
			"id = ?": t.ID,
		},
		Updates: tenant,
	})
	if err != nil {
		return err
	}

	return nil
}
