package services

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/ai"
)

type (
	aiService struct {
		ctx *ctx.Context
	}

	InterAiService interface {
		Chat(req interface{}) (interface{}, interface{})
	}
)

func newInterAiService(ctx *ctx.Context) InterAiService {
	return &aiService{
		ctx: ctx,
	}
}

func (a aiService) Chat(req interface{}) (interface{}, interface{}) {
	setting, err := a.ctx.DB.Setting().Get()
	if err != nil {
		return nil, err
	}

	if !setting.AiConfig.GetEnable() {
		return nil, fmt.Errorf("未开启 Ai 分析能力")
	}

	r := req.(*types.RequestAiChatContent)
	err = r.ValidateParams()
	if err != nil {
		return nil, err
	}

	client, err := a.ctx.Redis.ProviderPools().GetClient("AiClient")
	if err != nil {
		return "", err
	}

	aiClient := client.(ai.AiClient)
	prompt := setting.AiConfig.Prompt
	prompt = strings.ReplaceAll(prompt, "{{ RuleName }}", r.RuleName)
	prompt = strings.ReplaceAll(prompt, "{{ Content }}", r.Content)
	prompt = strings.ReplaceAll(prompt, "{{ SearchQL }}", r.SearchQL)
	r.Content = prompt

	switch r.Deep {
	case "true":
		r.Content = fmt.Sprintf("注意, 请深度思考下面的问题!\n%s", r.Content)
		completion, err := aiClient.ChatCompletion(a.ctx.Ctx, r.Content)
		if err != nil {
			return "", err
		}
		err = a.ctx.DB.Ai().Update(models.AiContentRecord{
			RuleId:  r.RuleId,
			Content: completion,
		})
		if err != nil {
			return nil, err
		}

		return completion, nil

	default:
		data, exist, err := a.ctx.DB.Ai().Get(r.RuleId)
		if err != nil && err != gorm.ErrRecordNotFound {
			return "", err
		}
		if exist {
			return data.Content, nil
		}

		completion, err := aiClient.ChatCompletion(a.ctx.Ctx, r.Content)
		if err != nil {
			return "", err
		}

		err = a.ctx.DB.Ai().Create(models.AiContentRecord{
			RuleId:  r.RuleId,
			Content: completion,
		})
		if err != nil {
			return nil, err
		}

		return completion, nil
	}
}
