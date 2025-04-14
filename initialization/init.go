package initialization

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"golang.org/x/sync/errgroup"
	"watchAlert/alert"
	"watchAlert/config"
	"watchAlert/internal/cache"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/internal/repo"
	"watchAlert/internal/services"
	"watchAlert/pkg/ai"
	"watchAlert/pkg/ctx"
)

func InitBasic() {

	// 初始化配置
	global.Config = config.InitConfig()

	dbRepo := repo.NewRepoEntry()
	rCache := cache.NewEntryCache(global.Config.Cache)
	ctx := ctx.NewContext(context.Background(), dbRepo, rCache, global.Config.Cache)

	services.NewServices(ctx)

	// 启用告警评估携程
	alert.Initialize(ctx)

	// 初始化权限数据
	InitPermissionsSQL(ctx)

	// 初始化角色数据
	InitUserRolesSQL(ctx)

	// 导入数据源 Client 到存储池
	importClientPools(ctx)

	if global.Config.Ldap.Enabled {
		// 定时同步LDAP用户任务
		go services.LdapService.SyncUsersCronjob()
	}

	r, err := ctx.DB.Setting().Get()
	if err != nil {
		logc.Error(ctx.Ctx, fmt.Sprintf("加载系统设置失败: %s", err.Error()))
		return
	}

	if r.AiConfig.GetEnable() {
		client, err := ai.NewAiClient(&r.AiConfig)
		if err != nil {
			logc.Error(ctx.Ctx, fmt.Sprintf("创建 Ai 客户端失败: %s", err.Error()))
			return
		}
		ctx.Cache.ProviderPools().SetClient("AiClient", client)
	}
}

func importClientPools(ctx *ctx.Context) {
	list, err := ctx.DB.Datasource().List(models.DatasourceQuery{})
	if err != nil {
		logc.Error(ctx.Ctx, err.Error())
		return
	}

	g := new(errgroup.Group)
	for _, datasource := range list {
		ds := datasource
		if !*ds.GetEnabled() {
			continue
		}
		g.Go(func() error {
			err := services.DatasourceService.WithAddClientToProviderPools(ds)
			if err != nil {
				logc.Error(ctx.Ctx, fmt.Sprintf("添加到 Client 存储池失败, err: %s", err.Error()))
				return err
			}
			return nil
		})
	}
}
