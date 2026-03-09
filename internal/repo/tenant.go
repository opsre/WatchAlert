package repo

import (
	"context"
	"watchAlert/internal/models"

	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
)

type (
	TenantRepo struct {
		entryRepo
	}

	InterTenantRepo interface {
		Create(t models.Tenant) error
		Update(t models.Tenant) error
		Delete(tenantId string) error
		List(userId string) (data []models.Tenant, err error)
		Get(tenantId string) (data models.Tenant, err error)
		CreateTenantLinkedUserRecord(t models.TenantLinkedUsers) error
		AddTenantLinkedUsers(tenantId string, users []models.TenantUser, userRole string) error
		RemoveTenantLinkedUsers(tenantId, userId string) error
		GetTenantLinkedUsers(tenantId string) (models.TenantLinkedUsers, error)
		DelTenantLinkedUserRecord(tenantId string) error
		GetTenantLinkedUserInfo(tenantId, userId string) (models.TenantUser, error)
		ChangeTenantUserRole(tenantId, userId, userRole string) error
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
		return err
	}

	var users = []models.TenantUser{
		{
			UserID:   "admin",
			UserName: "admin",
		},
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

		userData, _, err := tr.User().Get(u.UserID, "", "")
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

func (tr TenantRepo) Delete(tenantId string) error {
	getTenant, err := tr.Tenant().GetTenantLinkedUsers(tenantId)
	if err != nil {
		return err
	}

	for _, u := range getTenant.Users {
		err := tr.Tenant().RemoveTenantLinkedUsers(tenantId, u.UserID)
		if err != nil {
			return err
		}
	}

	err = tr.Tenant().DelTenantLinkedUserRecord(tenantId)
	if err != nil {
		return err
	}

	err = tr.g.Delete(Delete{
		Table: &models.Tenant{},
		Where: map[string]interface{}{
			"id = ?": tenantId,
		},
	})
	if err != nil {
		logc.Error(context.Background(), err)
		return err
	}
	return nil
}

func (tr TenantRepo) List(userId string) (data []models.Tenant, err error) {
	getUser, _, err := tr.User().Get(userId, "", "")
	if err != nil {
		return nil, err
	}

	var ts = &[]models.Tenant{}
	for _, tid := range getUser.Tenants {
		getT, err := tr.Tenant().Get(tid)
		if err != nil {
			return nil, err
		}
		*ts = append(*ts, getT)
	}

	return *ts, nil
}

func (tr TenantRepo) Get(tenantId string) (data models.Tenant, err error) {
	var d models.Tenant
	err = tr.db.Model(&models.Tenant{}).Where("id = ?", tenantId).First(&d).Error
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
func (tr TenantRepo) AddTenantLinkedUsers(tenantId string, users []models.TenantUser, userRole string) error {
	oldTenantUsers, err := tr.Tenant().GetTenantLinkedUsers(tenantId)
	if err != nil {
		return err
	}

	// 在新增成员时不会一并将角色写入，需要找到新增的成员，并且修改它的角色。
	for _, nUser := range users {
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
				UserRole: userRole,
			})
		}
	}

	// 更新租户表
	err = tr.g.Updates(Updates{
		Table: models.TenantLinkedUsers{},
		Where: map[string]interface{}{
			"id = ?": tenantId,
		},
		Updates: oldTenantUsers,
	})
	if err != nil {
		return err
	}

	// 更新用户表，新增租户ID
	for _, u := range users {
		userData, _, err := tr.User().Get(u.UserID, "", "")
		if err != nil {
			return err
		}

		var exist bool
		for _, tid := range userData.Tenants {
			if tid == tenantId {
				exist = true
			}
		}

		if !exist {
			userData.Tenants = append(userData.Tenants, tenantId)
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
func (tr TenantRepo) RemoveTenantLinkedUsers(tenantId, userId string) error {
	record, err := tr.GetTenantLinkedUsers(tenantId)
	if err != nil {
		return err
	}

	var newRecord []models.TenantUser
	// 移除租户中当前选择的用户，保留其他用户
	for _, u := range record.Users {
		if u.UserID == userId {
			continue
		}
		newRecord = append(newRecord, u)
	}
	record.Users = newRecord

	err = tr.g.Updates(Updates{
		Table: models.TenantLinkedUsers{},
		Where: map[string]interface{}{
			"id = ?": tenantId,
		},
		Updates: record,
	})
	if err != nil {
		return err
	}

	// 获取当前选择的用户详情
	userData, _, err := tr.User().Get(userId, "", "")
	if err != nil {
		return err
	}

	var newTenants = &[]string{}
	// 删除当前选择的租户，保留其他租户
	for _, tid := range userData.Tenants {
		if tid == tenantId {
			continue
		}
		*newTenants = append(*newTenants, tid)
	}

	userData.Tenants = *newTenants
	err = tr.g.Updates(Updates{
		Table: models.Member{},
		Where: map[string]interface{}{
			"user_id = ?": userId,
		},
		Updates: userData,
	})
	if err != nil {
		return err
	}

	return nil
}

// GetTenantLinkedUsers 获取租户关联的用户数据
func (tr TenantRepo) GetTenantLinkedUsers(tenantId string) (models.TenantLinkedUsers, error) {
	var d models.TenantLinkedUsers
	err := tr.db.Model(&models.TenantLinkedUsers{}).Where("id = ?", tenantId).First(&d).Error
	if err != nil {
		return d, err
	}

	return d, nil
}

// DelTenantLinkedUserRecord 删除租户关联表记录
func (tr TenantRepo) DelTenantLinkedUserRecord(tenantId string) error {
	err := tr.g.Delete(Delete{
		Table: &models.TenantLinkedUsers{},
		Where: map[string]interface{}{
			"id = ?": tenantId,
		},
	})
	if err != nil {
		logc.Error(context.Background(), err)
		return err
	}

	return nil
}

// GetTenantLinkedUserInfo 获取租户关联用户的详细信息
func (tr TenantRepo) GetTenantLinkedUserInfo(tenantId, userId string) (models.TenantUser, error) {
	var (
		tlu models.TenantLinkedUsers
		tu  models.TenantUser
	)

	err := tr.db.Model(&models.TenantLinkedUsers{}).Where("id = ?", tenantId).First(&tlu).Error
	if err != nil {
		return tu, err
	}

	for _, u := range tlu.Users {
		if u.UserID == userId {
			tu = u
			break
		}
	}

	return tu, nil
}

// ChangeTenantUserRole 修改用户角色
func (tr TenantRepo) ChangeTenantUserRole(tenantId, userId, userRole string) error {
	tenant, err := tr.GetTenantLinkedUsers(tenantId)
	if err != nil {
		return err
	}

	var users []models.TenantUser
	for _, u := range tenant.Users {
		if u.UserID != userId {
			users = append(users, u)
		} else {
			u.UserRole = userRole
			users = append(users, u)
		}
	}

	tenant.Users = users
	err = tr.g.Updates(Updates{
		Table: models.TenantLinkedUsers{},
		Where: map[string]interface{}{
			"id = ?": tenantId,
		},
		Updates: tenant,
	})
	if err != nil {
		return err
	}

	return nil
}
