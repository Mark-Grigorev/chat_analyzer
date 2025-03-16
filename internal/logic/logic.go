package logic

import (
	"context"

	"github.com/Mark-Grigorev/chat_analyzer/internal/clients/llm"
	telegram "github.com/Mark-Grigorev/chat_analyzer/internal/clients/telegram"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

type logic struct {
	tgBot        *tgBotAPI.BotAPI
	llm          llm.LLMClient
	updateConfig tgBotAPI.UpdateConfig
	chatIDs      []int64
	log          log.Logger
}

func New(
	telegram *telegram.Telegram,
	llm llm.LLMClient,
	log log.Logger,
) *logic {
	return &logic{
		tgBot:        telegram.Bot,
		updateConfig: telegram.UpdateConfig,
		llm:          llm,
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
			for _, i := range l.chatIDs {
				if update.Message.Chat.ID == l.chatIDs[i] {
					log.Debugf("info - msg - %s", update.Message.Text)
					response, err := l.llm.GetLLMResponseAboutMsg(ctx, "Проанализируй данное сообщение и ответь 1 если считаешь что это скорее всего не человек, и 0 если человек(учитывай системное сообщение)"+update.Message.Text)
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
	}
}
