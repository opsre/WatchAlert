package sender

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"
	"watchAlert/internal/ctx"

	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logc"

	"watchAlert/pkg/tools"
)

type (
	// FeiShuSender 飞书发送策略
	FeiShuSender struct{}

	FeiShuResponse struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
)

var FeiShuTestContent = fmt.Sprintf(`{
  "msg_type": "text",
  "content": {
  "text": "%s"
  }
}`, RobotTestContent)

func NewFeiShuSender() SendInter {
	return &FeiShuSender{}
}

func (f *FeiShuSender) Send(params SendParams) error {
	return f.post(params.Hook, params.Sign, params.GetSendMsg())
}

func (f *FeiShuSender) Test(params SendParams) error {
	msg := make(map[string]any)
	err := sonic.Unmarshal([]byte(FeiShuTestContent), &msg)
	if err != nil {
		logc.Errorf(ctx.Ctx, fmt.Sprintf("发送的内容解析失败, err: %s", err.Error()))
		return err
	}

	return f.post(params.Hook, params.Sign, msg)
}

func (f *FeiShuSender) post(hook, sign string, msg map[string]any) error {
	if sign != "" {
		signature, timestamp := generateFeishuSignature(sign)
		msg["sign"] = signature
		msg["timestamp"] = timestamp
	}

	msgByte := bytes.NewReader(tools.JsonMarshalToByte(msg))
	res, err := tools.Post(nil, hook, msgByte, 10)
	if err != nil {
		return err
	}

	var response FeiShuResponse
	if err := tools.ParseReaderBody(res.Body, &response); err != nil {
		return errors.New(fmt.Sprintf("Error unmarshalling Feishu response: %s", err.Error()))
	}
	if response.Code != 0 {
		return errors.New(response.Msg)
	}

	return nil
}

// generateFeishuSignature 生成 Feishu 签名
func generateFeishuSignature(secret string) (string, string) {
	// 1. Get timestamp
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// 2. Prepare the string to sign: {timestamp}\n{secret}
	stringToSign := fmt.Sprintf("%s\n%s", timestamp, secret)

	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	h.Write(data)

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), timestamp
}
