# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run all tests with race detector and coverage
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Run a single test
go test -v -run TestName ./internal/logic/

# Lint
golangci-lint run --timeout=5m

# Build
go build -o app cmd/*.go

# Docker build
docker build -f docker/Dockerfile --build-arg VERSION=1.0.0 -t chat_analyzer .
```

## Architecture

**chat_analyzer** — Telegram bot that monitors group chats for spam/scam messages using an OpenAI-compatible LLM. Detected spam is auto-deleted.

### Component Flow

```
main() → config.Read() → llm.New() → telegram.New() → settings.Load() → logic.Start()
```

`logic.Start()` runs the main event loop:
- **Private messages** → admin command handler (`/setprompt`, `/addchat`, `/removechat`, `/listchats`, `/settemperature`, `/status`, `/help`)
- **Group messages** (in configured chat IDs) → `llm.GetLLMResponseAboutMsg()` → if `"1"` → `tgBot.DeleteMessage()`

### Key Layers

- **`internal/config/`** — Env-based config via `envconfig`. Required vars: `DEBUG`, `LLM_URL`, `LLM_TOKEN`, `LLM_MODEL`, `LLM_TEMPERATURE`, `TG_TOKEN`, `TG_CHAT_IDS`, `TG_ADMIN_USER_ID`. Optional: `SETTINGS_PATH` (default: `settings.json`).
- **`internal/settings/`** — Thread-safe JSON file that persists runtime state: system prompt, allowed chat IDs, LLM temperature. Changes via admin commands are written here immediately.
- **`internal/clients/llm/`** — Wraps `langchaingo` (OpenAI-compatible). Sends each message with the system prompt and expects `"1"` (spam) or `"0"` (clean).
- **`internal/clients/telegram/`** — Wraps `go-telegram-bot-api`. Interface: `GetUpdatesChan()`, `Send()`, `DeleteMessage()`.
- **`internal/logic/`** — Orchestrates everything. Holds references to all clients and settings.

### Interfaces & Mocks

Both `LLM` and `Telegram` are interfaces with mocks in their `mocks/` subdirectory (generated with mockery). Use these in tests — see `internal/logic/logic_test.go` for examples.

### Exit Codes

| Code | Meaning |
|------|---------|
| 1 | Config read error |
| 2 | LLM client creation error |
| 3 | Telegram client creation error |
| 4 | Error reading from update channel |
| 5 | Settings load error |
