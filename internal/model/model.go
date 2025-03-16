package model

type Config struct {
	LLMConfig      LLM
	TelegramConfig Telegram
}

type LLM struct {
	URL         string  `envconfig:"LLM_URL" required:"true"`
	Token       string  `envconfig:"LLM_TOKEN" required:"true"`
	Model       string  `envconfig:"LLM_MODEL" required:"true"`
	Temperature float64 `envconfig:"LLM_TEMPERATURE" required:"true"`
}

type Telegram struct {
	Token   string  `envconfig:"TG_TOKEN" required:"true"`
	ChatIDS []int64 `envconfig:"TG_CHAT_IDS" required:"true"`
}
