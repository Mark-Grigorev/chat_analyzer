package settings

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type Settings struct {
	mu           sync.RWMutex
	SystemPrompt string  `json:"system_prompt"`
	ChatIDs      []int64 `json:"chat_ids"`
	Temperature  float64 `json:"temperature"`
	filePath     string  `json:"-"`
}

// Load reads from filePath. If not exists, uses defaults and creates the file.
func Load(filePath, defaultPrompt string, defaultChatIDs []int64, defaultTemp float64) (*Settings, error) {
	s := &Settings{filePath: filePath}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		// File doesn't exist: use defaults and create it
		s.SystemPrompt = defaultPrompt
		s.ChatIDs = defaultChatIDs
		s.Temperature = defaultTemp
		if err = s.save(); err != nil {
			return nil, err
		}
		return s, nil
	}

	if err = json.Unmarshal(data, s); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Settings) GetSystemPrompt() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.SystemPrompt
}

func (s *Settings) SetSystemPrompt(prompt string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.SystemPrompt = prompt
	return s.save()
}

func (s *Settings) GetChatIDs() []int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ids := make([]int64, len(s.ChatIDs))
	copy(ids, s.ChatIDs)
	return ids
}

func (s *Settings) AddChatID(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, existing := range s.ChatIDs {
		if existing == id {
			return nil
		}
	}
	s.ChatIDs = append(s.ChatIDs, id)
	return s.save()
}

func (s *Settings) RemoveChatID(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	filtered := s.ChatIDs[:0]
	for _, existing := range s.ChatIDs {
		if existing != id {
			filtered = append(filtered, existing)
		}
	}
	s.ChatIDs = filtered
	return s.save()
}

func (s *Settings) GetTemperature() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Temperature
}

func (s *Settings) SetTemperature(t float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Temperature = t
	return s.save()
}

// save must be called with the lock already held.
func (s *Settings) save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0644)
}
