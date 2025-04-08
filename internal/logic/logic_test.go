package logic_test

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"

	llm "github.com/Mark-Grigorev/chat_analyzer/internal/clients/llm/mocks"
	tg "github.com/Mark-Grigorev/chat_analyzer/internal/clients/telegram/mocks"
	"github.com/Mark-Grigorev/chat_analyzer/internal/config"
	"github.com/Mark-Grigorev/chat_analyzer/internal/logic"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stretchr/testify/mock"
)

func setEnv(t *testing.T) {
	t.Setenv("LLM_URL", "llmurl")
	t.Setenv("LLM_TOKEN", "token")
	t.Setenv("LLM_MODEL", "model")
	t.Setenv("LLM_TEMPERATURE", "0.01")
	t.Setenv("TG_TOKEN", "token")
	t.Setenv("TG_CHAT_IDS", "1")
	t.Setenv("DEBUG", "true")
}
func TestLogicOK(t *testing.T) {
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
	llmMock.On("GetLLMResponseAboutMsg", mock.Anything, mock.MatchedBy(func(s string) bool {
		return strings.Contains(s, update.Message.Text)
	})).Return("0", nil)

	tgMock.On("Send", mock.MatchedBy(func(msg tgbotapi.MessageConfig) bool {
		return msg.ChatID == update.Message.Chat.ID && msg.Text == "0"
	})).Return(tgbotapi.Message{}, nil)

	logic.New(config.Read(), tgMock, llmMock, slog.Default()).Start(context.Background())
	tgMock.AssertExpectations(t)
	llmMock.AssertExpectations(t)
}

func TestLogic_BadUpdateChannel(t *testing.T) {
	setEnv(t)
	tgMock := tg.NewTelegram(t)
	tgMock.On("GetUpdatesChan").Return(nil, errors.New("fail"))

	llmMock := llm.NewLLM(t)

	go logic.New(config.Read(), tgMock, llmMock, slog.Default()).Start(context.Background())
	time.Sleep(10 * time.Millisecond)

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

	go logic.New(config.Read(), tgMock, llmMock, slog.Default()).Start(context.Background())
	time.Sleep(10 * time.Millisecond)

	tgMock.AssertExpectations(t)
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

	logic.New(config.Read(), tgMock, llmMock, slog.Default()).Start(context.Background())
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
	llmMock.On("GetLLMResponseAboutMsg", mock.Anything, mock.Anything).Return("", errors.New("llm fail"))

	logic.New(config.Read(), tgMock, llmMock, slog.Default()).Start(context.Background())
	tgMock.AssertExpectations(t)
	llmMock.AssertExpectations(t)
}

func TestLogic_SendError(t *testing.T) {
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
	tgMock.On("Send", mock.Anything).Return(tgbotapi.Message{}, errors.New("send fail"))

	llmMock := llm.NewLLM(t)
	llmMock.On("GetLLMResponseAboutMsg", mock.Anything, mock.Anything).Return("0", nil)

	logic.New(config.Read(), tgMock, llmMock, slog.Default()).Start(context.Background())
	tgMock.AssertExpectations(t)
	llmMock.AssertExpectations(t)
}
