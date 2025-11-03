package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logc"
)

const (
	// Redis 消息通道
	ChannelRuleReload        = "w8t:rule:reload"
	ChannelFaultCenterReload = "w8t:faultcenter:reload"
	ChannelProbingReload     = "w8t:probing:reload"
)

// 操作类型
const (
	ActionCreate  = "create"
	ActionUpdate  = "update"
	ActionDelete  = "delete"
	ActionEnable  = "enable"
	ActionDisable = "disable"
)

// ReloadMessage 重载消息
type ReloadMessage struct {
	Action   string `json:"action"`   // create, update, delete, enable, disable
	ID       string `json:"id"`       // 规则/故障中心/拨测 ID
	TenantID string `json:"tenantId"` // 租户 ID
	Name     string `json:"name"`     // 名称（用于日志）
}

// PublishReloadMessage 发布重载消息
func PublishReloadMessage(ctx context.Context, client *redis.Client, channel string, msg ReloadMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal reload message: %v", err)
	}

	err = client.Publish(channel, string(data)).Err()
	if err != nil {
		return fmt.Errorf("failed to publish reload message: %v", err)
	}

	logc.Infof(ctx, "[Follower] 向 Leader 发布重载消息: channel=%s, action=%s, id=%s, name=%s",
		channel, msg.Action, msg.ID, msg.Name)

	return nil
}

// SubscribeReloadMessages 订阅重载消息
func SubscribeReloadMessages(ctx context.Context, client *redis.Client, channel string, handler func(msg ReloadMessage)) {
	pubsub := client.Subscribe(channel)
	defer pubsub.Close()

	logc.Infof(ctx, "[Leader] 开始订阅消息: channel=%s", channel)

	// 等待订阅确认
	_, err := pubsub.Receive()
	if err != nil {
		logc.Errorf(ctx, "Failed to subscribe to channel %s: %v", channel, err)
		return
	}

	// 接收消息
	ch := pubsub.Channel()
	for {
		select {
		case redisMsg := <-ch:
			var msg ReloadMessage
			if err := json.Unmarshal([]byte(redisMsg.Payload), &msg); err != nil {
				logc.Errorf(ctx, "Failed to unmarshal reload message: %v", err)
				continue
			}

			logc.Infof(ctx, "[Leader] 收到重载消息: action=%s, id=%s, name=%s",
				msg.Action, msg.ID, msg.Name)

			// 调用处理函数
			handler(msg)

		case <-ctx.Done():
			logc.Infof(ctx, "[Leader] 停止订阅消息: channel=%s", channel)
			return
		}
	}
}
