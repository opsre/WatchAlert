package services

import (
	"context"
	"watchAlert/internal/ctx"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/pkg/ai"
)

type (
	settingService struct {
		ctx *ctx.Context
	}

	InterSettingService interface {
		Save(req interface{}) (interface{}, interface{})
		Get() (interface{}, interface{})
	}
)

func newInterSettingService(ctx *ctx.Context) InterSettingService {
	return settingService{
		ctx: ctx,
	}
}

func (a settingService) Save(req interface{}) (interface{}, interface{}) {
	r := req.(*models.Settings)
	dbConf, err := a.ctx.DB.Setting().Get()
	if err != nil {
		return nil, err
	}

	if a.ctx.DB.Setting().Check() {
		err := a.ctx.DB.Setting().Update(*r)
		if err != nil {
			return nil, err
		}
	} else {
		err := a.ctx.DB.Setting().Create(*r)
		if err != nil {
			return nil, err
		}
	}

	const mark = "SyncLdapUserJob"
	if r.AuthType != nil && *r.AuthType == models.SettingLdapAuth && *dbConf.AuthType != models.SettingLdapAuth {
		if cancel, exists := a.ctx.ContextMap[mark]; exists {
			cancel()
			delete(a.ctx.ContextMap, mark)
		}
		c, cancel := context.WithCancel(context.Background())
		a.ctx.ContextMap[mark] = cancel
		// 定时同步LDAP用户任务
		go LdapService.SyncUsersCronjob(c)
	} else {
		if cancel, exists := a.ctx.ContextMap[mark]; exists {
			cancel()
			delete(a.ctx.ContextMap, mark)
		}
	}

	if r.AiConfig.GetEnable() {
		client, err := ai.NewAiClient(&r.AiConfig)
		if err != nil {
			return nil, err
		}
		a.ctx.Redis.ProviderPools().SetClient("AiClient", client)
	}

	return nil, nil
}

func (a settingService) Get() (interface{}, interface{}) {
	get, err := a.ctx.DB.Setting().Get()
	if err != nil {
		return nil, err
	}
	get.AppVersion = global.Version

	return get, nil
}
