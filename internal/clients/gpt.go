package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Mark-Grigorev/chat_analyzer/internal/model"
)

const (
	sendMessage = "/v1/chat/completions"
	timeOut     = 30 * time.Second
	maxTokens   = 50
)

type Gpt struct {
	url       string
	token     string
	model     string
	maxTokens int
	client    *http.Client
}

func NewGPT(cfg model.ChatGPTConfig) *Gpt {
	return &Gpt{
		url:       cfg.URL,
		token:     cfg.Token,
		model:     cfg.Model,
		maxTokens: maxTokens,
		client: &http.Client{
			Timeout: timeOut,
		},
	}
}

func (c *Gpt) SendMessage(ctx context.Context, prompt string) (string, error) {
	requestBody, err := json.Marshal(model.ChatGPTRequest{
		Model: c.model,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{"system", "Нужно проанализировать текст на спам."},
			{"user", prompt},
		},
		MaxTokens: c.maxTokens,
	})
	if err != nil {
		return "", fmt.Errorf("ошибка сериализации JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.url+sendMessage,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ошибка: статус-код %d", resp.StatusCode)
	}

	var response model.ChatGPTResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("ошибка декодирования JSON: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("пустой ответ от API")
	}

	return response.Choices[0].Message.Content, nil
}
