package cache

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logc"
	"golang.org/x/net/context"
	"sync"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

type (
	FaultCenterCache struct {
		rc *redis.Client
		sync.RWMutex
	}

	FaultCenterCacheInterface interface {
		PushFaultCenterInfo(center models.FaultCenter)
		GetFaultCenterInfo(faultCenterInfoKey string) models.FaultCenter
		RemoveFaultCenterInfo(faultCenterInfoKey string)
	}
)

// newFaultCenterCacheInterface 创建一个新的 FaultCenterCache 实例
func newFaultCenterCacheInterface(r *redis.Client) FaultCenterCacheInterface {
	return &FaultCenterCache{
		rc: r,
	}
}

// PushFaultCenterInfo 添加 Info 数据
func (f *FaultCenterCache) PushFaultCenterInfo(center models.FaultCenter) {
	err := f.rc.Set(center.GetFaultCenterInfoKey(), tools.JsonMarshal(center), 0).Err()
	if err != nil {
		logc.Errorf(context.Background(), err.Error())
		return
	}
}

// GetFaultCenterInfo 获取 Info 数据
func (f *FaultCenterCache) GetFaultCenterInfo(faultCenterInfoKey string) models.FaultCenter {
	result, err := f.rc.Get(faultCenterInfoKey).Result()
	if err != nil {
		return models.FaultCenter{}
	}

	var fc models.FaultCenter
	_ = json.Unmarshal([]byte(result), &fc)
	return fc
}

// RemoveFaultCenterInfo 删除 Info 数据
func (f *FaultCenterCache) RemoveFaultCenterInfo(faultCenterInfoKey string) {
	f.rc.Del(faultCenterInfoKey)
}
