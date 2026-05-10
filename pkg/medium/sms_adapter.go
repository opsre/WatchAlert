package medium

import (
	"context"
	"fmt"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/medium/sms"

	"github.com/zeromicro/go-zero/core/logc"
)

// SMSAdapter 短信通知适配器
type SMSAdapter struct {
	notifier types.NotificationInterface
}

// NewAliyunSMSSender 创建阿里云短信发送器
func NewAliyunSMSSender(accessKeyId, accessKeySecret, signName, templateCode string) (SendInter, error) {
	notifier := sms.NewAliyunSMSNotifier(accessKeyId, accessKeySecret, signName, templateCode)
	return &SMSAdapter{notifier: notifier}, nil
}

// NewTencentSMSSender 创建腾讯云短信发送器
func NewTencentSMSSender(appKey, sdkAppId string, templateId int, sign string) (SendInter, error) {
	notifier := sms.NewTencentSMSNotifier(appKey, sdkAppId, templateId, sign)
	return &SMSAdapter{notifier: notifier}, nil
}

// Send 实现SendInter接口
func (s *SMSAdapter) Send(params SendParams) error {
	// 构建用户列表
	var toUsers []models.Member
	for _, phone := range params.SMS.To {
		toUsers = append(toUsers, models.Member{
			UserId:   phone,
			UserName: phone,
			Phone:    phone,
		})
	}

	// 构建消息
	message := &types.Message{
		ToUsers: toUsers,
		Labels:  params.GetSendMsg(),
		Content: params.Content,
	}

	// 发送通知
	result, err := s.notifier.Notify(context.Background(), message)
	if err != nil {
		return fmt.Errorf("短信通知发送失败: %v", err)
	}

	if !result.Success {
		return fmt.Errorf("短信通知发送失败: %s", result.Message)
	}

	logc.Info(context.Background(), fmt.Sprintf("短信通知发送成功: %s", result.Message))
	return nil
}

// Test 实现SendInter接口
func (s *SMSAdapter) Test(params SendParams) error {
	// 构建测试消息
	message := &types.Message{
		Labels: map[string]any{
			"test":     "true",
			"template": "test",
		},
		Content: RobotTestContent,
	}

	if len(params.SMS.To) > 0 {
		for _, phone := range params.SMS.To {
			message.ToUsers = append(message.ToUsers, models.Member{
				Phone: phone,
			})
		}
	}

	result, err := s.notifier.Notify(context.Background(), message)
	if err != nil {
		return fmt.Errorf("短信通知测试失败: %v", err)
	}

	if !result.Success {
		return fmt.Errorf("短信通知测试失败: %s", result.Message)
	}

	return nil
}
