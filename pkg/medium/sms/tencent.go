package sms

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

	"github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/logc"
)

// TencentSMSNotifier 腾讯云短信通知器
type TencentSMSNotifier struct {
	AppKey     string
	SdkAppId   string
	TemplateId int
	Sign       string
}

// Mobiles 手机号码结构
type Mobiles struct {
	Mobile     string `json:"mobile"`
	Nationcode string `json:"nationcode"`
}

// TXmessage 腾讯云短信消息结构
type TXmessage struct {
	Ext    string    `json:"ext"`
	Extend string    `json:"extend"`
	Params []string  `json:"params"`
	Sig    string    `json:"sig"`
	Sign   string    `json:"sign"`
	Tel    []Mobiles `json:"tel"`
	Time   int       `json:"time"`
	Tpl_id int       `json:"tpl_id"`
}

// NewTencentSMSNotifier 创建腾讯云短信通知器
func NewTencentSMSNotifier(appKey, sdkAppId string, templateId int, sign string) *TencentSMSNotifier {
	return &TencentSMSNotifier{
		AppKey:     appKey,
		SdkAppId:   sdkAppId,
		TemplateId: templateId,
		Sign:       sign,
	}
}

// GetType 获取通知器类型
func (t *TencentSMSNotifier) GetType() types.NotificationType {
	return types.SMS
}

// GetProvider 获取通知器提供商
func (t *TencentSMSNotifier) GetProvider() types.NotificationProvider {
	return types.TencentSms
}

// Notify 发送腾讯云短信通知
func (t *TencentSMSNotifier) Notify(ctx context.Context, message *types.Message) (*types.NotificationResult, error) {
	if err := t.validate(); err != nil {
		logrus.Errorf("腾讯云短信配置验证失败: %v", err)
		return &types.NotificationResult{
			Success: false,
			Message: fmt.Sprintf("配置验证失败: %v", err),
		}, err
	}

	// 检查必要参数
	if len(message.ToUsers) == 0 {
		return &types.NotificationResult{
			Success: false,
			Message: "手机号码不能为空",
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

	// 调用腾讯云短信发送函数
	for _, to := range message.ToUsers {
		if to.Phone == "" {
			continue
		}

		result, err := t.Post(string(content), to.Phone, "sms")
		if err != nil {
			hasError = true
		}
		results = append(results, fmt.Sprintf("用户: %s, %s", to.Phone, result))
	}

	success := !hasError && len(results) > 0
	resultMessage := fmt.Sprintf("发送结果: %v", results)

	if success {
		logc.Info(ctx, "腾讯云短信通知发送成功")
	}

	return &types.NotificationResult{
		Success: success,
		Message: resultMessage,
		Data:    results,
	}, nil
}

// Post 发送腾讯云短信
func (t *TencentSMSNotifier) Post(Messages, PhoneNumbers, logsign string) (string, error) {
	// 检查配置是否完整
	if t.AppKey == "" || t.SdkAppId == "" || t.TemplateId == 0 || t.Sign == "" {
		logrus.Info("腾讯云短信接口配置不完整")
		return "", fmt.Errorf("腾讯云短信接口配置不完整")
	}

	TXmobiles := []Mobiles{}
	mobiles := splitPhoneNumbers(PhoneNumbers)
	for _, m := range mobiles {
		TXmobiles = append(TXmobiles, Mobiles{
			Mobile:     m,
			Nationcode: "86",
		})
	}

	strRand := "7226249334"
	strTime := strconv.FormatInt(time.Now().Unix(), 10)
	intTime, _ := strconv.Atoi(strTime)
	sig := t.getSha256Code("appkey=" + t.AppKey + "&random=" + strRand + "&time=" + strTime + "&mobile=" + PhoneNumbers)
	url := "https://yun.tim.qq.com/v5/tlssmssvr/sendmultisms2?sdkappid=" + t.SdkAppId + "&random=" + strRand

	u := TXmessage{
		Ext:    logsign,
		Extend: "",
		Params: []string{Messages},
		Sig:    sig,
		Sign:   t.Sign,
		Tel:    TXmobiles,
		Time:   intTime,
		Tpl_id: t.TemplateId,
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)

	// 设置请求头
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	// 发送POST请求
	response, err := tools.Post(headers, url, bytes.NewReader(b.Bytes()), 10)
	if err != nil {
		logrus.Error("发送腾讯云短信请求失败: " + err.Error())
		return "", fmt.Errorf("发送腾讯云短信请求失败: %s", err.Error())
	}

	body, _ := io.ReadAll(response.Body)
	return string(body), nil
}

// getSha256Code 计算SHA256摘要
func (t *TencentSMSNotifier) getSha256Code(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// validate 验证配置参数
func (t *TencentSMSNotifier) validate() error {
	if t.AppKey == "" {
		return fmt.Errorf("AppKey 不能为空")
	}
	if t.SdkAppId == "" {
		return fmt.Errorf("SdkAppId 不能为空")
	}
	if t.Sign == "" {
		return fmt.Errorf("Sign 不能为空")
	}
	if t.TemplateId == 0 {
		return fmt.Errorf("TemplateCode 不能为空")
	}
	return nil
}

// splitPhoneNumbers 分割手机号码
func splitPhoneNumbers(phoneNumbers string) []string {
	if phoneNumbers == "" {
		return nil
	}

	// 去重和清理
	unique := make(map[string]bool)
	var result []string

	parts := strings.Split(phoneNumbers, ",")
	for _, phone := range parts {
		phone = strings.TrimSpace(phone)
		if phone != "" && !unique[phone] {
			unique[phone] = true
			result = append(result, phone)
		}
	}

	return result
}
