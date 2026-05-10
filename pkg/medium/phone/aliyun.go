package phone

import (
	"context"
	"encoding/json"
	"fmt"
	"watchAlert/internal/types"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dyvmsapi"
	"github.com/zeromicro/go-zero/core/logc"
)

// AliyunPhoneNotifier 阿里云电话通知器
type AliyunPhoneNotifier struct {
	AccessKeyId      string
	AccessKeySecret  string
	CalledShowNumber string
	TtsCode          string
	Region           string
}

// NewAliyunPhoneNotifier 创建阿里云电话通知器
func NewAliyunPhoneNotifier(accessKeyId, accessKeySecret, calledShowNumber, ttsCode string) *AliyunPhoneNotifier {
	return &AliyunPhoneNotifier{
		AccessKeyId:      accessKeyId,
		AccessKeySecret:  accessKeySecret,
		CalledShowNumber: calledShowNumber,
		TtsCode:          ttsCode,
		Region:           "cn-hangzhou",
	}
}

// GetType 获取通知器类型
func (a *AliyunPhoneNotifier) GetType() types.NotificationType {
	return types.Phone
}

// GetProvider 获取通知器提供商
func (a *AliyunPhoneNotifier) GetProvider() types.NotificationProvider {
	return types.AliyunPhone
}

// Notify 发送阿里云电话通知
func (a *AliyunPhoneNotifier) Notify(ctx context.Context, message *types.Message) (*types.NotificationResult, error) {
	logc.Infof(ctx, "开始发送阿里云电话通知, Users: %v, Labels: %v", message.ToUsers, message.Labels)

	// 验证配置
	if err := a.validate(); err != nil {
		logc.Errorf(ctx, "阿里云电话配置验证失败: %v", err)
		return &types.NotificationResult{
			Success: false,
			Message: fmt.Sprintf("配置验证失败: %v", err),
		}, err
	}

	// 检查必要参数
	if len(message.ToUsers) == 0 {
		logc.Errorf(ctx, "电话号码不能为空")
		return &types.NotificationResult{
			Success: false,
			Message: "电话号码不能为空",
		}, fmt.Errorf("电话号码不能为空")
	}

	// 提取消息内容
	content, err := json.Marshal(message.Labels)
	if err != nil {
		return &types.NotificationResult{
			Success: false,
			Message: "Labels JSON 转换失败: " + err.Error(),
		}, nil
	}

	// 批量发送电话通知
	var results []string
	var hasError bool

	for _, to := range message.ToUsers {
		if to.Phone == "" {
			continue
		}

		result := a.Call(ctx, string(content), to.Phone)
		results = append(results, fmt.Sprintf("用户: %s, %s", to.Phone, result))

		// 简单判断是否有错误（根据实际返回结果优化）
		if result != "OK" && result != "成功" {
			hasError = true
		}
	}

	success := !hasError && len(results) > 0
	resultMessage := fmt.Sprintf("发送结果: %v", results)

	if success {
		logc.Infof(ctx, "阿里云电话通知发送成功")
	}

	return &types.NotificationResult{
		Success: success,
		Message: resultMessage,
		Data:    results,
	}, nil
}

// Call 拨打电话
func (a *AliyunPhoneNotifier) Call(ctx context.Context, content, phoneNumber string) string {
	logc.Infof(ctx, "开始拨打电话到 %s，内容: %s", phoneNumber, content)

	// 创建阿里云电话客户端
	client, err := dyvmsapi.NewClientWithAccessKey(a.Region, a.AccessKeyId, a.AccessKeySecret)
	if err != nil {
		errMsg := fmt.Sprintf("创建阿里云电话客户端失败: %s", err.Error())
		logc.Error(ctx, errMsg)
		return errMsg
	}

	// 创建拨打电话请求
	request := dyvmsapi.CreateSingleCallByTtsRequest()
	request.Scheme = "https"
	request.CalledShowNumber = a.CalledShowNumber
	request.CalledNumber = phoneNumber
	request.TtsCode = a.TtsCode

	// 构建TTS参数
	ttsParam := map[string]string{
		"content": content,
	}
	ttsParamJson, err := json.Marshal(ttsParam)
	if err != nil {
		errMsg := fmt.Sprintf("构建TTS参数失败: %s", err.Error())
		logc.Error(ctx, errMsg)
		return errMsg
	}
	request.TtsParam = string(ttsParamJson)

	// 拨打电话
	response, err := client.SingleCallByTts(request)
	if err != nil {
		errMsg := fmt.Sprintf("拨打阿里云电话失败: %s", err.Error())
		return errMsg
	}

	logc.Infof(ctx, "电话拨打完成，响应: %s", response.Message)
	return response.Message
}

// validate 验证配置参数
func (a *AliyunPhoneNotifier) validate() error {
	if a.AccessKeyId == "" {
		return fmt.Errorf("AccessKeyId 不能为空")
	}
	if a.AccessKeySecret == "" {
		return fmt.Errorf("AccessKeySecret 不能为空")
	}
	if a.CalledShowNumber == "" {
		return fmt.Errorf("CalledShowNumber 不能为空")
	}
	if a.TtsCode == "" {
		return fmt.Errorf("TtsCode 不能为空")
	}
	return nil
}
