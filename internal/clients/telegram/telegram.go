package clients

import (
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	tgTimeOut = 60
)

type Telegram struct {
	Bot          *tgBotAPI.BotAPI
	UpdateConfig tgBotAPI.UpdateConfig
}

func New(token string) (*Telegram, error) {
	bot, err := tgBotAPI.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	updateConfig := tgBotAPI.NewUpdate(0)
	updateConfig.Timeout = tgTimeOut
	return &Telegram{
		Bot:          bot,
		UpdateConfig: updateConfig,
	}, nil
}
