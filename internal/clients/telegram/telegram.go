package clients

import (
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	tgTimeOut = 60
)

type Telegram interface {
	GetUpdatesChan() (tgBotAPI.UpdatesChannel, error)
	Send(message tgBotAPI.MessageConfig) (tgBotAPI.Message, error)
}
type Client struct {
	bot          *tgBotAPI.BotAPI
	updateConfig tgBotAPI.UpdateConfig
}

func New(token string) (*Client, error) {
	bot, err := tgBotAPI.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	updateConfig := tgBotAPI.NewUpdate(0)
	updateConfig.Timeout = tgTimeOut
	return &Client{
		bot:          bot,
		updateConfig: updateConfig,
	}, nil
}

func (c *Client) GetUpdatesChan() (tgBotAPI.UpdatesChannel, error) {
	return c.bot.GetUpdatesChan(c.updateConfig)
}

func (c *Client) Send(message tgBotAPI.MessageConfig) (tgBotAPI.Message, error) {
	return c.bot.Send(message)
}
