package settings_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Mark-Grigorev/chat_analyzer/internal/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func load(t *testing.T, chatIDs []int64) (*settings.Settings, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "settings.json")
	s, err := settings.Load(path, "default prompt", chatIDs, 0.5)
	require.NoError(t, err)
	return s, path
}

func TestLoad_NewFile_UsesDefaults(t *testing.T) {
	s, path := load(t, []int64{1, 2})

	assert.Equal(t, "default prompt", s.GetSystemPrompt())
	assert.Equal(t, []int64{1, 2}, s.GetChatIDs())
	assert.Equal(t, 0.5, s.GetTemperature())

	// файл должен быть создан на диске
	_, err := os.Stat(path)
	assert.NoError(t, err)
}

func TestLoad_ExistingFile_ReadsValues(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")

	data, _ := json.Marshal(map[string]any{
		"system_prompt": "loaded prompt",
		"chat_ids":      []int64{10, 20},
		"temperature":   0.7,
	})
	require.NoError(t, os.WriteFile(path, data, 0644))

	s, err := settings.Load(path, "ignored", nil, 0.0)
	require.NoError(t, err)

	assert.Equal(t, "loaded prompt", s.GetSystemPrompt())
	assert.Equal(t, []int64{10, 20}, s.GetChatIDs())
	assert.InDelta(t, 0.7, s.GetTemperature(), 1e-9)
}

func TestLoad_InvalidJSON_ReturnsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	require.NoError(t, os.WriteFile(path, []byte("not json"), 0644))

	_, err := settings.Load(path, "p", nil, 0)
	assert.Error(t, err)
}

func TestSetSystemPrompt_PersistsToDisk(t *testing.T) {
	s, path := load(t, nil)

	require.NoError(t, s.SetSystemPrompt("new prompt"))
	assert.Equal(t, "new prompt", s.GetSystemPrompt())

	// перечитываем с диска
	s2, err := settings.Load(path, "", nil, 0)
	require.NoError(t, err)
	assert.Equal(t, "new prompt", s2.GetSystemPrompt())
}

func TestAddChatID_AddsAndPersists(t *testing.T) {
	s, path := load(t, []int64{1})

	require.NoError(t, s.AddChatID(2))
	assert.Equal(t, []int64{1, 2}, s.GetChatIDs())

	s2, err := settings.Load(path, "", nil, 0)
	require.NoError(t, err)
	assert.Equal(t, []int64{1, 2}, s2.GetChatIDs())
}

func TestAddChatID_IgnoresDuplicate(t *testing.T) {
	s, _ := load(t, []int64{1})

	require.NoError(t, s.AddChatID(1))
	assert.Equal(t, []int64{1}, s.GetChatIDs())
}

func TestRemoveChatID_RemovesAndPersists(t *testing.T) {
	s, path := load(t, []int64{1, 2, 3})

	require.NoError(t, s.RemoveChatID(2))
	assert.Equal(t, []int64{1, 3}, s.GetChatIDs())

	s2, err := settings.Load(path, "", nil, 0)
	require.NoError(t, err)
	assert.Equal(t, []int64{1, 3}, s2.GetChatIDs())
}

func TestRemoveChatID_NotExisting_NoError(t *testing.T) {
	s, _ := load(t, []int64{1})

	require.NoError(t, s.RemoveChatID(999))
	assert.Equal(t, []int64{1}, s.GetChatIDs())
}

func TestSetTemperature_PersistsToDisk(t *testing.T) {
	s, path := load(t, nil)

	require.NoError(t, s.SetTemperature(0.9))
	assert.InDelta(t, 0.9, s.GetTemperature(), 1e-9)

	s2, err := settings.Load(path, "", nil, 0)
	require.NoError(t, err)
	assert.InDelta(t, 0.9, s2.GetTemperature(), 1e-9)
}

func TestGetChatIDs_ReturnsCopy(t *testing.T) {
	s, _ := load(t, []int64{1, 2})

	ids := s.GetChatIDs()
	ids[0] = 999

	assert.Equal(t, []int64{1, 2}, s.GetChatIDs())
}
