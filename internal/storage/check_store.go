package storage

import "api-status-check/internal/model"

const (
	CheckHistoryFile = "data/check_history.json"
	Recent100File    = "data/recent_100_checks.json"
)

type CheckStore struct {
	store *JSONStore
}

func NewCheckStore() *CheckStore {
	return &CheckStore{store: NewJSONStore()}
}

func (s *CheckStore) AppendHistory(record model.CheckRecord) error {
	var history []model.CheckRecord
	if err := s.store.ReadJSON(CheckHistoryFile, &history); err != nil {
		return err
	}
	if history == nil {
		history = []model.CheckRecord{}
	}
	history = append(history, record)
	return s.store.WriteJSON(CheckHistoryFile, history)
}

func (s *CheckStore) GetHistory() ([]model.CheckRecord, error) {
	var history []model.CheckRecord
	if err := s.store.ReadJSON(CheckHistoryFile, &history); err != nil {
		return nil, err
	}
	if history == nil {
		history = []model.CheckRecord{}
	}
	return history, nil
}

func (s *CheckStore) UpdateRecent100(record model.CheckRecord) error {
	var recent []model.CheckRecord
	if err := s.store.ReadJSON(Recent100File, &recent); err != nil {
		return err
	}
	if recent == nil {
		recent = []model.CheckRecord{}
	}
	recent = append(recent, record)
	if len(recent) > 100 {
		recent = recent[len(recent)-100:]
	}
	return s.store.WriteJSON(Recent100File, recent)
}

func (s *CheckStore) GetRecent100() ([]model.CheckRecord, error) {
	var recent []model.CheckRecord
	if err := s.store.ReadJSON(Recent100File, &recent); err != nil {
		return nil, err
	}
	if recent == nil {
		recent = []model.CheckRecord{}
	}
	return recent, nil
}

func (s *CheckStore) DeleteByAPIID(apiID string) error {
	var history []model.CheckRecord
	if err := s.store.ReadJSON(CheckHistoryFile, &history); err != nil {
		return err
	}
	filtered := []model.CheckRecord{}
	for _, r := range history {
		if r.APIID != apiID {
			filtered = append(filtered, r)
		}
	}
	if err := s.store.WriteJSON(CheckHistoryFile, filtered); err != nil {
		return err
	}

	var recent []model.CheckRecord
	if err := s.store.ReadJSON(Recent100File, &recent); err != nil {
		return err
	}
	filteredRecent := []model.CheckRecord{}
	for _, r := range recent {
		if r.APIID != apiID {
			filteredRecent = append(filteredRecent, r)
		}
	}
	return s.store.WriteJSON(Recent100File, filteredRecent)
}
