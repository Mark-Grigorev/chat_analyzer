package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Mark-Grigorev/chat_analyzer/internal/clients/llm"
	telegram "github.com/Mark-Grigorev/chat_analyzer/internal/clients/telegram"
	"github.com/Mark-Grigorev/chat_analyzer/internal/config"
	"github.com/Mark-Grigorev/chat_analyzer/internal/logic"
	"github.com/Mark-Grigorev/chat_analyzer/internal/settings"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	code := run()
	if code != 0 {
		os.Exit(code)
	}
}

func run() int {
	var err error
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetLogLoggerLevel(slog.LevelDebug)

	cfg, err := config.Read()
	if err != nil {
		log.Error(err.Error())
		return 1
	}

	gpt, err := llm.New(cfg.LLMConfig)
	if err != nil {
		log.Error(err.Error())
		return 2
	}

	tg, err := telegram.New(cfg.TelegramConfig.Token)
	if err != nil {
		log.Error(err.Error())
		return 3
	}

	s, err := settings.Load(
		cfg.SettingsPath,
		llm.DefaultSystemPrompt,
		cfg.TelegramConfig.ChatIDS,
		cfg.LLMConfig.Temperature,
	)
	if err != nil {
		log.Error(err.Error())
		return 5
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err = logic.New(cfg, tg, gpt, s, log).Start(ctx); err != nil {
		log.Error(err.Error())
		return 4
	}

	return 0
}
