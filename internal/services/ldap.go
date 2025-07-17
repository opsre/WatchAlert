package services

import (
	"fmt"
	"time"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
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

	var users []ldapUser
	pageSize := uint32(50)
	var pagingControl *ldap.ControlPaging

	logc.Infof(l.ctx.Ctx, "开始分页查询LDAP用户，每页大小: %d", pageSize)

	pageNum := 1
	for {
		logc.Infof(l.ctx.Ctx, "正在查询第 %d 页用户...", pageNum)

		searchRequest := ldap.NewSearchRequest(
			lc.BaseDN,
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, // size limit 设为 0，由分页控制
			"(&(objectClass=person)(objectClass=user))",
			[]string{"uid", "cn", "sAMAccountName", "mobile", "mail"},
			nil,
		)

		// 设置分页控制
		if pagingControl != nil {
			searchRequest.Controls = []ldap.Control{pagingControl}
		} else {
			pagingControl = ldap.NewControlPaging(pageSize)
			searchRequest.Controls = []ldap.Control{pagingControl}
		}

		sr, err := auth.Search(searchRequest)
		if err != nil {
			logc.Errorf(l.ctx.Ctx, fmt.Sprintf("LDAP 第 %d 页用户搜索失败, err: %s", pageNum, err.Error()))
			return nil, err
		}

		// 处理当前页的搜索结果
		pageUserCount := 0
		for _, entry := range sr.Entries {
			uid := entry.GetAttributeValue("sAMAccountName")
			if uid == "" {
				uid = entry.GetAttributeValue("uid")
			}
			if uid == "" {
				uid = entry.GetAttributeValue("cn")
			}
			if uid == "" {
				continue
			}
			users = append(users, ldapUser{
				Uid:    uid,
				Mobile: entry.GetAttributeValue("mobile"),
				Mail:   entry.GetAttributeValue("mail"),
			})
			pageUserCount++
		}

		logc.Infof(l.ctx.Ctx, "第 %d 页查询完成，获取到 %d 个用户", pageNum, pageUserCount)

		updatedControl := ldap.FindControl(sr.Controls, ldap.ControlTypePaging)
		if updatedControl == nil {
			logc.Infof(l.ctx.Ctx, "没有更多页面，查询结束")
			break
		}

		pagingControl = updatedControl.(*ldap.ControlPaging)
		if len(pagingControl.Cookie) == 0 {
			logc.Infof(l.ctx.Ctx, "分页Cookie为空，查询结束")
			break
		}

		pageNum++

		if pageNum > 100 {
			logc.Errorf(l.ctx.Ctx, "查询页数超过100页，可能存在问题，停止查询")
			break
		}
	}

	logc.Infof(l.ctx.Ctx, "LDAP 用户同步完成，共获取 %d 个用户", len(users))
	return users, nil
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
