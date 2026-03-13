package botroutes

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	telegram "github.com/Mark-Grigorev/chat_analyzer/internal/clients/telegram"
	"github.com/Mark-Grigorev/chat_analyzer/internal/settings"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Router struct {
	tgBot    telegram.Telegram
	settings *settings.Settings
	adminID  int64
	log      *slog.Logger
}

func New(tgBot telegram.Telegram, s *settings.Settings, adminID int64, log *slog.Logger) *Router {
	return &Router{
		tgBot:    tgBot,
		settings: s,
		adminID:  adminID,
		log:      log,
	}
}

func (r *Router) HandleAdminMessage(ctx context.Context, msg *tgbotapi.Message) {
	if msg.From == nil || int64(msg.From.ID) != r.adminID {
		if _, err := r.tgBot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Unauthorized")); err != nil {
			r.log.Error("send error - " + err.Error())
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
		if err := r.settings.SetSystemPrompt(arg); err != nil {
			r.log.Error("setprompt error - " + err.Error())
			reply = "Error saving prompt"
		} else {
			reply = "System prompt updated"
		}

	case "/addchat":
		id, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			reply = "Invalid chat ID"
		} else if err = r.settings.AddChatID(id); err != nil {
			r.log.Error("addchat error - " + err.Error())
			reply = "Error saving chat ID"
		} else {
			reply = fmt.Sprintf("Chat %d added", id)
		}

	case "/removechat":
		id, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			reply = "Invalid chat ID"
		} else if err = r.settings.RemoveChatID(id); err != nil {
			r.log.Error("removechat error - " + err.Error())
			reply = "Error removing chat ID"
		} else {
			reply = fmt.Sprintf("Chat %d removed", id)
		}

	case "/listchats":
		ids := r.settings.GetChatIDs()
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
		} else if err = r.settings.SetTemperature(t); err != nil {
			r.log.Error("settemperature error - " + err.Error())
			reply = "Error saving temperature"
		} else {
			reply = fmt.Sprintf("Temperature set to %.2f", t)
		}

	case "/status":
		prompt := r.settings.GetSystemPrompt()
		if len(prompt) > 100 {
			prompt = prompt[:100]
		}
		reply = fmt.Sprintf("Prompt: %s\nChats: %v\nTemperature: %.2f",
			prompt,
			r.settings.GetChatIDs(),
			r.settings.GetTemperature(),
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

	if _, err := r.tgBot.Send(tgbotapi.NewMessage(msg.Chat.ID, reply)); err != nil {
		r.log.Error("send error - " + err.Error())
	}
}
