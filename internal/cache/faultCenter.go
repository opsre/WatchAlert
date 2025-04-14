package cache

import (
	"encoding/json"
	"sync"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

type (
	FaultCenterCache struct {
		//rc    *redis.Client
		cache Cache
		sync.RWMutex
	}

	FaultCenterCacheInterface interface {
		PushFaultCenterInfo(center models.FaultCenter)
		GetFaultCenterInfo(faultCenterInfoKey string) models.FaultCenter
		RemoveFaultCenterInfo(faultCenterInfoKey string)
	}
)

// newFaultCenterCacheInterface 创建一个新的 FaultCenterCache 实例
func newFaultCenterCacheInterface(c Cache) FaultCenterCacheInterface {
	return &FaultCenterCache{
		cache: c,
	}
}

// PushFaultCenterInfo 添加 Info 数据
func (f *FaultCenterCache) PushFaultCenterInfo(center models.FaultCenter) {
	f.cache.SetKey(center.GetFaultCenterInfoKey(), tools.JsonMarshal(center), 0)
}

// GetFaultCenterInfo 获取 Info 数据
func (f *FaultCenterCache) GetFaultCenterInfo(faultCenterInfoKey string) models.FaultCenter {
	result, err := f.cache.GetKey(faultCenterInfoKey)
	if err != nil {
		return models.FaultCenter{}
	}

	var fc models.FaultCenter
	_ = json.Unmarshal([]byte(result), &fc)
	return fc
}

// RemoveFaultCenterInfo 删除 Info 数据
func (f *FaultCenterCache) RemoveFaultCenterInfo(faultCenterInfoKey string) {
	f.cache.DeleteKey(faultCenterInfoKey)
}
