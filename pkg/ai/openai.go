package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"watchAlert/config"
	"watchAlert/pkg/tools"
)

type (

	// OpenaiClient 客户端结构
	OpenaiClient struct {
		Url       string
		ApiKey    string
		Model     string
		Timeout   int
		Stream    bool
		MaxTokens int
	}

	// OpenaiOption 配置选项
	OpenaiOption func(*OpenaiClient)

	OpenaiRequest struct {
		Model       string           `json:"model"`
		Messages    []*OpenaiMessage `json:"messages"`
		Stream      bool             `json:"stream,omitempty"`
		MaxTokens   int              `json:"max_tokens,omitempty"`
		Temperature float64          `json:"temperature,omitempty"`
	}

	// OpenaiMessage 消息结构
	OpenaiMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	// OpenaiResponse 响应结构
	OpenaiResponse struct {
		ID      string `json:"id"`
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	// OpenaiStreamChunk 流式响应块
	OpenaiStreamChunk struct {
		Choices []struct {
			Delta struct {
				Content string `json:"content"`
			} `json:"delta"`
		} `json:"choices"`
	}
)

// NewOpenAIClient 工厂方法
func NewOpenAIClient(config *config.OpenAIConfig, opt ...OpenaiOption) AiClient {
	openAi := &OpenaiClient{
		Url:       config.Url,
		ApiKey:    config.AppKey,
		MaxTokens: config.MaxTokens,
		Model:     config.Model,
	}
	for _, o := range opt {
		o(openAi)
	}
	return openAi
}

// WithOpenAiTimeout 设置超时时间
func WithOpenAiTimeout(timeout int) OpenaiOption {
	return func(o *OpenaiClient) {
		o.Timeout = timeout
	}
}

func (o *OpenaiClient) ChatCompletion(_ context.Context, prompt string) (string, error) {
	// 构造请求消息
	messages := []*OpenaiMessage{
		{Role: "user", Content: prompt},
	}
	// 组装请求参数
	reqParams := OpenaiRequest{
		Model:    o.Model,
		Messages: messages,
		Stream:   false,
	}

	bodyBytes, _ := json.Marshal(reqParams)
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + o.ApiKey
	response, err := tools.Post(headers, o.Url, bytes.NewReader(bodyBytes), o.Timeout)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		errorBody, _ := io.ReadAll(response.Body)
		var errResp OpenaiResponse
		_ = json.Unmarshal(errorBody, &errResp)
		return "", fmt.Errorf("OpenAI API错误: %d - %s", response.StatusCode, errResp.Error.Message)
	}

	// 解析响应
	var result OpenaiResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("响应解析失败: %w", err)
	}

	// 检查有效响应
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("无有效返回内容")
	}

	return result.Choices[0].Message.Content, nil
}

func (o *OpenaiClient) StreamCompletion(ctx context.Context, prompt string) (<-chan string, error) {
	reqParams := OpenaiRequest{
		Model:     o.Model,
		Messages:  []*OpenaiMessage{{Role: "user", Content: prompt}},
		Stream:    true,
		MaxTokens: o.MaxTokens,
	}
	bodyBytes, _ := json.Marshal(reqParams)
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + o.ApiKey

	response, err := tools.Post(headers, o.Url, bytes.NewReader(bodyBytes), o.Timeout)
	if err != nil {
		return nil, fmt.Errorf("流式请求失败: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		errorBody, _ := io.ReadAll(response.Body)
		var errResp OpenaiResponse
		_ = json.Unmarshal(errorBody, &errResp)
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
					if err := json.Unmarshal([]byte(content), &chunk); err != nil {
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

func (o *OpenaiClient) Check(_ context.Context) error {
	if o.Url == "" || o.ApiKey == "" {
		return fmt.Errorf("OpenAI API配置错误")
	}

	if o.Timeout == 0 {
		return fmt.Errorf("OpenAI API超时时间未设置")
	}
	return nil
}
