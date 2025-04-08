package logic

import (
	"context"

	"github.com/Mark-Grigorev/chat_analyzer/internal/clients/llm"
	telegram "github.com/Mark-Grigorev/chat_analyzer/internal/clients/telegram"
	"github.com/Mark-Grigorev/chat_analyzer/internal/config"

	"log/slog"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	llmRulePromt = "Проанализируй данное сообщение и ответь 1 если считаешь что это скорее всего не человек, и 0 если человек(учитывай системное сообщение)"
)

type logic struct {
	tgBot   telegram.Telegram
	llm     llm.LLM
	chatIDs []int64
	log     *slog.Logger
}

func New(
	config *config.Config,
	telegram telegram.Telegram,
	llm llm.LLM,
	log *slog.Logger,
) *logic {
	return &logic{
		tgBot:   telegram,
		llm:     llm,
		log:     log,
		chatIDs: config.TelegramConfig.ChatIDS,
	}
}

func (l *logic) Start(ctx context.Context) {
	l.log.Info("Bot starting")
	updates, err := l.tgBot.GetUpdatesChan()
	if err != nil {
		l.log.Error("error - " + err.Error())
		return
	}

	for update := range updates {
		if update.Message != nil {
			for i := range l.chatIDs {
				if update.Message.Chat.ID == l.chatIDs[i] {
					l.log.Debug("info - msg - " + update.Message.Text)
					resp, err := l.llm.GetLLMResponseAboutMsg(ctx, llmRulePromt+update.Message.Text)
					if err != nil {
						l.log.Error("llm error - " + err.Error())
						continue
					}
					_, err = l.tgBot.Send(tgBotAPI.NewMessage(update.Message.Chat.ID, resp))
					if err != nil {
						l.log.Error("new msg error - " + err.Error())
						continue
					}
				}

			}

		}
	}
}
