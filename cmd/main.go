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
	code := run()
	if code != 0 {
		os.Exit(code)
	}
}

func run() int {
	cfg := config.Read()

	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetLogLoggerLevel(slog.LevelDebug)

	gpt, err := llm.New(cfg.LLMConfig)
	if err != nil {
		log.Error(err.Error())
		return 1
	}

	tg, err := telegram.New(cfg.TelegramConfig.Token)
	if err != nil {
		log.Error(err.Error())
		return 2
	}

	logic.New(cfg, tg, gpt, log).Start(context.Background())

	return 0
}
