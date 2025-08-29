package services

import (
	"context"
	"fmt"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"

	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logc"
	"gopkg.in/ldap.v2"
)

type ldapService struct {
	ldapConfig models.LdapConfig
	ctx        *ctx.Context
}

type InterLdapService interface {
	ListUsers() ([]ldapUser, error)
	SyncUserToW8t()
	Login(username, password string) error
	SyncUsersCronjob(ctx context.Context)
}

func newInterLdapService(ctx *ctx.Context) InterLdapService {
	setting, err := ctx.DB.Setting().Get()
	if err != nil {
		return nil
	}

	return &ldapService{
		ctx:        ctx,
		ldapConfig: setting.LdapConfig,
	}
}

func (l ldapService) getAdminAuth() (*ldap.Conn, error) {
	ls, err := ldap.Dial("tcp", l.ldapConfig.Address)
	if err != nil {
		logc.Errorf(l.ctx.Ctx, fmt.Sprintf("无法连接 LDAP 服务器, Address: %s, err: %s", l.ldapConfig.Address, err.Error()))
		return nil, err
	}

	err = ls.Bind(l.ldapConfig.AdminUser, l.ldapConfig.AdminPass)
	if err != nil {
		logc.Errorf(l.ctx.Ctx, fmt.Sprintf("LDAP 管理员绑定失败 err: %s", err.Error()))
		return nil, err
	}

	return ls, nil
}

type ldapUser struct {
	Uid    string `json:"uid"`
	Mobile string `json:"mobile"`
	Mail   string `json:"mail"`
}

func (l ldapService) ListUsers() ([]ldapUser, error) {
	auth, err := l.getAdminAuth()
	if err != nil {
		return nil, err
	}
	defer auth.Close()

	logc.Infof(l.ctx.Ctx, "开始LDAP分页查询用户...")

	var totalResults []ldapUser
	pageSize := uint32(500)
	pages := 0
	pagingControl := ldap.NewControlPaging(pageSize)

	for {
		pages++

		// 创建搜索请求
		searchRequest := ldap.NewSearchRequest(
			l.ldapConfig.BaseDN,
			ldap.ScopeWholeSubtree,
			ldap.NeverDerefAliases,
			0, 0, false,
			"(objectClass=person)",
			[]string{"sAMAccountName", "cn", "mail", "mobile"},
			[]ldap.Control{pagingControl},
		)

		searchResult, err := auth.Search(searchRequest)
		if err != nil {
			logc.Errorf(l.ctx.Ctx, fmt.Sprintf("第 %d 页查询失败: %s", pages, err.Error()))
			return nil, err
		}

		pageUserCount := 0
		for _, entry := range searchResult.Entries {
			uid := entry.GetAttributeValue("sAMAccountName")
			if uid == "" {
				uid = entry.GetAttributeValue("cn")
			}
			if uid == "" {
				continue
			}

			totalResults = append(totalResults, ldapUser{
				Uid:    uid,
				Mobile: entry.GetAttributeValue("mobile"),
				Mail:   entry.GetAttributeValue("mail"),
			})
			pageUserCount++
		}

		var nextPageControl *ldap.ControlPaging
		for _, control := range searchResult.Controls {
			if control.GetControlType() == ldap.ControlTypePaging {
				nextPageControl = control.(*ldap.ControlPaging)
				break
			}
		}

		if nextPageControl == nil || len(nextPageControl.Cookie) == 0 {
			break
		}

		pagingControl = &ldap.ControlPaging{
			PagingSize: pageSize,
			Cookie:     nextPageControl.Cookie,
		}

		if pages >= 50 {
			break
		}
	}

	logc.Infof(l.ctx.Ctx, "LDAP分页查询完成，共 %d 页，获取 %d 个用户", pages, len(totalResults))
	return totalResults, nil
}

func (l ldapService) SyncUserToW8t() {
	users, err := l.ListUsers()
	if err != nil {
		logc.Errorf(l.ctx.Ctx, err.Error())
		return
	}

	for _, u := range users {
		_, b, _ := l.ctx.DB.User().Get("", "", u.Mail)
		if b {
			continue
		}
		uid := tools.RandUid()
		m := models.Member{
			UserId:   uid,
			UserName: u.Uid,
			Email:    u.Mail,
			Phone:    u.Mobile,
			CreateBy: "LDAP",
			CreateAt: time.Now().Unix(),
			Tenants:  []string{"default"},
		}
		err = l.ctx.DB.User().Create(m)
		if err != nil {
			logc.Errorf(l.ctx.Ctx, err.Error())
			return
		}

		err = l.ctx.DB.Tenant().AddTenantLinkedUsers("default",
			[]models.TenantUser{
				{
					UserID:   uid,
					UserName: u.Mail,
				},
			},
			l.ldapConfig.DefaultUserRole,
		)
		if err != nil {
			logc.Errorf(l.ctx.Ctx, err.Error())
			return
		}
	}
}

func (l ldapService) Login(username, password string) error {
	auth, err := l.getAdminAuth()
	if err != nil {
		logc.Errorf(l.ctx.Ctx, err.Error())
		return err
	}
	defer auth.Close()

	// 先搜索用户，获取真实的DN
	searchRequest := ldap.NewSearchRequest(
		l.ldapConfig.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		1, 0, false,
		fmt.Sprintf("(sAMAccountName=%s)", ldap.EscapeFilter(username)),
		[]string{"dn"},
		nil,
	)

	searchResult, err := auth.Search(searchRequest)
	if err != nil {
		logc.Errorf(l.ctx.Ctx, fmt.Sprintf("LDAP 搜索用户失败, username: %s, err: %s", username, err.Error()))
		return err
	}

	if len(searchResult.Entries) == 0 {
		logc.Errorf(l.ctx.Ctx, fmt.Sprintf("LDAP 用户不存在, username: %s", username))
		return fmt.Errorf("用户不存在")
	}

	userDN := searchResult.Entries[0].DN
	logc.Infof(l.ctx.Ctx, fmt.Sprintf("找到用户DN: %s", userDN))

	err = auth.Bind(userDN, password)
	if err != nil {
		logc.Errorf(l.ctx.Ctx, fmt.Sprintf("LDAP 用户登陆失败, username: %s, DN: %s, err: %s", username, userDN, err.Error()))
		return err
	}

	logc.Infof(l.ctx.Ctx, fmt.Sprintf("LDAP 用户登陆成功, username: %s", username))
	return nil
}


func (l ldapService) SyncUsersCronjob(ctx context.Context) {
	c := cron.New()
	_, err := c.AddFunc(l.ldapConfig.Cronjob, func() {
		l.SyncUserToW8t()
	})
	if err != nil {
		logc.Errorf(ctx, err.Error())
		return
	}
	c.Start()
	defer c.Stop()

	select {
	case <-ctx.Done():
		logc.Infof(ctx, "停止 LDAP 用户同步!")
		return
	}
}
