package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logc"
)

const (
	// LeaderElectionKey Leader 选举的 Redis Key
	LeaderElectionKey = "w8t:leader"
	// LeaderTTL Leader 锁的过期时间（秒）
	LeaderTTL = 5
	// LeaderRenewInterval Leader 锁续期间隔（秒）
	LeaderRenewInterval = 2
	// LeaderCheckInterval 检查 Leader 状态的间隔（秒）
	LeaderCheckInterval = 1
)

// LeaderElector Leader 选举器
type LeaderElector struct {
	client         *redis.Client
	ctx            context.Context
	instanceID     string
	isLeader       bool
	cancelRenew    context.CancelFunc
	onBecomeLeader func()
	onLoseLeader   func()
}

// NewLeaderElector 创建 Leader 选举器
func NewLeaderElector(ctx context.Context, client *redis.Client, onBecomeLeader, onLoseLeader func()) *LeaderElector {
	return &LeaderElector{
		client:         client,
		ctx:            ctx,
		instanceID:     uuid.New().String(),
		isLeader:       false,
		onBecomeLeader: onBecomeLeader,
		onLoseLeader:   onLoseLeader,
	}
}

// Start 开始 Leader 选举
func (le *LeaderElector) Start() {
	logc.Infof(le.ctx, "实例 ID: %s", le.instanceID)

	go le.electionLoop()
}

// electionLoop 选举循环
func (le *LeaderElector) electionLoop() {
	ticker := time.NewTicker(time.Second * LeaderCheckInterval)
	defer ticker.Stop()

	// 选举
	le.tryBecomeLeader()

	for {
		select {
		case <-ticker.C:
			le.checkLeaderStatus()
		case <-le.ctx.Done():
			le.resign()
			return
		}
	}
}

// tryBecomeLeader 尝试成为 Leader
func (le *LeaderElector) tryBecomeLeader() {
	// 尝试获取 Leader 锁
	ok, err := le.client.SetNX(LeaderElectionKey, le.instanceID, time.Duration(LeaderTTL)*time.Second).Result()
	if err != nil {
		logc.Errorf(le.ctx, "Leader 选举失败: %v", err)
		return
	}

	if ok {
		// 成功获取 Leader 锁
		le.promoteToLeader()
	} else {
		// 未获取到锁，检查当前 Leader
		currentLeader, err := le.client.Get(LeaderElectionKey).Result()
		if err == redis.Nil {
			// Leader 已失效，重试
			le.tryBecomeLeader()
		} else if err != nil {
			logc.Errorf(le.ctx, "获取当前 Leader 失败: %v", err)
		} else if currentLeader == le.instanceID {
			// 自己已经是 Leader，但可能是从 Follower 恢复
			if !le.isLeader {
				le.promoteToLeader()
			}
		} else {
			// 其他实例是 Leader
			if le.isLeader {
				le.demoteToFollower()
			}
		}
	}
}

// promoteToLeader 提升为 Leader
func (le *LeaderElector) promoteToLeader() {
	if le.isLeader {
		return
	}

	logc.Infof(le.ctx, "当前实例成为 Leader，ID: %s", le.instanceID)
	le.isLeader = true

	// 启动心跳续期
	le.startHeartbeat()

	// 成为 Leader 时的操作（加载规则）
	if le.onBecomeLeader != nil {
		go le.onBecomeLeader()
	}
}

// demoteToFollower 降级为 Follower
func (le *LeaderElector) demoteToFollower() {
	if !le.isLeader {
		return
	}

	logc.Infof(le.ctx, "当前实例降级为 Follower，ID: %s", le.instanceID)
	le.isLeader = false

	// 停止心跳续期
	if le.cancelRenew != nil {
		le.cancelRenew()
	}

	// 失去 Leader 时的操作（停止所有任务）
	if le.onLoseLeader != nil {
		go le.onLoseLeader()
	}
}

// startHeartbeat 启动心跳续期
func (le *LeaderElector) startHeartbeat() {
	renewCtx, cancel := context.WithCancel(le.ctx)
	le.cancelRenew = cancel

	go func() {
		ticker := time.NewTicker(time.Duration(LeaderRenewInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := le.renewLeadership(); err != nil {
					logc.Errorf(le.ctx, "Leader 心跳续期失败: %v", err)
					// 续期失败，失去 Leader 身份
					le.demoteToFollower()
					return
				}
			case <-renewCtx.Done():
				return
			}
		}
	}()
}

// renewLeadership 续期 Leader 身份
func (le *LeaderElector) renewLeadership() error {
	// Lua 脚本确保只续期自己持有的锁
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("expire", KEYS[1], ARGV[2])
		else
			return 0
		end
	`

	result, err := le.client.Eval(script, []string{LeaderElectionKey}, le.instanceID, LeaderTTL).Result()
	if err != nil {
		return fmt.Errorf("续期失败: %v", err)
	}

	if result.(int64) == 0 {
		return fmt.Errorf("当前实例已不是 Leader")
	}

	return nil
}

// checkLeaderStatus 检查 Leader 状态
func (le *LeaderElector) checkLeaderStatus() {
	currentLeader, err := le.client.Get(LeaderElectionKey).Result()

	if err == redis.Nil {
		// Leader 已失效，尝试成为 Leader
		if !le.isLeader {
			logc.Infof(le.ctx, "检测到 Leader 缺失，尝试竞选...")
			le.tryBecomeLeader()
		}
	} else if err != nil {
		logc.Errorf(le.ctx, "检查 Leader 状态失败: %v", err)
	} else if currentLeader != le.instanceID && le.isLeader {
		// 自己以为是 Leader，但实际上不是
		logc.Infof(le.ctx, "检测到 Leader 变更，主动降级")
		le.demoteToFollower()
	} else if currentLeader == le.instanceID && !le.isLeader {
		// 自己是 Leader 但状态未更新
		logc.Infof(le.ctx, "重新确认 Leader 身份")
		le.promoteToLeader()
	}
}

// resign 主动辞去 Leader
func (le *LeaderElector) resign() {
	if !le.isLeader {
		return
	}

	logc.Infof(le.ctx, "当前实例主动辞去 Leader，ID: %s", le.instanceID)

	// Lua 脚本确保只删除自己持有的锁
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`

	_, err := le.client.Eval(script, []string{LeaderElectionKey}, le.instanceID).Result()
	if err != nil {
		logc.Errorf(le.ctx, "辞去 Leader 失败: %v", err)
	}

	le.demoteToFollower()
}

// IsLeader 判断当前实例是否是 Leader
func (le *LeaderElector) IsLeader() bool {
	return le.isLeader
}

// GetLeaderID 获取当前 Leader 的实例 ID
func (le *LeaderElector) GetLeaderID() (string, error) {
	leaderID, err := le.client.Get(LeaderElectionKey).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("当前没有 Leader")
	}
	if err != nil {
		return "", err
	}
	return leaderID, nil
}

// GetInstanceID 获取当前实例 ID
func (le *LeaderElector) GetInstanceID() string {
	return le.instanceID
}
