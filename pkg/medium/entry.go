package medium

import (
	"fmt"
	"strconv"
	"time"
	"watchAlert/internal/ctx"

	"github.com/bytedance/sonic"

	"watchAlert/internal/models"

	"github.com/zeromicro/go-zero/core/logc"
)

type (
	// SendParams 定义发送参数
	SendParams struct {
		// 基础
		TenantId string
		EventId  string
		RuleName string
		Severity string
		// 通知
		NoticeType string
		NoticeId   string
		NoticeName string
		// 恢复通知
		IsRecovered bool
		// hook 地址
		Hook string
		// 邮件
		Email models.Email
		// 短信
		SMS models.SMS
		// 电话
		Phone models.Phone
		// 消息
		Content string
		// 签名
		Sign string `json:"sign,omitempty"`
	}

	// SendInter 发送通知的接口
	SendInter interface {
		Send(params SendParams) error
		Test(params SendParams) error
	}
)

const RobotTestContent = "这是一条来自 WatchAlert 的测试消息"

// Sender 发送通知的主函数
func Sender(ctx *ctx.Context, sendParams SendParams) error {
	// 根据通知类型获取对应的发送器
	sender, err := senderFactory(sendParams.NoticeType)
	if err != nil {
		return fmt.Errorf("Send alarm failed, %s", err.Error())
	}

	// 发送通知
	if err := sender.Send(sendParams); err != nil {
		addRecord(ctx, sendParams, 1, sendParams.Content, err.Error())
		return fmt.Errorf("Send alarm failed to %s, err: %s", sendParams.NoticeType, err.Error())
	}

	// 记录成功发送的日志
	addRecord(ctx, sendParams, 0, sendParams.Content, "success")
	logc.Info(ctx.Ctx, fmt.Sprintf("Send alarm ok, msg: %s", sendParams.Content))
	return nil
}

// Tester 发送测试消息
func Tester(ctx *ctx.Context, sendParams SendParams) error {
	sender, err := senderFactory(sendParams.NoticeType)
	if err != nil {
		return fmt.Errorf("Send alarm failed, %s", err.Error())
	}

	// 发送通知
	if err := sender.Test(sendParams); err != nil {
		return fmt.Errorf("Test alarm failed to %s, err: %s", sendParams.NoticeType, err.Error())
	}

	return nil
}

// senderFactory 创建发送器的工厂函数
func senderFactory(noticeType string) (SendInter, error) {
	switch noticeType {
	case "Email":
		return NewEmailSender()
	case "FeiShu":
		return NewFeiShuSender(), nil
	case "DingDing":
		return NewDingSender(), nil
	case "WeChat":
		return NewWeChatSender(), nil
	case "WebHook":
		return NewWebHookSender(), nil
	case "Slack":
		return NewSlackSender(), nil
	case "SMS":
		return getSMSSender()
	case "Phone":
		return getPhoneSender()
	default:
		return nil, fmt.Errorf("无效的通知类型: %s", noticeType)
	}
}

// addRecord 记录通知发送结果
func addRecord(ctx *ctx.Context, sendParams SendParams, status int, msg, errMsg string) {
	err := ctx.DB.Notice().AddRecord(models.NoticeRecord{
		EventId:  sendParams.EventId,
		Date:     time.Now().Format("2006-01-02"),
		CreateAt: time.Now().Unix(),
		TenantId: sendParams.TenantId,
		RuleName: sendParams.RuleName,
		NType:    sendParams.NoticeType,
		NObj:     sendParams.NoticeId,
		Severity: sendParams.Severity,
		Status:   status,
		AlarmMsg: msg,
		ErrMsg:   errMsg,
	})
	if err != nil {
		logc.Errorf(ctx.Ctx, "Add notice record failed, err: %s", err.Error())
	}
}

// GetSendMsg 发送内容
func (s *SendParams) GetSendMsg() map[string]any {
	msg := make(map[string]any)
	if s == nil || s.Content == "" {
		return msg
	}
	err := sonic.Unmarshal([]byte(s.Content), &msg)
	if err != nil {
		logc.Errorf(ctx.Ctx, "发送的内容解析失败, err: %s", err.Error())
		return msg
	}
	return msg
}

// getSMSSender 获取短信发送器配置
func getSMSSender() (SendInter, error) {
	setting, err := ctx.DB.Setting().Get()
	if err != nil {
		return nil, fmt.Errorf("获取系统配置失败: %v", err)
	}

	// 根据配置选择短信提供商
	switch setting.CommunicationConfig.SMS.Provider {
	case "aliyun":
		return NewAliyunSMSSender(
			setting.CommunicationConfig.SMS.Aliyun.AccessKeyId,
			setting.CommunicationConfig.SMS.Aliyun.AccessKeySecret,
			setting.CommunicationConfig.SMS.Aliyun.SignName,
			setting.CommunicationConfig.SMS.Aliyun.TemplateCode,
		)
	case "tencent":
		templateId, _ := strconv.Atoi(setting.CommunicationConfig.SMS.Tencent.TemplateId)
		return NewTencentSMSSender(
			setting.CommunicationConfig.SMS.Tencent.AppKey,
			setting.CommunicationConfig.SMS.Tencent.SdkAppId,
			templateId,
			setting.CommunicationConfig.SMS.Tencent.Sign,
		)
	}

	return nil, fmt.Errorf("未配置短信提供商")
}

// getPhoneSender 获取电话发送器配置
func getPhoneSender() (SendInter, error) {
	setting, err := ctx.DB.Setting().Get()
	if err != nil {
		return nil, fmt.Errorf("获取系统配置失败: %v", err)
	}

	// 根据配置选择电话提供商
	switch setting.CommunicationConfig.Phone.Provider {
	case "aliyun":
		return NewAliyunPhoneSender(
			setting.CommunicationConfig.Phone.Aliyun.AccessKeyId,
			setting.CommunicationConfig.Phone.Aliyun.AccessKeySecret,
			setting.CommunicationConfig.Phone.Aliyun.CalledShowNumber,
			setting.CommunicationConfig.Phone.Aliyun.TtsCode,
		)
	case "tencent":
		return NewTencentPhoneSender(
			setting.CommunicationConfig.Phone.Tencent.SecretID,
			setting.CommunicationConfig.Phone.Tencent.SecretKey,
			setting.CommunicationConfig.Phone.Tencent.AppID,
		)
	}

	return nil, fmt.Errorf("未配置电话提供商")
}
