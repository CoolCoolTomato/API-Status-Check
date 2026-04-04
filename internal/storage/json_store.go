package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type JSONStore struct {
	mu sync.RWMutex
}

func NewJSONStore() *JSONStore {
	return &JSONStore{}
}

func (s *JSONStore) ReadJSON(path string, v interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	return json.Unmarshal(data, v)
}

func (s *JSONStore) WriteJSON(path string, v interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, path)
}
