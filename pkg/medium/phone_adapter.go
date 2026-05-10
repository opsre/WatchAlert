package medium

import (
	"context"
	"fmt"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/medium/phone"

	"github.com/zeromicro/go-zero/core/logc"
)

// PhoneAdapter 电话通知适配器
type PhoneAdapter struct {
	notifier types.NotificationInterface
}

// NewAliyunPhoneSender 创建阿里云电话发送器
func NewAliyunPhoneSender(accessKeyId, accessKeySecret, calledShowNumber, ttsCode string) (SendInter, error) {
	notifier := phone.NewAliyunPhoneNotifier(accessKeyId, accessKeySecret, calledShowNumber, ttsCode)
	return &PhoneAdapter{notifier: notifier}, nil
}

// NewTencentPhoneSender 创建腾讯云电话发送器
func NewTencentPhoneSender(secretID, secretKey, appID string) (SendInter, error) {
	notifier := phone.NewTencentPhoneNotifier(secretID, secretKey, appID)
	return &PhoneAdapter{notifier: notifier}, nil
}

// Send 实现SendInter接口
func (p *PhoneAdapter) Send(params SendParams) error {
	// 构建用户列表
	var toUsers []models.Member
	for _, phone := range params.Phone.To {
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
	result, err := p.notifier.Notify(context.Background(), message)
	if err != nil {
		return fmt.Errorf("电话通知发送失败: %v", err)
	}

	if !result.Success {
		return fmt.Errorf("电话通知发送失败: %s", result.Message)
	}

	logc.Info(context.Background(), fmt.Sprintf("电话通知发送成功: %s", result.Message))
	return nil
}

// Test 实现SendInter接口
func (p *PhoneAdapter) Test(params SendParams) error {
	// 构建测试消息
	message := &types.Message{
		Labels: map[string]any{
			"test":     "true",
			"template": "test",
		},
		Content: RobotTestContent,
	}

	if len(params.Phone.To) > 0 {
		for _, phone := range params.Phone.To {
			message.ToUsers = append(message.ToUsers, models.Member{
				Phone: phone,
			})
		}
	}

	result, err := p.notifier.Notify(context.Background(), message)
	if err != nil {
		return fmt.Errorf("电话通知测试失败: %v", err)
	}

	if !result.Success {
		return fmt.Errorf("电话通知测试失败: %s", result.Message)
	}

	return nil
}
