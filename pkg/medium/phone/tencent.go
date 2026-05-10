package phone

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"watchAlert/internal/types"
	"watchAlert/pkg/tools"

	"github.com/zeromicro/go-zero/core/logc"
)

// TencentPhoneNotifier 腾讯云电话通知器
type TencentPhoneNotifier struct {
	SecretID  string
	SecretKey string
	AppID     string
}

// NewTencentPhoneNotifier 创建腾讯云电话通知器
func NewTencentPhoneNotifier(secretID, secretKey, appID string) *TencentPhoneNotifier {
	return &TencentPhoneNotifier{
		SecretID:  secretID,
		SecretKey: secretKey,
		AppID:     appID,
	}
}

// GetType 获取通知器类型
func (t *TencentPhoneNotifier) GetType() types.NotificationType {
	return types.Phone
}

// GetProvider 获取通知器提供商
func (t *TencentPhoneNotifier) GetProvider() types.NotificationProvider {
	return types.TencentPhone
}

// Notify 发送腾讯云电话通知
func (t *TencentPhoneNotifier) Notify(ctx context.Context, message *types.Message) (*types.NotificationResult, error) {
	logc.Infof(ctx, "开始发送腾讯云电话通知, Users: %v, Labels: %v", message.ToUsers, message.Labels)

	// 验证配置
	if err := t.validate(); err != nil {
		logc.Errorf(ctx, "腾讯云电话配置验证失败: %v", err)
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

		result := t.Call(ctx, string(content), to.Phone)
		results = append(results, fmt.Sprintf("用户: %s, %s", to.Phone, result))

		// 简单判断是否有错误（根据实际返回结果优化）
		if !strings.Contains(result, "success") && !strings.Contains(result, "成功") && !strings.Contains(result, "OK") {
			hasError = true
		}
	}

	success := !hasError && len(results) > 0
	resultMessage := fmt.Sprintf("发送结果: %v", results)

	if success {
		logc.Infof(ctx, "腾讯云电话通知发送成功")
	}

	return &types.NotificationResult{
		Success: success,
		Message: resultMessage,
		Data:    results,
	}, nil
}

// Call 拨打腾讯云电话的内部方法
func (t *TencentPhoneNotifier) Call(ctx context.Context, content, phoneNumber string) string {
	logc.Infof(ctx, "开始拨打电话到 %s，内容: %s", phoneNumber, content)

	// 腾讯云语音通知API接口地址
	url := "https://cloud.tim.qq.com/v5/tlsvoicesvr/sendtvoice?sdkappid=" + t.AppID + "&random=7226249334"

	// 准备请求数据
	reqData := map[string]interface{}{
		"tel": map[string]string{
			"nationcode": "86", // 国家码，中国为86
			"mobile":     phoneNumber,
		},
		"prompttype": 2,       // 语音通知类型，2表示播放文本转语音
		"promptfile": content, // 语音内容，如果是文本转语音则传入文本内容
		"playtimes":  2,       // 播放次数
	}

	// 计算签名
	timeStr := strconv.FormatInt(time.Now().Unix(), 10)
	timeInt, _ := strconv.Atoi(timeStr)

	// 计算签名，腾讯云语音通知的签名方式与短信类似
	strRand := "7226249334"
	sigContent := "appkey=" + t.SecretKey + "&random=" + strRand + "&time=" + timeStr + "&mobile=" + phoneNumber
	sig := getSha256Code(sigContent)
	reqData["sig"] = sig
	reqData["time"] = timeInt

	// 将请求数据转换为JSON
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		errMsg := fmt.Sprintf("构建腾讯云电话请求数据失败: %s", err.Error())
		logc.Error(ctx, errMsg)
		return errMsg
	}

	// 设置请求头
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	// 发送POST请求
	response, err := tools.Post(headers, url, bytes.NewReader(jsonData), 10)
	if err != nil {
		errMsg := fmt.Sprintf("拨打腾讯云电话请求失败: %s", err.Error())
		logc.Error(ctx, errMsg)
		return errMsg
	}

	respStr, _ := io.ReadAll(response.Body)
	logc.Infof(ctx, "腾讯云电话响应: %s", string(respStr))

	// 检查响应是否包含成功标识
	if strings.Contains(string(respStr), "\"result\":0") || strings.Contains(string(respStr), "success") { // 根据腾讯云API响应格式判断
		return "success: " + string(respStr)
	} else {
		return "error: " + string(respStr)
	}
}

// getSha256Code 计算SHA256摘要
func getSha256Code(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// validate 验证配置参数
func (t *TencentPhoneNotifier) validate() error {
	if t.SecretID == "" {
		return fmt.Errorf("SecretID 不能为空")
	}
	if t.SecretKey == "" {
		return fmt.Errorf("SecretKey 不能为空")
	}
	if t.AppID == "" {
		return fmt.Errorf("AppID 不能为空")
	}
	return nil
}
