package types

import (
	"context"
	"watchAlert/internal/models"
)

// NotificationType 通知类型枚举
type NotificationType string

const (
	Email    NotificationType = "Email"
	SMS      NotificationType = "SMS"
	Phone    NotificationType = "Phone"
	FeiShu   NotificationType = "FeiShu"
	DingDing NotificationType = "DingDing"
	WeChat   NotificationType = "WeChat"
	WebHook  NotificationType = "WebHook"
	Slack    NotificationType = "Slack"
)

// NotificationProvider 通知提供商枚举
type NotificationProvider string

const (
	AliyunSms    NotificationProvider = "AliyunSms"
	TencentSms   NotificationProvider = "TencentSms"
	AliyunPhone  NotificationProvider = "AliyunPhone"
	TencentPhone NotificationProvider = "TencentPhone"
)

// Message 通知消息结构
type Message struct {
	ToUsers []models.Member `json:"toUsers"`
	Labels  map[string]any  `json:"labels"`
	Content string          `json:"content"`
	// Subject string          `json:"subject"`
}

// NotificationResult 通知发送结果
type NotificationResult struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NotificationInterface 通知发送器接口
type NotificationInterface interface {
	GetType() NotificationType
	GetProvider() NotificationProvider
	Notify(ctx context.Context, message *Message) (*NotificationResult, error)
}
