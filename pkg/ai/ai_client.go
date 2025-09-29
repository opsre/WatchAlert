package ai

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"io"
	"log"
	"net/http"
	"strings"
	"watchAlert/internal/models"

	"watchAlert/pkg/tools"
)

// NewAiClient 工厂方法
func NewAiClient(config *models.AiConfig) (AiClient, error) {
	ai := &AiConfig{
		Url:       config.Url,
		ApiKey:    config.AppKey,
		MaxTokens: config.MaxTokens,
		Model:     config.Model,
		Timeout:   config.Timeout,
	}

	err := ai.Check(context.Background())
	if err != nil {
		return nil, err
	}

	return ai, nil
}

func (o *AiConfig) ChatCompletion(_ context.Context, prompt string) (string, error) {
	// 构建请求参数
	reqParams := Request{
		Model: o.Model,
		Messages: []*Message{
			{
				Role:    "system",
				Content: "您是站点可靠性工程 (SRE) 可观测性监控专家、资深 DevOps 工程师、资深运维专家",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream:    false,
		MaxTokens: o.MaxTokens,
	}

	bodyBytes, _ := sonic.Marshal(reqParams)
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + o.ApiKey
	response, err := tools.Post(headers, o.Url, bytes.NewReader(bodyBytes), o.Timeout)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		errorBody, _ := io.ReadAll(response.Body)
		var errResp Response
		_ = sonic.Unmarshal(errorBody, &errResp)
		return "", fmt.Errorf("API 请求错误: %d - %s", response.StatusCode, errResp.Error.Message)
	}

	// 解析响应
	var result Response
	err = tools.ParseReaderBody(response.Body, &result)
	if err != nil {
		return "", err
	}

	// 检查有效响应
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("无有效返回内容")
	}

	return result.Choices[0].Message.Content, nil
}

func (o *AiConfig) StreamCompletion(ctx context.Context, prompt string) (<-chan string, error) {
	reqParams := Request{
		Model: o.Model,
		Messages: []*Message{
			{Role: "user", Content: prompt},
		},
		Stream:    true,
		MaxTokens: o.MaxTokens,
	}
	bodyBytes, _ := sonic.Marshal(reqParams)
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + o.ApiKey

	response, err := tools.Post(headers, o.Url, bytes.NewReader(bodyBytes), o.Timeout)
	if err != nil {
		return nil, fmt.Errorf("流式请求失败: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		errorBody, _ := io.ReadAll(response.Body)
		var errResp Response
		_ = sonic.Unmarshal(errorBody, &errResp)
		return nil, fmt.Errorf("OpenAI API错误: %d - %s", response.StatusCode, errResp.Error.Message)
	}

	// 创建流式通道
	streamChan := make(chan string)

	go func() {
		defer close(streamChan)
		defer response.Body.Close()
		select {
		case <-ctx.Done():
			return
		default:
			scanner := bufio.NewScanner(response.Body)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "data: ") {
					content := strings.TrimPrefix(line, "data: ")
					if strings.TrimSpace(content) == "[DONE]" {
						continue
					}

					var chunk StreamChunk
					if err := sonic.Unmarshal([]byte(content), &chunk); err != nil {
						log.Printf("解析错误: %v | 内容: %s", err, content)
						continue
					}

					// 拼接内容
					if len(chunk.Choices) > 0 {
						streamChan <- chunk.Choices[0].Delta.Content
					}
				}
			}
		}
	}()
	return streamChan, nil
}

func (o *AiConfig) Check(_ context.Context) error {
	if o.Url == "" || o.ApiKey == "" {
		return fmt.Errorf("OpenAI API配置错误")
	}

	if o.Timeout == 0 {
		return fmt.Errorf("OpenAI API超时时间未设置")
	}
	return nil
}
