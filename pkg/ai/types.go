package ai

import (
	"context"
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

	AiConfig struct {
		Url       string
		ApiKey    string
		Model     string
		Timeout   int
		Stream    bool
		MaxTokens int
	}

	Request struct {
		Model       string     `json:"model"`
		Messages    []*Message `json:"messages"`
		Stream      bool       `json:"stream,omitempty"`
		MaxTokens   int        `json:"max_tokens,omitempty"`
		Temperature float64    `json:"temperature,omitempty"`
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

	// Response 响应结构
	Response struct {
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
)
