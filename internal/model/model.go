package model

type AppConfig struct {
	ChatGPTConfig
	TelegramConfig
}

type ChatGPTConfig struct {
	URL   string
	Token string
	Model string
}

type TelegramConfig struct {
	Token string
}
