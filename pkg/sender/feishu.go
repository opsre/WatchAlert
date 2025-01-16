package sender

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

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

func NewFeiShuSender() SendInter {
	return &FeiShuSender{}
}

func (f *FeiShuSender) Send(params SendParams) error {
	msg := params.GetSendMsg()
	if params.Sign != "" {
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		signature, err := generateSignature(params.Sign, timestamp)
		if err != nil {
			return err
		}
		msg["sign"] = signature
		msg["timestamp"] = timestamp
	}

	msgStr, _ := json.Marshal(msg)

	msgByte := bytes.NewReader(msgStr)

	res, err := tools.Post(nil, params.Hook, msgByte, 10)
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

// generateSignature 生成签名
func generateSignature(secret string, timestamp string) (string, error) {
	//timestamp + key 做sha256, 再进行base64 encode
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}
