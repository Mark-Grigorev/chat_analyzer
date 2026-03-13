package logic

import (
	"context"
	"fmt"
	"strings"

	botroutes "github.com/Mark-Grigorev/chat_analyzer/internal/bot-routes"
	"github.com/Mark-Grigorev/chat_analyzer/internal/clients/llm"
	telegram "github.com/Mark-Grigorev/chat_analyzer/internal/clients/telegram"
	"github.com/Mark-Grigorev/chat_analyzer/internal/config"
	"github.com/Mark-Grigorev/chat_analyzer/internal/settings"

	"log/slog"
)

type Logic interface {
	Start(ctx context.Context) error
}

type logic struct {
	tgBot    telegram.Telegram
	llm      llm.LLM
	settings *settings.Settings
	router   *botroutes.Router
	log      *slog.Logger
}

func New(
	cfg *config.Config,
	tg telegram.Telegram,
	llm llm.LLM,
	s *settings.Settings,
	log *slog.Logger,
) *logic {
	return &logic{
		tgBot:    tg,
		llm:      llm,
		settings: s,
		router:   botroutes.New(tg, s, cfg.TelegramConfig.AdminUserID, log),
		log:      log,
	}
}

func (l *logic) Start(ctx context.Context) error {
	op := "[Start]"
	l.log.Info("Bot starting")
	updates, err := l.tgBot.GetUpdatesChan()
	if err != nil {
		return fmt.Errorf("%s - %s", op, err)
	}

	for {
		select {
		case <-ctx.Done():
			l.log.Info("Bot shutting down", "reason", ctx.Err())
			return nil
		case update, ok := <-updates:
			if !ok {
				return nil
			}

			if update.Message == nil {
				continue
			}

			if update.Message.Chat.Type == "private" {
				l.router.HandleAdminMessage(ctx, update.Message)
				continue
			}

			chatIDs := l.settings.GetChatIDs()
			for i := range chatIDs {
				if update.Message.Chat.ID == chatIDs[i] {
					l.log.Debug("info - msg - " + update.Message.Text)
					resp, err := l.llm.GetLLMResponseAboutMsg(
						ctx,
						l.settings.GetSystemPrompt(),
						update.Message.Text,
						l.settings.GetTemperature(),
					)
					if err != nil {
						l.log.Error("llm error - " + err.Error())
						continue
					}
					if strings.TrimSpace(resp) == "1" {
						l.log.Info("scam detected, deleting message",
							"chat_id", update.Message.Chat.ID,
							"message_id", update.Message.MessageID,
							"message", update.Message.Text,
						)
						if err = l.tgBot.DeleteMessage(update.Message.Chat.ID, update.Message.MessageID); err != nil {
							l.log.Error("delete message error - " + err.Error())
						}
					}
				}
			}
		}
	}
}
