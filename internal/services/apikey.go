package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
)

type apiKeyService struct {
	ctx *ctx.Context
}

type InterApiKeyService interface {
	Create(req interface{}) (interface{}, interface{})
	List(req interface{}) (interface{}, interface{})
	Get(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Delete(req interface{}) (interface{}, interface{})
	GetApiKeyByUserId(userId string) ([]types.ResponseApiKeyInfo, error)
}

func newInterApiKeyService(ctx *ctx.Context) InterApiKeyService {
	return &apiKeyService{
		ctx: ctx,
	}
}

func (aks apiKeyService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestApiKeyCreate)

	// 生成随机API密钥
	apiKey, err := generateApiKey()
	if err != nil {
		return nil, fmt.Errorf("生成API密钥失败: %v", err)
	}

	// 检查是否提供了用户ID，如果没有则从JWT中获取
	userId := r.UserId
	if userId == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	model := models.ApiKey{
		UserId:      userId,
		Name:        r.Name,
		Description: r.Description,
		Key:         apiKey,
		CreatedAt:   time.Now(),
	}

	err = aks.ctx.DB.ApiKey().Create(model)
	if err != nil {
		return nil, err
	}

	// 返回不包含敏感信息的结果
	result := types.ResponseApiKeyInfo{
		ID:          model.ID,
		UserId:      model.UserId,
		Name:        model.Name,
		Description: model.Description,
		Key:         model.Key,
		CreatedAt:   model.CreatedAt.Unix(),
	}

	return result, nil
}

func (aks apiKeyService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestApiKeyQuery)

	data, err := aks.ctx.DB.ApiKey().List(r.UserId)
	if err != nil {
		return nil, err
	}

	var result []types.ResponseApiKeyInfo
	for _, item := range data {
		result = append(result, types.ResponseApiKeyInfo{
			ID:          item.ID,
			UserId:      item.UserId,
			Name:        item.Name,
			Description: item.Description,
			Key:         item.Key,
			CreatedAt:   item.CreatedAt.Unix(),
		})
	}

	return result, nil
}

func (aks apiKeyService) Get(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestApiKeyQuery)

	data, _, err := aks.ctx.DB.ApiKey().Get(r.ID, r.UserId)
	if err != nil {
		return nil, err
	}

	result := types.ResponseApiKeyInfo{
		ID:          data.ID,
		UserId:      data.UserId,
		Name:        data.Name,
		Description: data.Description,
		Key:         data.Key,
		CreatedAt:   data.CreatedAt.Unix(),
	}

	return result, nil
}

func (aks apiKeyService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestApiKeyUpdate)

	// 首先检查API密钥是否存在且属于当前用户
	existing, _, err := aks.ctx.DB.ApiKey().Get(r.ID, r.UserId)
	if err != nil {
		return nil, err
	}

	model := models.ApiKey{
		ID:          r.ID,
		UserId:      existing.UserId, // 确保不能更改所属用户
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   existing.CreatedAt,
	}

	err = aks.ctx.DB.ApiKey().Update(model)
	if err != nil {
		return nil, err
	}

	// 返回不包含敏感信息的结果
	result := types.ResponseApiKeyInfo{
		ID:          model.ID,
		UserId:      model.UserId,
		Name:        model.Name,
		Description: model.Description,
		Key:         model.Key,
		CreatedAt:   model.CreatedAt.Unix(),
	}

	return result, nil
}

func (aks apiKeyService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestApiKeyQuery)

	err := aks.ctx.DB.ApiKey().Delete(r.ID, r.UserId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetApiKeyByUserId 根据用户ID获取API密钥列表
func (aks apiKeyService) GetApiKeyByUserId(userId string) ([]types.ResponseApiKeyInfo, error) {
	data, err := aks.ctx.DB.ApiKey().List(userId)
	if err != nil {
		return nil, err
	}

	var result []types.ResponseApiKeyInfo
	for _, item := range data {
		result = append(result, types.ResponseApiKeyInfo{
			ID:          item.ID,
			UserId:      item.UserId,
			Name:        item.Name,
			Description: item.Description,
			Key:         item.Key,
			CreatedAt:   item.CreatedAt.Unix(),
		})
	}

	return result, nil
}

// generateApiKey 生成随机API密钥
func generateApiKey() (string, error) {
	bytes := make([]byte, 32) // 256位密钥
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
