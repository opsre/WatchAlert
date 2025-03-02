package ai

import (
	"context"
	"fmt"
	"time"
	"watchAlert/config"
)

const (
	DefaultTimeout        = 30 * time.Second
	DeepSeek       string = "deepseek"
	OpenAI         string = "openai"
)

type (
	// AiClient is the interface for AI chatbot clients.
	AiClient interface {
		// ChatCompletion returns the completion of the given input text.
		ChatCompletion(context.Context, string) (string, error)
		// StreamCompletion returns a channel that streams the completion of the given input text.
		StreamCompletion(context.Context, string) (<-chan string, error)
		// Check checks the health of the AI chatbot client.
		Check(context.Context) error
	}

	// Message is a message
	Message struct {
		Role    string `json:"role"` // system/user/assistant
		Content string `json:"content"`
	}

	// StreamChunk 流式响应结构
	StreamChunk struct {
		Choices []struct {
			Delta struct {
				Content string `json:"content"`
			} `json:"delta"`
		} `json:"choices"`
	}
)

// NewAiClient  new ai client
func NewAiClient(c *config.AiConfig) (AiClient, error) {
	switch c.Type {
	case OpenAI:
		return NewOpenAIClient(c.OpenAI, WithOpenAiTimeout(c.Timeout)), nil
	case DeepSeek:
		return NewDeepSeekClient(c.DeepSeek, WithDeepSeekTimeout(c.Timeout)), nil
	default:
		return nil, fmt.Errorf("unsupported ai type: %s", c.Type)
	}
}
