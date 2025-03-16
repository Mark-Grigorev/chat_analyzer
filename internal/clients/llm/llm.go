package llm

import (
	"context"
	"net/http"
	"time"

	"github.com/Mark-Grigorev/chat_analyzer/internal/model"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

const (
	systemType = `
	Ты ассистент в онлайн чате, анализирующий входящие сообщения. 
	Твоя задача понять, на сколько входящее сообщение является попыткой обмануть/ввести в заблуждение, или является предложением мутного заработка.`
)

type LLMClient interface {
	GetLLMResponseAboutMsg(ctx context.Context, promt string) (string, error)
}

type Client struct {
	llm         *openai.LLM
	temperature float64
}

func New(cfg model.LLM) (*Client, error) {
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
		return &Client{}, nil
	}
	return &Client{
		llm:         llm,
		temperature: cfg.Temperature,
	}, nil
}

func (c *Client) GetLLMResponseAboutMsg(ctx context.Context, promt string) (string, error) {
	messages := []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{llms.TextContent{Text: systemType}},
		},
		{
			Role:  llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{llms.TextPart(promt)},
		},
	}

	resp, err := c.llm.GenerateContent(ctx, messages)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Content, nil
}
