package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Mark-Grigorev/chat_analyzer/internal/clients/llm"
	telegram "github.com/Mark-Grigorev/chat_analyzer/internal/clients/telegram"
	"github.com/Mark-Grigorev/chat_analyzer/internal/config"
	"github.com/Mark-Grigorev/chat_analyzer/internal/logic"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	cfg := config.Read()
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	gpt, err := llm.New(cfg.LLMConfig)
	if err != nil {
		log.Error(err.Error())
		return
	}
	tg, err := telegram.New(cfg.TelegramConfig.Token)
	if err != nil {
		log.Error(err.Error())
		return
	}
	logic.New(cfg, tg, gpt, log).Start(context.Background())
}
