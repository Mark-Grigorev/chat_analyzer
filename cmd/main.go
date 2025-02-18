package main

import (
	"chat_analyzer/internal/clients"
	"chat_analyzer/internal/config"
	"chat_analyzer/internal/logic"
	"context"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.Read()
	gpt := clients.NewGPT(cfg.ChatGPTConfig)
	tg, err := clients.NewTelegram(cfg.TelegramConfig.Token)
	if err != nil {
		log.Fatal(err.Error())
	}
	logic.New(tg, gpt, *logrus.New()).Start(context.Background())
}
