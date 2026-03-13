package llm

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Mark-Grigorev/chat_analyzer/internal/config"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

const DefaultSystemPrompt = "Ты ассистент в онлайн чате, анализирующий входящие сообщения.\nТвоя задача — определить, является ли сообщение скамом: попыткой обмануть, ввести в заблуждение или предложением сомнительного заработка.\nОтветь строго одной цифрой: 1 — если сообщение является скамом, 0 — если это обычное сообщение."

type LLM interface {
	GetLLMResponseAboutMsg(ctx context.Context, systemPrompt, userMessage string, temperature float64) (string, error)
}

type Client struct {
	llm *openai.LLM
}

func New(cfg config.LLMConfig) (*Client, error) {
	op := "[NewLLM]"
	llm, err := openai.New(
		openai.WithBaseURL(cfg.URL),
		openai.WithToken(cfg.Token),
		openai.WithModel(cfg.Model),
		openai.WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				MaxIdleConns:    100,
				IdleConnTimeout: 60 * time.Second,
			},
			Timeout: 60 * time.Second,
		}),
	)
	if err != nil {
		return &Client{}, fmt.Errorf("%s - %s", op, err)
	}
	return &Client{
		llm: llm,
	}, nil
}

func (c *Client) GetLLMResponseAboutMsg(ctx context.Context, systemPrompt, userMessage string, temperature float64) (string, error) {
	messages := []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{llms.TextContent{Text: systemPrompt}},
		},
		{
			Role:  llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{llms.TextPart(userMessage)},
		},
	}

	resp, err := c.llm.GenerateContent(
		ctx,
		messages,
		llms.WithTemperature(temperature),
	)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Content, nil
}
