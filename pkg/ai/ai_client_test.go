package ai

import (
	"context"
	"testing"
	"watchAlert/internal/models"
)

func TestOpenAiChatCompletion(t *testing.T) {

	c := &models.AiConfig{
		Type:      "openai",
		Url:       "https://free.v36.cm/v1/chat/completions",
		AppKey:    "sk-ZauXu0adURlYp8JBCa7e2a52C7Fd433c9eE33a094226CeEf",
		Model:     "gpt-4o-mini",
		MaxTokens: 2048,
	}
	client := NewOpenAIClient(c, WithOpenAiTimeout(300))
	client.Check(context.Background())

	resp, err := client.ChatCompletion(context.Background(), "你好，你是谁？")
	if err != nil {
		t.Fatal("completion err", err)
		return
	}

	t.Log("resp:", resp)
}

func TestOpenAiStreamCompletion(t *testing.T) {

	c := &models.AiConfig{
		Type:      "openai",
		Url:       "https://free.v36.cm/v1/chat/completions",
		AppKey:    "sk-ZauXu0adURlYp8JBCa7e2a52C7Fd433c9eE33a094226CeEf",
		Model:     "gpt-4o-mini",
		MaxTokens: 2048,
	}
	client := NewOpenAIClient(c, WithOpenAiTimeout(30))
	client.Check(context.Background())

	resp, err := client.StreamCompletion(context.Background(), "你好，你是谁？")
	if err != nil {
		t.Fatal("streamCompletion err", err)
		return
	}

	var received []string
	for part := range resp {
		received = append(received, part)
	}
	t.Log("received:", received)
}
