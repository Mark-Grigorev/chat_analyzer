package clients

import (
	"fmt"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	tgTimeOut = 60
)

type Telegram interface {
	GetUpdatesChan() (tgBotAPI.UpdatesChannel, error)
	Send(message tgBotAPI.MessageConfig) (tgBotAPI.Message, error)
	DeleteMessage(chatID int64, messageID int) error
}
type Client struct {
	bot          *tgBotAPI.BotAPI
	updateConfig tgBotAPI.UpdateConfig
}

func New(token string) (*Client, error) {
	op := "[New]"
	bot, err := tgBotAPI.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("%s - %s", op, err)
	}

	updateConfig := tgBotAPI.NewUpdate(0)
	updateConfig.Timeout = tgTimeOut
	return &Client{
		bot:          bot,
		updateConfig: updateConfig,
	}, nil
}

func (c *Client) GetUpdatesChan() (tgBotAPI.UpdatesChannel, error) {
	op := "[GetUpdatesChan]"
	ch, err := c.bot.GetUpdatesChan(c.updateConfig)
	if err != nil {
		return nil, fmt.Errorf("%s - %s", op, err)
	}
	return ch, err
}

func (c *Client) Send(message tgBotAPI.MessageConfig) (tgBotAPI.Message, error) {
	op := "[Send]"
	msg, err := c.bot.Send(message)
	if err != nil {
		return tgBotAPI.Message{}, fmt.Errorf("%s - %s", op, err)
	}
	return msg, nil
}

func (c *Client) DeleteMessage(chatID int64, messageID int) error {
	op := "[DeleteMessage]"
	_, err := c.bot.Send(tgBotAPI.NewDeleteMessage(chatID, messageID))
	if err != nil {
		return fmt.Errorf("%s - %s", op, err)
	}
	return nil
}
