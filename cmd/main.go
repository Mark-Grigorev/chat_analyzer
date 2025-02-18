package main

import (
	"context"
	"log"

	"github.com/Mark-Grigorev/chat_analyzer/internal/clients"
	"github.com/Mark-Grigorev/chat_analyzer/internal/config"
	"github.com/Mark-Grigorev/chat_analyzer/internal/logic"

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
