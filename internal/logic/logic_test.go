package logic_test

import (
	"context"
	"errors"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	llm "github.com/Mark-Grigorev/chat_analyzer/internal/clients/llm/mocks"
	tg "github.com/Mark-Grigorev/chat_analyzer/internal/clients/telegram/mocks"
	"github.com/Mark-Grigorev/chat_analyzer/internal/config"
	"github.com/Mark-Grigorev/chat_analyzer/internal/logic"
	"github.com/Mark-Grigorev/chat_analyzer/internal/settings"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setEnv(t *testing.T) {
	t.Setenv("LLM_URL", "llmurl")
	t.Setenv("LLM_TOKEN", "token")
	t.Setenv("LLM_MODEL", "model")
	t.Setenv("LLM_TEMPERATURE", "0.01")
	t.Setenv("TG_TOKEN", "token")
	t.Setenv("TG_CHAT_IDS", "1")
	t.Setenv("DEBUG", "true")
	t.Setenv("TG_ADMIN_USER_ID", "100")
}

func newTestSettings(t *testing.T, chatIDs []int64) *settings.Settings {
	dir := t.TempDir()
	s, err := settings.Load(filepath.Join(dir, "settings.json"), "system prompt", chatIDs, 0.01)
	require.NoError(t, err)
	return s
}

func TestLogicOK_NotScam(t *testing.T) {
	setEnv(t)
	update := tgbotapi.Update{
		UpdateID: 1,
		Message: &tgbotapi.Message{
			Text: "hello",
			Chat: &tgbotapi.Chat{ID: 1},
		},
	}
	updates := make(chan tgbotapi.Update, 1)
	updates <- update
	close(updates)

	tgMock := tg.NewTelegram(t)
	tgMock.On("GetUpdatesChan").Return(tgbotapi.UpdatesChannel(updates), nil)

	llmMock := llm.NewLLM(t)
	llmMock.On("GetLLMResponseAboutMsg", mock.Anything, mock.Anything, update.Message.Text, mock.Anything).Return("0", nil)

	cfg, err := config.Read()
	require.NoError(t, err)
	assert.NoError(t, logic.New(cfg, tgMock, llmMock, newTestSettings(t, []int64{1}), slog.Default()).Start(context.Background()))
	tgMock.AssertExpectations(t)
	llmMock.AssertExpectations(t)
}

func TestLogicOK_ScamDetected(t *testing.T) {
	setEnv(t)
	update := tgbotapi.Update{
		UpdateID: 1,
		Message: &tgbotapi.Message{
			MessageID: 42,
			Text:      "Заработай миллион за день, пиши в лс",
			Chat:      &tgbotapi.Chat{ID: 1},
		},
	}
	updates := make(chan tgbotapi.Update, 1)
	updates <- update
	close(updates)

	tgMock := tg.NewTelegram(t)
	tgMock.On("GetUpdatesChan").Return(tgbotapi.UpdatesChannel(updates), nil)
	tgMock.On("DeleteMessage", update.Message.Chat.ID, update.Message.MessageID).Return(nil)

	llmMock := llm.NewLLM(t)
	llmMock.On("GetLLMResponseAboutMsg", mock.Anything, mock.Anything, update.Message.Text, mock.Anything).Return("1", nil)

	cfg, err := config.Read()
	require.NoError(t, err)
	assert.NoError(t, logic.New(cfg, tgMock, llmMock, newTestSettings(t, []int64{1}), slog.Default()).Start(context.Background()))
	tgMock.AssertExpectations(t)
	llmMock.AssertExpectations(t)
}

func TestLogic_DeleteError(t *testing.T) {
	setEnv(t)
	update := tgbotapi.Update{
		UpdateID: 1,
		Message: &tgbotapi.Message{
			MessageID: 42,
			Text:      "Заработай миллион за день, пиши в лс",
			Chat:      &tgbotapi.Chat{ID: 1},
		},
	}
	updates := make(chan tgbotapi.Update, 1)
	updates <- update
	close(updates)

	tgMock := tg.NewTelegram(t)
	tgMock.On("GetUpdatesChan").Return(tgbotapi.UpdatesChannel(updates), nil)
	tgMock.On("DeleteMessage", update.Message.Chat.ID, update.Message.MessageID).Return(errors.New("delete fail"))

	llmMock := llm.NewLLM(t)
	llmMock.On("GetLLMResponseAboutMsg", mock.Anything, mock.Anything, update.Message.Text, mock.Anything).Return("1", nil)

	cfg, err := config.Read()
	require.NoError(t, err)
	// ошибка удаления логируется, но не останавливает бота
	assert.NoError(t, logic.New(cfg, tgMock, llmMock, newTestSettings(t, []int64{1}), slog.Default()).Start(context.Background()))
	tgMock.AssertExpectations(t)
	llmMock.AssertExpectations(t)
}

func TestLogic_BadUpdateChannel(t *testing.T) {
	setEnv(t)
	tgMock := tg.NewTelegram(t)
	tgMock.On("GetUpdatesChan").Return(nil, errors.New("fail"))

	llmMock := llm.NewLLM(t)
	cfg, err := config.Read()
	require.NoError(t, err)
	errCh := make(chan error)

	go func() {
		errCh <- logic.New(cfg, tgMock, llmMock, newTestSettings(t, []int64{1}), slog.Default()).Start(context.Background())
	}()

	time.Sleep(10 * time.Millisecond)
	assert.Error(t, <-errCh, "fail")
	tgMock.AssertExpectations(t)
}

func TestLogic_MessageNil(t *testing.T) {
	setEnv(t)
	update := tgbotapi.Update{UpdateID: 1}
	updates := make(chan tgbotapi.Update, 1)
	updates <- update

	tgMock := tg.NewTelegram(t)
	tgMock.On("GetUpdatesChan").Return(tgbotapi.UpdatesChannel(updates), nil)

	llmMock := llm.NewLLM(t)

	cfg, err := config.Read()
	require.NoError(t, err)

	errCh := make(chan error, 1)

	go func() {
		errCh <- logic.New(cfg, tgMock, llmMock, newTestSettings(t, []int64{1}), slog.Default()).Start(context.Background())
	}()

	time.Sleep(10 * time.Millisecond)
	select {
	case err := <-errCh:
		assert.Error(t, err, "any error")
		tgMock.AssertExpectations(t)
	default:
		tgMock.AssertExpectations(t)
	}
}

func TestLogic_WrongChatID(t *testing.T) {
	setEnv(t)
	update := tgbotapi.Update{
		UpdateID: 1,
		Message: &tgbotapi.Message{
			Text: "hello",
			Chat: &tgbotapi.Chat{ID: 999}, // не входит в TG_CHAT_IDS (1)
		},
	}
	updates := make(chan tgbotapi.Update, 1)
	updates <- update
	close(updates)

	tgMock := tg.NewTelegram(t)
	tgMock.On("GetUpdatesChan").Return(tgbotapi.UpdatesChannel(updates), nil)

	llmMock := llm.NewLLM(t)

	cfg, err := config.Read()
	require.NoError(t, err)

	assert.NoError(t, logic.New(cfg, tgMock, llmMock, newTestSettings(t, []int64{1}), slog.Default()).Start(context.Background()))
	tgMock.AssertExpectations(t)
	llmMock.AssertExpectations(t)
}

func TestLogic_LLMError(t *testing.T) {
	setEnv(t)
	update := tgbotapi.Update{
		UpdateID: 1,
		Message: &tgbotapi.Message{
			Text: "hello",
			Chat: &tgbotapi.Chat{ID: 1},
		},
	}
	updates := make(chan tgbotapi.Update, 1)
	updates <- update
	close(updates)

	tgMock := tg.NewTelegram(t)
	tgMock.On("GetUpdatesChan").Return(tgbotapi.UpdatesChannel(updates), nil)

	llmMock := llm.NewLLM(t)
	llmMock.On("GetLLMResponseAboutMsg", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("", errors.New("llm fail"))

	cfg, err := config.Read()
	require.NoError(t, err)

	assert.NoError(t, logic.New(cfg, tgMock, llmMock, newTestSettings(t, []int64{1}), slog.Default()).Start(context.Background()))
	tgMock.AssertExpectations(t)
	llmMock.AssertExpectations(t)
}

func TestLogic_AdminSetPrompt(t *testing.T) {
	setEnv(t)

	updates := make(chan tgbotapi.Update, 1)
	update := tgbotapi.Update{
		UpdateID: 1,
		Message: &tgbotapi.Message{
			Text: "/setprompt new prompt text",
			Chat: &tgbotapi.Chat{ID: 200, Type: "private"},
			From: &tgbotapi.User{ID: 100},
		},
	}
	updates <- update
	close(updates)

	tgMock := tg.NewTelegram(t)
	tgMock.On("GetUpdatesChan").Return(tgbotapi.UpdatesChannel(updates), nil)
	tgMock.On("Send", mock.Anything).Return(tgbotapi.Message{}, nil)

	llmMock := llm.NewLLM(t)

	cfg, err := config.Read()
	require.NoError(t, err)

	s := newTestSettings(t, []int64{1})
	assert.NoError(t, logic.New(cfg, tgMock, llmMock, s, slog.Default()).Start(context.Background()))

	assert.Equal(t, "new prompt text", s.GetSystemPrompt())
	tgMock.AssertExpectations(t)
	llmMock.AssertExpectations(t)
}
