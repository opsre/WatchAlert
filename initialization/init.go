package initialization

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"golang.org/x/sync/errgroup"
	"watchAlert/alert"
	"watchAlert/config"
	"watchAlert/internal/cache"
	"watchAlert/internal/ctx"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/internal/repo"
	"watchAlert/internal/services"
	"watchAlert/pkg/ai"
	"watchAlert/pkg/tools"
)

func InitBasic() {

	// 初始化配置
	global.Config = config.InitConfig()

	dbRepo := repo.NewRepoEntry()
	rCache := cache.NewEntryCache()
	ctx := ctx.NewContext(context.Background(), dbRepo, rCache)

	services.NewServices(ctx)

	// 启用告警评估携程
	alert.Initialize(ctx)

	// 初始化权限数据
	InitPermissionsSQL(ctx)

	// 初始化角色数据
	InitUserRolesSQL(ctx)

	// 导入数据源 Client 到存储池
	importClientPools(ctx)

	// 定时任务，清理历史通知记录和历史拨测数据
	go gcHistoryData(ctx)

	r, err := ctx.DB.Setting().Get()
	if err != nil {
		logc.Error(ctx.Ctx, fmt.Sprintf("加载系统设置失败: %s", err.Error()))
		return
	}

	if *r.AuthType == models.SettingLdapAuth {
		const mark = "SyncLdapUserJob"
		c, cancel := context.WithCancel(context.Background())
		ctx.ContextMap[mark] = cancel
		go services.LdapService.SyncUsersCronjob(c)
	}

	if r.AiConfig.GetEnable() {
		client, err := ai.NewAiClient(&r.AiConfig)
		if err != nil {
			logc.Error(ctx.Ctx, fmt.Sprintf("创建 Ai 客户端失败: %s", err.Error()))
			return
		}
		ctx.Redis.ProviderPools().SetClient("AiClient", client)
	}
}

func importClientPools(ctx *ctx.Context) {
	list, err := ctx.DB.Datasource().List("", "", "", "")
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

func gcHistoryData(ctx *ctx.Context) {
	// gc probe history data and notice history record
	tools.NewCronjob("00 00 */1 * *", func() {
		err := ctx.DB.Probing().DeleteRecord()
		if err != nil {
			logc.Errorf(ctx.Ctx, "fail to delete probe history data, %s", err.Error())
		} else {
			logc.Info(ctx.Ctx, "success delete probe history data")
		}

		err = ctx.DB.Notice().DeleteRecord()
		if err != nil {
			logc.Errorf(ctx.Ctx, "fail to delete notice history record, %s", err.Error())
		} else {
			logc.Info(ctx.Ctx, "success delete notice history record")
		}
	})

}
