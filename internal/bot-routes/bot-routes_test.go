package botroutes_test

import (
	"context"
	"log/slog"
	"path/filepath"
	"testing"

	botroutes "github.com/Mark-Grigorev/chat_analyzer/internal/bot-routes"
	tg "github.com/Mark-Grigorev/chat_analyzer/internal/clients/telegram/mocks"
	"github.com/Mark-Grigorev/chat_analyzer/internal/settings"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const adminID = int64(100)

func newTestSettings(t *testing.T, chatIDs []int64) *settings.Settings {
	t.Helper()
	s, err := settings.Load(filepath.Join(t.TempDir(), "settings.json"), "system prompt", chatIDs, 0.5)
	require.NoError(t, err)
	return s
}

func newRouter(tgMock *tg.Telegram, s *settings.Settings) *botroutes.Router {
	return botroutes.New(tgMock, s, adminID, slog.Default())
}

func adminMsg(chatID int64, userID int, text string) *tgbotapi.Message {
	return &tgbotapi.Message{
		Text: text,
		Chat: &tgbotapi.Chat{ID: chatID},
		From: &tgbotapi.User{ID: userID},
	}
}

func expectSend(tgMock *tg.Telegram, text string) {
	tgMock.On("Send", mock.MatchedBy(func(m tgbotapi.MessageConfig) bool {
		return m.Text == text
	})).Return(tgbotapi.Message{}, nil).Once()
}

func expectSendContains(tgMock *tg.Telegram, substr string) {
	tgMock.On("Send", mock.MatchedBy(func(m tgbotapi.MessageConfig) bool {
		return len(m.Text) > 0 && containsStr(m.Text, substr)
	})).Return(tgbotapi.Message{}, nil).Once()
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}

// --- Unauthorized ---

func TestHandleAdminMessage_NilFrom_Unauthorized(t *testing.T) {
	tgMock := tg.NewTelegram(t)
	msg := &tgbotapi.Message{Text: "/help", Chat: &tgbotapi.Chat{ID: 1}, From: nil}
	expectSend(tgMock, "Unauthorized")

	newRouter(tgMock, newTestSettings(t, nil)).HandleAdminMessage(context.Background(), msg)
}

func TestHandleAdminMessage_WrongUserID_Unauthorized(t *testing.T) {
	tgMock := tg.NewTelegram(t)
	msg := adminMsg(1, 999, "/help")
	expectSend(tgMock, "Unauthorized")

	newRouter(tgMock, newTestSettings(t, nil)).HandleAdminMessage(context.Background(), msg)
}

// --- /setprompt ---

func TestHandleAdminMessage_SetPrompt(t *testing.T) {
	tgMock := tg.NewTelegram(t)
	s := newTestSettings(t, nil)
	expectSend(tgMock, "System prompt updated")

	newRouter(tgMock, s).HandleAdminMessage(context.Background(), adminMsg(1, 100, "/setprompt new text"))

	assert.Equal(t, "new text", s.GetSystemPrompt())
}

// --- /addchat ---

func TestHandleAdminMessage_AddChat_Valid(t *testing.T) {
	tgMock := tg.NewTelegram(t)
	s := newTestSettings(t, nil)
	expectSend(tgMock, "Chat 42 added")

	newRouter(tgMock, s).HandleAdminMessage(context.Background(), adminMsg(1, 100, "/addchat 42"))

	assert.Equal(t, []int64{42}, s.GetChatIDs())
}

func TestHandleAdminMessage_AddChat_InvalidArg(t *testing.T) {
	tgMock := tg.NewTelegram(t)
	expectSend(tgMock, "Invalid chat ID")

	newRouter(tgMock, newTestSettings(t, nil)).HandleAdminMessage(context.Background(), adminMsg(1, 100, "/addchat notanumber"))
}

// --- /removechat ---

func TestHandleAdminMessage_RemoveChat_Valid(t *testing.T) {
	tgMock := tg.NewTelegram(t)
	s := newTestSettings(t, []int64{10, 20})
	expectSend(tgMock, "Chat 10 removed")

	newRouter(tgMock, s).HandleAdminMessage(context.Background(), adminMsg(1, 100, "/removechat 10"))

	assert.Equal(t, []int64{20}, s.GetChatIDs())
}

func TestHandleAdminMessage_RemoveChat_InvalidArg(t *testing.T) {
	tgMock := tg.NewTelegram(t)
	expectSend(tgMock, "Invalid chat ID")

	newRouter(tgMock, newTestSettings(t, nil)).HandleAdminMessage(context.Background(), adminMsg(1, 100, "/removechat bad"))
}

// --- /listchats ---

func TestHandleAdminMessage_ListChats_Empty(t *testing.T) {
	tgMock := tg.NewTelegram(t)
	expectSend(tgMock, "No chats configured")

	newRouter(tgMock, newTestSettings(t, nil)).HandleAdminMessage(context.Background(), adminMsg(1, 100, "/listchats"))
}

func TestHandleAdminMessage_ListChats_WithChats(t *testing.T) {
	tgMock := tg.NewTelegram(t)
	expectSendContains(tgMock, "Chats:")

	newRouter(tgMock, newTestSettings(t, []int64{1, 2})).HandleAdminMessage(context.Background(), adminMsg(1, 100, "/listchats"))
}

// --- /settemperature ---

func TestHandleAdminMessage_SetTemperature_Valid(t *testing.T) {
	tgMock := tg.NewTelegram(t)
	s := newTestSettings(t, nil)
	expectSend(tgMock, "Temperature set to 0.80")

	newRouter(tgMock, s).HandleAdminMessage(context.Background(), adminMsg(1, 100, "/settemperature 0.8"))

	assert.InDelta(t, 0.8, s.GetTemperature(), 1e-9)
}

func TestHandleAdminMessage_SetTemperature_NonNumeric(t *testing.T) {
	tgMock := tg.NewTelegram(t)
	expectSend(tgMock, "Invalid temperature, must be between 0.0 and 1.0")

	newRouter(tgMock, newTestSettings(t, nil)).HandleAdminMessage(context.Background(), adminMsg(1, 100, "/settemperature hot"))
}

func TestHandleAdminMessage_SetTemperature_OutOfRange(t *testing.T) {
	tgMock := tg.NewTelegram(t)
	expectSend(tgMock, "Invalid temperature, must be between 0.0 and 1.0")

	newRouter(tgMock, newTestSettings(t, nil)).HandleAdminMessage(context.Background(), adminMsg(1, 100, "/settemperature 1.5"))
}

// --- /status ---

func TestHandleAdminMessage_Status(t *testing.T) {
	tgMock := tg.NewTelegram(t)
	s := newTestSettings(t, []int64{5})
	expectSendContains(tgMock, "Prompt:")

	newRouter(tgMock, s).HandleAdminMessage(context.Background(), adminMsg(1, 100, "/status"))
}

// --- /help ---

func TestHandleAdminMessage_Help(t *testing.T) {
	tgMock := tg.NewTelegram(t)
	expectSendContains(tgMock, "/setprompt")

	newRouter(tgMock, newTestSettings(t, nil)).HandleAdminMessage(context.Background(), adminMsg(1, 100, "/help"))
}

// --- unknown ---

func TestHandleAdminMessage_UnknownCommand(t *testing.T) {
	tgMock := tg.NewTelegram(t)
	expectSend(tgMock, "Unknown command. Use /help")

	newRouter(tgMock, newTestSettings(t, nil)).HandleAdminMessage(context.Background(), adminMsg(1, 100, "/unknown"))
}
