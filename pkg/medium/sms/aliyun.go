package sms

import (
	"context"
	"encoding/json"
	"fmt"

	"watchAlert/internal/types"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/logc"
)

// AliyunSMSNotifier 阿里云短信通知器
type AliyunSMSNotifier struct {
	AccessKeyId     string
	AccessKeySecret string
	SignName        string
	TemplateCode    string
}

// NewAliyunSMSNotifier 创建阿里云短信通知器
func NewAliyunSMSNotifier(accessKeyId, accessKeySecret, signName, templateCode string) *AliyunSMSNotifier {
	return &AliyunSMSNotifier{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		SignName:        signName,
		TemplateCode:    templateCode,
	}
}

// GetType 获取通知器类型
func (a *AliyunSMSNotifier) GetType() types.NotificationType {
	return types.SMS
}

// GetProvider 获取通知器提供商
func (a *AliyunSMSNotifier) GetProvider() types.NotificationProvider {
	return types.AliyunSms
}

// Notify 发送阿里云短信通知
func (a *AliyunSMSNotifier) Notify(ctx context.Context, message *types.Message) (*types.NotificationResult, error) {
	// 验证配置
	if err := a.validate(); err != nil {
		logc.Errorf(ctx, "阿里云短信配置验证失败: %v", err)
		return &types.NotificationResult{
			Success: false,
			Message: fmt.Sprintf("配置验证失败: %v", err),
		}, err
	}

	// 检查必要参数
	if len(message.ToUsers) == 0 {
		return &types.NotificationResult{
			Success: false,
			Message: "电话号码不能为空",
		}, nil
	}

	content, err := json.Marshal(message.Labels)
	if err != nil {
		return &types.NotificationResult{
			Success: false,
			Message: "Labels JSON 转换失败: " + err.Error(),
		}, nil
	}

	var results []string
	var hasError bool

	// 调用阿里云短信发送函数
	for _, to := range message.ToUsers {
		if to.Phone == "" {
			continue
		}

		result, err := a.Post(string(content), to.Phone, "sms")
		if err != nil {
			hasError = true
		}
		results = append(results, fmt.Sprintf("用户: %s, %s", to.Phone, result))
	}

	success := !hasError && len(results) > 0
	resultMessage := fmt.Sprintf("发送结果: %v", results)

	if success {
		logc.Info(ctx, "阿里云短信通知发送成功")
	}

	return &types.NotificationResult{
		Success: success,
		Message: resultMessage,
		Data:    results,
	}, nil
}

// Post 发送阿里云短信
func (a *AliyunSMSNotifier) Post(Messages, PhoneNumbers, logsign string) (string, error) {
	// 创建阿里云短信客户端
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", a.AccessKeyId, a.AccessKeySecret)
	if err != nil {
		logrus.Error("创建阿里云短信客户端失败, ", err.Error())
		return "", fmt.Errorf("创建阿里云短信客户端失败: %s", err.Error())
	}

	// 创建发送短信请求
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = PhoneNumbers
	request.SignName = a.SignName
	request.TemplateCode = a.TemplateCode
	request.TemplateParam = Messages

	// 发送短信
	response, err := client.SendSms(request)
	if err != nil {
		logrus.Error("阿里云短信发送失败, ", err.Error())
		return "", fmt.Errorf("阿里云短信发送失败: %s", err.Error())
	}

	logrus.Info("阿里云短信发送成功, ", response)

	return response.Message, nil
}

// validate 验证配置参数
func (a *AliyunSMSNotifier) validate() error {
	if a.AccessKeyId == "" {
		return fmt.Errorf("AccessKeyId 不能为空")
	}
	if a.AccessKeySecret == "" {
		return fmt.Errorf("AccessKeySecret 不能为空")
	}
	if a.SignName == "" {
		return fmt.Errorf("SignName 不能为空")
	}
	if a.TemplateCode == "" {
		return fmt.Errorf("TemplateCode 不能为空")
	}
	return nil
}
