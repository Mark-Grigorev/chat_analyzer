package config_test

import (
	"fmt"
	"testing"

	"github.com/Mark-Grigorev/chat_analyzer/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	const (
		llmURL         = "llmurl"
		llmToken       = "token"
		llmModel       = "model"
		llmTemperature = "0.01"
		tgToken        = "token"
		tgIDs          = "55,22,11"
	)
	expectedChatIDs := []int64{55, 22, 11}

	t.Setenv("LLM_URL", llmURL)
	t.Setenv("LLM_TOKEN", llmToken)
	t.Setenv("LLM_MODEL", llmModel)
	t.Setenv("LLM_TEMPERATURE", llmTemperature)
	t.Setenv("TG_TOKEN", tgToken)
	t.Setenv("TG_CHAT_IDS", tgIDs)

	cfg := config.Read()
	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg)

	assert.Equal(t, cfg.LLMConfig.URL, llmURL)
	assert.Equal(t, cfg.LLMConfig.Token, llmToken)
	assert.Equal(t, cfg.LLMConfig.Model, llmModel)
	assert.Equal(t, fmt.Sprintf("%v", cfg.LLMConfig.Temperature), llmTemperature)
	assert.Equal(t, cfg.TelegramConfig.Token, tgToken)
	assert.ElementsMatch(t, cfg.TelegramConfig.ChatIDS, expectedChatIDs)
}
