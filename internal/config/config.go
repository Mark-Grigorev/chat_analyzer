package config

import (
	"chat_analyzer/internal/model"
	"log"
	"os"
)

func Read() *model.AppConfig {
	var config model.AppConfig
	config.ChatGPTConfig = readChatGPTConfig()
	config.TelegramConfig = readTelegramConfig()
	return &config
}

func readChatGPTConfig() model.ChatGPTConfig {
	var config model.ChatGPTConfig
	config.URL = getEnv("GPT_URL")
	config.Token = getEnv("GPT_TOKEN")
	config.Model = getEnv("GPT_MODEL")
	return config
}

func readTelegramConfig() model.TelegramConfig {
	var config model.TelegramConfig
	config.Token = getEnv("TELEGRAM_TOKEN")
	return config
}

func getEnv(key string) string {
	var data string
	if data = os.Getenv(key); data == "" {
		log.Fatalf("не указан %s", key)
	}
	return data
}
