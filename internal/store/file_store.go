package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"blockchain-bot/internal/model"
)

// --- Interfaces

type Store interface {
	Save(string) error
	GetLastData() (*model.Data, error)
}

// --- Implementations

type FileStore struct {
	fileName string
	path     string
	m        sync.RWMutex
}

func NewFileStore(fileName, path string) (*FileStore, error) {
	err := os.MkdirAll(path, 0o755)
	if err != nil {
		return nil, err
	}

	return &FileStore{
		fileName: fileName,
		path:     path,
	}, nil
}

func (s *FileStore) GetLastData() (*model.Data, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	filepath := filepath.Join(s.path, s.fileName)
	d, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var data model.Data
	if err := json.Unmarshal(d, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (s *FileStore) SetLastData(hash string, blockNumber int64) error {
	s.m.Lock()
	defer s.m.Unlock()

	filepath := filepath.Join(s.path, s.fileName)
	data := &model.Data{
		PrevHash:    hash,
		BlockNumber: blockNumber,
	}

	d, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath, d, 0o644); err != nil {
		return err
	}

	return nil
}
