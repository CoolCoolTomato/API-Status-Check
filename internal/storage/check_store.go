package storage

import (
	"api-status-check/internal/model"
	"os"
	"path/filepath"
	"sort"
)

const ChecksBaseDir = "data/checks"

type CheckStore struct {
	store *JSONStore
}

func NewCheckStore() *CheckStore {
	return &CheckStore{store: NewJSONStore()}
}

func channelDir(apiID string) string {
	return filepath.Join(ChecksBaseDir, apiID)
}

func historyFile(apiID string) string {
	return filepath.Join(channelDir(apiID), "history.json")
}

func recent100File(apiID string) string {
	return filepath.Join(channelDir(apiID), "recent_100.json")
}

func (s *CheckStore) AppendHistory(record model.CheckRecord) error {
	hFile := historyFile(record.APIID)
	var history []model.CheckRecord
	if err := s.store.ReadJSON(hFile, &history); err != nil {
		return err
	}
	history = append(history, record)
	return s.store.WriteJSON(hFile, history)
}

func (s *CheckStore) UpdateRecent100(record model.CheckRecord) error {
	rFile := recent100File(record.APIID)
	var recent []model.CheckRecord
	if err := s.store.ReadJSON(rFile, &recent); err != nil {
		return err
	}
	recent = append(recent, record)
	if len(recent) > 100 {
		recent = recent[len(recent)-100:]
	}
	return s.store.WriteJSON(rFile, recent)
}

func (s *CheckStore) GetHistory() ([]model.CheckRecord, error) {
	return s.readAllChannels("history.json")
}

func (s *CheckStore) GetRecent100() ([]model.CheckRecord, error) {
	return s.readAllChannels("recent_100.json")
}

func (s *CheckStore) readAllChannels(filename string) ([]model.CheckRecord, error) {
	entries, err := os.ReadDir(ChecksBaseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []model.CheckRecord{}, nil
		}
		return nil, err
	}
	var all []model.CheckRecord
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		var records []model.CheckRecord
		if err := s.store.ReadJSON(filepath.Join(ChecksBaseDir, e.Name(), filename), &records); err != nil {
			return nil, err
		}
		all = append(all, records...)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].CheckedAt.Before(all[j].CheckedAt)
	})
	if all == nil {
		all = []model.CheckRecord{}
	}
	return all, nil
}

func (s *CheckStore) DeleteByAPIID(apiID string) error {
	return os.RemoveAll(channelDir(apiID))
}
