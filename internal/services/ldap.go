package services

import (
	"fmt"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"

	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logc"
	"gopkg.in/ldap.v2"
)

type ldapService struct {
	ctx *ctx.Context
}

type InterLdapService interface {
	ListUsers() ([]ldapUser, error)
	SyncUserToW8t()
	Login(username, password string) error
	SyncUsersCronjob()
}

func newInterLdapService(ctx *ctx.Context) InterLdapService {
	return &ldapService{
		ctx: ctx,
	}
}

func (l ldapService) getAdminAuth() (*ldap.Conn, error) {
	ls, err := ldap.Dial("tcp", global.Config.Ldap.Address)
	if err != nil {
		logc.Errorf(l.ctx.Ctx, fmt.Sprintf("无法连接 LDAP 服务器, Address: %s, err: %s", global.Config.Ldap.Address, err.Error()))
		return nil, err
	}

	err = ls.Bind(global.Config.Ldap.AdminUser, global.Config.Ldap.AdminPass)
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
	lc := global.Config.Ldap

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
		logc.Infof(l.ctx.Ctx, "正在查询第 %d 页，页面大小: %d", pages, pageSize)

		// 创建搜索请求
		searchRequest := ldap.NewSearchRequest(
			lc.BaseDN,
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

		logc.Infof(l.ctx.Ctx, "第 %d 页完成，获取到 %d 个用户，总计: %d", pages, pageUserCount, len(totalResults))

		var nextPageControl *ldap.ControlPaging
		for _, control := range searchResult.Controls {
			if control.GetControlType() == ldap.ControlTypePaging {
				nextPageControl = control.(*ldap.ControlPaging)
				break
			}
		}

		if nextPageControl == nil || len(nextPageControl.Cookie) == 0 {
			logc.Infof(l.ctx.Ctx, "没有更多页面，查询完成")
			break
		}

		pagingControl = &ldap.ControlPaging{
			PagingSize: pageSize,
			Cookie:     nextPageControl.Cookie,
		}

		logc.Infof(l.ctx.Ctx, "找到下一页Cookie，长度: %d", len(nextPageControl.Cookie))

		if pages >= 50 {
			logc.Errorf(l.ctx.Ctx, "查询页数超过50页，停止查询")
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
		_, b, _ := l.ctx.DB.User().Get(models.MemberQuery{Query: u.Mail})
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

		err = l.ctx.DB.Tenant().AddTenantLinkedUsers(models.TenantLinkedUsers{
			ID:       "default",
			UserRole: global.Config.Ldap.DefaultUserRole,
			Users: []models.TenantUser{
				{
					UserID:   uid,
					UserName: u.Mail,
				},
			},
		})
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

	userDn := fmt.Sprintf("%s=%s,%s", global.Config.Ldap.UserPrefix, username, global.Config.Ldap.UserDN)
	err = auth.Bind(userDn, password)
	if err != nil {
		logc.Errorf(l.ctx.Ctx, fmt.Sprintf("LDAP 用户登陆失败, err: %s", err.Error()))
		return err
	}

	return nil
}

func (l ldapService) SyncUsersCronjob() {
	c := cron.New()
	_, err := c.AddFunc(global.Config.Ldap.Cronjob, func() {
		l.SyncUserToW8t()
	})
	if err != nil {
		logc.Errorf(ctx.Ctx, err.Error())
		return
	}
	c.Start()
	defer c.Stop()

	select {}
}
