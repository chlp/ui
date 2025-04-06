package filestore

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type fileStore struct {
	path string
	mu   sync.RWMutex
}

func NewFileStore(filePath string) *fileStore {
	if filePath == "" {
		return nil
	}
	return &fileStore{path: filePath}
}

func (f *fileStore) SaveJSON(v interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return SaveJSON(f.path, v)
}

func SaveJSON(path string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (f *fileStore) LoadJSON(v interface{}) error {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return LoadJSON(f.path, v)
}

func LoadJSON(path string, v interface{}) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return os.ErrNotExist
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
