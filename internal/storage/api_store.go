package storage

import (
	"api-status-check/internal/model"
	"errors"
)

const APIConfigFile = "data/apis.json"

type APIStore struct {
	store *JSONStore
}

func NewAPIStore() *APIStore {
	return &APIStore{store: NewJSONStore()}
}

func (s *APIStore) GetAll() ([]model.APIConfig, error) {
	var configs []model.APIConfig
	if err := s.store.ReadJSON(APIConfigFile, &configs); err != nil {
		return nil, err
	}
	if configs == nil {
		configs = []model.APIConfig{}
	}
	return configs, nil
}

func (s *APIStore) GetByID(id string) (*model.APIConfig, error) {
	configs, err := s.GetAll()
	if err != nil {
		return nil, err
	}
	for _, cfg := range configs {
		if cfg.ID == id {
			return &cfg, nil
		}
	}
	return nil, errors.New("api config not found")
}

func (s *APIStore) Save(config model.APIConfig) error {
	configs, err := s.GetAll()
	if err != nil {
		return err
	}
	configs = append(configs, config)
	return s.store.WriteJSON(APIConfigFile, configs)
}

func (s *APIStore) Update(config model.APIConfig) error {
	configs, err := s.GetAll()
	if err != nil {
		return err
	}
	found := false
	for i, cfg := range configs {
		if cfg.ID == config.ID {
			configs[i] = config
			found = true
			break
		}
	}
	if !found {
		return errors.New("api config not found")
	}
	return s.store.WriteJSON(APIConfigFile, configs)
}

func (s *APIStore) Delete(id string) error {
	configs, err := s.GetAll()
	if err != nil {
		return err
	}
	newConfigs := []model.APIConfig{}
	for _, cfg := range configs {
		if cfg.ID != id {
			newConfigs = append(newConfigs, cfg)
		}
	}
	return s.store.WriteJSON(APIConfigFile, newConfigs)
}
