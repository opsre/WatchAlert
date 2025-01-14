package sender

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"time"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
)

type (
	// SendParams 定义发送参数
	SendParams struct {
		// 基础
		TenantId string
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
		// 消息
		Content string
		// 事件
		Event interface{}
		// 电话号码
		PhoneNumber []string
	}

	// SendInter 发送通知的接口
	SendInter interface {
		Send(params SendParams) error
	}
)

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
	addRecord(ctx, sendParams, 0, sendParams.Content, "")
	logc.Info(ctx.Ctx, fmt.Sprintf("Send alarm ok, msg: %s", sendParams.Content))
	return nil
}

// senderFactory 创建发送器的工厂函数
func senderFactory(noticeType string) (SendInter, error) {
	switch noticeType {
	case "Email":
		return NewEmailSender(), nil
	case "FeiShu":
		return NewFeiShuSender(), nil
	case "DingDing":
		return NewDingSender(), nil
	case "WeChat":
		return NewWeChatSender(), nil
	case "CustomHook":
		return NewWebHookSender(), nil
	case "PhoneCall":
		return NewPhoneCallSender(), nil
	default:
		return nil, fmt.Errorf("无效的通知类型: %s", noticeType)
	}
}

// addRecord 记录通知发送结果
func addRecord(ctx *ctx.Context, sendParams SendParams, status int, msg, errMsg string) {
	err := ctx.DB.Notice().AddRecord(models.NoticeRecord{
		Date:     time.Now().Format("2006-01-02"),
		CreateAt: time.Now().Unix(),
		TenantId: sendParams.TenantId,
		RuleName: sendParams.RuleName,
		NType:    sendParams.NoticeType,
		NObj:     sendParams.NoticeName + " (" + sendParams.NoticeId + ")",
		Severity: sendParams.Severity,
		Status:   status,
		AlarmMsg: msg,
		ErrMsg:   errMsg,
	})
	if err != nil {
		logc.Errorf(ctx.Ctx, fmt.Sprintf("Add notice record failed, err: %s", err.Error()))
	}
}
