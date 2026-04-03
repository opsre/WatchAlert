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
	ctx *ctx.Context
}

type InterLdapService interface {
	ListUsers(ldapConfig models.LdapConfig) ([]ldapUser, error)
	SyncUserToW8t(ldapConfig models.LdapConfig)
	Login(username, password string) error
	SyncUsersCronjob(ctx context.Context, ldapConfig models.LdapConfig)
	SyncNow() error
}

func newInterLdapService(ctx *ctx.Context) InterLdapService {
	return &ldapService{
		ctx: ctx,
	}
}

func (l ldapService) getAdminAuth(ldapConfig models.LdapConfig) (*ldap.Conn, error) {
	ls, err := ldap.Dial("tcp", ldapConfig.Address)
	if err != nil {
		logc.Errorf(l.ctx.Ctx, fmt.Sprintf("无法连接 LDAP 服务器, Address: %s, err: %s", ldapConfig.Address, err.Error()))
		return nil, err
	}

	err = ls.Bind(ldapConfig.AdminUser, ldapConfig.AdminPass)
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

func (l ldapService) ListUsers(ldapConfig models.LdapConfig) ([]ldapUser, error) {
	logc.Infof(l.ctx.Ctx, "开始 LDAP 的分页查询用户...")

	auth, err := l.getAdminAuth(ldapConfig)
	if err != nil {
		return nil, err
	}
	defer auth.Close()

	var totalResults []ldapUser
	pageSize := uint32(500)
	pages := 0
	pagingControl := ldap.NewControlPaging(pageSize)

	// 构建查询过滤器
	listFilter := "(objectClass=person)"
	if ldapConfig.Filter != "" {
		listFilter = ldapConfig.Filter
	}

	// 属性列表
	attributes := []string{"sAMAccountName", "cn", "mail", "mobile"}

	for {
		pages++

		// 创建搜索请求
		searchRequest := ldap.NewSearchRequest(
			ldapConfig.BaseDN,
			ldap.ScopeWholeSubtree,
			ldap.NeverDerefAliases,
			0, 0, false,
			listFilter,
			attributes,
			[]ldap.Control{pagingControl},
		)

		searchResult, err := auth.Search(searchRequest)
		if err != nil {
			logc.Error(l.ctx.Ctx, fmt.Sprintf("第 %d 页查询失败: %s", pages, err.Error()))
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

	logc.Infof(l.ctx.Ctx, "完成 LDAP 分页查询，共 %d 页，获取 %d 个用户", pages, len(totalResults))
	return totalResults, nil
}

func (l ldapService) SyncUserToW8t(ldapConfig models.LdapConfig) {
	users, err := l.ListUsers(ldapConfig)
	if err != nil {
		logc.Error(l.ctx.Ctx, err.Error())
		return
	}

	for _, u := range users {
		_, b, _ := l.ctx.DB.User().Get("", u.Uid, u.Mail, u.Mobile)
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
			Tenants:  []string{ldapConfig.DefaultTenant},
		}
		err = l.ctx.DB.User().Create(m)
		if err != nil {
			logc.Error(l.ctx.Ctx, err.Error())
			return
		}

		err = l.ctx.DB.Tenant().AddTenantLinkedUsers(ldapConfig.DefaultTenant,
			[]models.TenantUser{
				{
					UserID:   uid,
					UserName: u.Mail,
				},
			},
			ldapConfig.DefaultUserRole,
		)
		if err != nil {
			logc.Error(l.ctx.Ctx, err.Error())
			return
		}
	}
}

func (l ldapService) Login(username, password string) error {
	if username == "" || password == "" {
		return fmt.Errorf("LDAP 用户名或密码不能为空")
	}

	settings, err := l.ctx.DB.Setting().Get()
	if err != nil {
		logc.Error(l.ctx.Ctx, "获取 LDAP 配置失败: %s", err.Error())
		return fmt.Errorf("获取 LDAP 配置失败: %s", err.Error())
	}

	auth, err := l.getAdminAuth(settings.LdapConfig)
	if err != nil {
		logc.Error(l.ctx.Ctx, "LDAP 连接失败: %s", err.Error())
		return fmt.Errorf("LDAP 连接失败: %s", err.Error())
	}
	defer auth.Close()

	// 先搜索用户，获取真实的DN
	loginFilter := "(objectClass=person)"
	if settings.LdapConfig.Filter != "" {
		loginFilter = settings.LdapConfig.Filter
	}

	searchRequest := ldap.NewSearchRequest(
		settings.LdapConfig.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		1, 0, false,
		loginFilter,
		[]string{"dn"},
		nil,
	)

	searchResult, err := auth.Search(searchRequest)
	if err != nil {
		logc.Error(l.ctx.Ctx, fmt.Sprintf("LDAP 搜索用户失败, username: %s, err: %s", username, err.Error()))
		return fmt.Errorf("LDAP 搜索用户失败: %s", err.Error())
	}

	if len(searchResult.Entries) == 0 {
		logc.Error(l.ctx.Ctx, fmt.Sprintf("LDAP 用户不存在, username: %s", username))
		return fmt.Errorf("LDAP 用户不存在")
	}

	userDN := searchResult.Entries[0].DN
	err = auth.Bind(userDN, password)
	if err != nil {
		logc.Error(l.ctx.Ctx, fmt.Sprintf("LDAP 用户登陆失败, username: %s, DN: %s, err: %s", username, userDN, err.Error()))
		return fmt.Errorf("LDAP 用户登陆失败: %s", err.Error())
	}

	logc.Info(l.ctx.Ctx, fmt.Sprintf("LDAP 用户登陆成功, username: %s, DN: %s", username, userDN))
	return nil
}

func (l ldapService) SyncUsersCronjob(ctx context.Context, ldapConfig models.LdapConfig) {
	if ldapConfig.Cronjob == "" {
		logc.Error(ctx, "LDAP 同步 Cron 表达式为空，跳过启动")
		return
	}

	c := cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
	))

	entryID, err := c.AddFunc(ldapConfig.Cronjob, func() {
		// 创建一个带有超时的 Context 用于单次同步任务，防止卡死
		taskCtx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		logc.Infof(taskCtx, "触发 LDAP 定时同步任务")
		l.SyncUserToW8t(ldapConfig)
	})

	if err != nil {
		logc.Error(ctx, "添加 LDAP 同步定时任务失败: %v", err)
		return
	}

	logc.Infof(ctx, "启动 LDAP 定时同步任务, EntryID: %d, Cron: %s", entryID, ldapConfig.Cronjob)

	c.Start()

	<-ctx.Done()
	logc.Infof(ctx, "停止 LDAP 定时同步任务")
	c.Stop()
}

func (l ldapService) SyncNow() error {
	setting, err := l.ctx.DB.Setting().Get()
	if err != nil {
		logc.Error(l.ctx.Ctx, "获取LDAP配置失败: %s", err.Error())
		return err
	}

	if setting.LdapConfig.Cronjob == "" {
		return fmt.Errorf("LDAP 未配置定时同步")
	}

	l.SyncUserToW8t(setting.LdapConfig)
	return nil
}
