package cache

import (
	"github.com/bytedance/sonic"
	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logc"
	"golang.org/x/net/context"
	"sync"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

type (
	TopologyCache struct {
		rc *redis.Client
		sync.RWMutex
	}

	TopologyCacheInterface interface {
		PushTopologyInfo(topology models.Topology)
		GetTopologyInfo(topologyInfoKey models.TopologyCacheKey) models.Topology
		RemoveTopologyInfo(topologyInfoKey models.TopologyCacheKey)
	}
)

// newTopologyCacheInterface 创建一个新的 TopologyCache 实例
func newTopologyCacheInterface(r *redis.Client) TopologyCacheInterface {
	return &TopologyCache{
		rc: r,
	}
}

// PushTopologyInfo 添加拓扑图信息到缓存
func (t *TopologyCache) PushTopologyInfo(topology models.Topology) {
	err := t.rc.Set(string(models.BuildTopologyCacheKey(topology.TenantId, topology.ID)), tools.JsonMarshalToString(topology), 0).Err()
	if err != nil {
		logc.Errorf(context.Background(), err.Error())
		return
	}
}

// GetTopologyInfo 从缓存中获取拓扑图信息
func (t *TopologyCache) GetTopologyInfo(topologyInfoKey models.TopologyCacheKey) models.Topology {
	result, err := t.rc.Get(string(topologyInfoKey)).Result()
	if err != nil {
		return models.Topology{}
	}

	var topology models.Topology
	_ = sonic.Unmarshal([]byte(result), &topology)
	return topology
}

// RemoveTopologyInfo 从缓存中删除拓扑图信息
func (t *TopologyCache) RemoveTopologyInfo(topologyInfoKey models.TopologyCacheKey) {
	t.rc.Del(string(topologyInfoKey))
}