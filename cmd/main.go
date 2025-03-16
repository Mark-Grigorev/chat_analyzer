package main

import (
	"context"
	"log"

	"github.com/Mark-Grigorev/chat_analyzer/internal/clients/llm"
	telegram "github.com/Mark-Grigorev/chat_analyzer/internal/clients/telegram"
	"github.com/Mark-Grigorev/chat_analyzer/internal/config"
	"github.com/Mark-Grigorev/chat_analyzer/internal/logic"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.Read()
	gpt, err := llm.New(cfg.LLMConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	tg, err := telegram.New(cfg.TelegramConfig.Token)
	if err != nil {
		log.Fatal(err.Error())
	}
	logic.New(tg, gpt, *logrus.New()).Start(context.Background())
}
