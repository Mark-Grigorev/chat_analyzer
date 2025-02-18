package logic

import (
	"context"
	"fmt"

	"github.com/Mark-Grigorev/chat_analyzer/internal/clients"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

type logic struct {
	tgBot        *tgBotAPI.BotAPI
	chatGPT      *clients.Gpt
	updateConfig tgBotAPI.UpdateConfig
	log          log.Logger
}

func New(
	telegram *clients.Telegram,
	chatGPT *clients.Gpt,
	log log.Logger,
) *logic {
	return &logic{
		tgBot:        telegram.Bot,
		updateConfig: telegram.UpdateConfig,
		chatGPT:      chatGPT,
		log:          log,
	}
}

func (l *logic) Start(ctx context.Context) {
	l.log.Infoln("Bot starting")
	logPrefix := "[Start]"
	updates, err := l.tgBot.GetUpdatesChan(l.updateConfig)
	if err != nil {
		l.log.Fatalf("%s error - %s", logPrefix, err.Error())
	}

	for update := range updates {
		if update.Message != nil {
			fmt.Println(update.Message)
			response, err := l.chatGPT.SendMessage(ctx, update.Message.Text)
			if err != nil {
				l.log.Fatalf("%s error - %s", logPrefix, err.Error())
			}
			msg := tgBotAPI.NewMessage(update.Message.Chat.ID, response)
			_, err = l.tgBot.Send(msg)
			if err != nil {
				l.log.Fatalf("%s errors - %s", logPrefix, err.Error())
			}
		}
	}
}
