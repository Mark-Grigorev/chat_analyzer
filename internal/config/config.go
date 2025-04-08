package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug          string `envconfig:"DEBUG" required:"true"`
	LLMConfig      LLMConfig
	TelegramConfig TelegramConfig
}

type LLMConfig struct {
	URL         string  `envconfig:"LLM_URL" required:"true"`
	Token       string  `envconfig:"LLM_TOKEN" required:"true"`
	Model       string  `envconfig:"LLM_MODEL" required:"true"`
	Temperature float64 `envconfig:"LLM_TEMPERATURE" required:"true"`
}

type TelegramConfig struct {
	Token   string  `envconfig:"TG_TOKEN" required:"true"`
	ChatIDS []int64 `envconfig:"TG_CHAT_IDS" required:"true"`
}

func Read() *Config {
	var config Config
	envconfig.Process("", &config)
	return &config
}
