package logic

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Mark-Grigorev/chat_analyzer/internal/clients/llm"
	telegram "github.com/Mark-Grigorev/chat_analyzer/internal/clients/telegram"
	"github.com/Mark-Grigorev/chat_analyzer/internal/config"
	"github.com/Mark-Grigorev/chat_analyzer/internal/settings"

	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Logic interface {
	Start(ctx context.Context) error
}

type logic struct {
	tgBot    telegram.Telegram
	llm      llm.LLM
	settings *settings.Settings
	adminID  int64
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
		adminID:  cfg.TelegramConfig.AdminUserID,
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
				l.handleAdminMessage(ctx, update.Message)
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

func (l *logic) handleAdminMessage(ctx context.Context, msg *tgbotapi.Message) {
	if msg.From == nil || int64(msg.From.ID) != l.adminID {
		_, err := l.tgBot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Unauthorized"))
		if err != nil {
			l.log.Error("send error - " + err.Error())
		}
		return
	}

	text := strings.TrimSpace(msg.Text)
	parts := strings.SplitN(text, " ", 2)
	command := parts[0]
	arg := ""
	if len(parts) == 2 {
		arg = strings.TrimSpace(parts[1])
	}

	var reply string

	switch command {
	case "/setprompt":
		if err := l.settings.SetSystemPrompt(arg); err != nil {
			l.log.Error("setprompt error - " + err.Error())
			reply = "Error saving prompt"
		} else {
			reply = "System prompt updated"
		}

	case "/addchat":
		id, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			reply = "Invalid chat ID"
		} else if err = l.settings.AddChatID(id); err != nil {
			l.log.Error("addchat error - " + err.Error())
			reply = "Error saving chat ID"
		} else {
			reply = fmt.Sprintf("Chat %d added", id)
		}

	case "/removechat":
		id, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			reply = "Invalid chat ID"
		} else if err = l.settings.RemoveChatID(id); err != nil {
			l.log.Error("removechat error - " + err.Error())
			reply = "Error removing chat ID"
		} else {
			reply = fmt.Sprintf("Chat %d removed", id)
		}

	case "/listchats":
		ids := l.settings.GetChatIDs()
		if len(ids) == 0 {
			reply = "No chats configured"
		} else {
			sb := strings.Builder{}
			sb.WriteString("Chats:\n")
			for _, id := range ids {
				sb.WriteString(fmt.Sprintf("  %d\n", id))
			}
			reply = sb.String()
		}

	case "/settemperature":
		t, err := strconv.ParseFloat(arg, 64)
		if err != nil || t < 0.0 || t > 1.0 {
			reply = "Invalid temperature, must be between 0.0 and 1.0"
		} else if err = l.settings.SetTemperature(t); err != nil {
			l.log.Error("settemperature error - " + err.Error())
			reply = "Error saving temperature"
		} else {
			reply = fmt.Sprintf("Temperature set to %.2f", t)
		}

	case "/status":
		prompt := l.settings.GetSystemPrompt()
		if len(prompt) > 100 {
			prompt = prompt[:100]
		}
		reply = fmt.Sprintf("Prompt: %s\nChats: %v\nTemperature: %.2f",
			prompt,
			l.settings.GetChatIDs(),
			l.settings.GetTemperature(),
		)

	case "/help":
		reply = "/setprompt <text> - set system prompt\n" +
			"/addchat <id> - add chat ID\n" +
			"/removechat <id> - remove chat ID\n" +
			"/listchats - list chat IDs\n" +
			"/settemperature <0.0-1.0> - set temperature\n" +
			"/status - show current settings\n" +
			"/help - show this help"

	default:
		reply = "Unknown command. Use /help"
	}

	if _, err := l.tgBot.Send(tgbotapi.NewMessage(msg.Chat.ID, reply)); err != nil {
		l.log.Error("send error - " + err.Error())
	}
}
