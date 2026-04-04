package service

import (
	"api-status-check/internal/model"
	"api-status-check/internal/storage"
	"api-status-check/internal/util"
)

type APIService struct {
	store      *storage.APIStore
	checkStore *storage.CheckStore
}

func NewAPIService() *APIService {
	return &APIService{store: storage.NewAPIStore(), checkStore: storage.NewCheckStore()}
}

func (s *APIService) Create(name, tag, apiURL, token, modelName string, enabled bool) (*model.APIConfig, error) {
	config := model.APIConfig{
		ID:        util.GenerateUUID(),
		Name:      name,
		Tag:       tag,
		APIURL:    apiURL,
		Token:     token,
		Model:     modelName,
		Enabled:   enabled,
		CreatedAt: util.Now(),
		UpdatedAt: util.Now(),
	}
	if err := s.store.Save(config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (s *APIService) GetAll() ([]model.APIConfig, error) {
	return s.store.GetAll()
}

func (s *APIService) GetByID(id string) (*model.APIConfig, error) {
	return s.store.GetByID(id)
}

func (s *APIService) Update(id, name, tag, apiURL, token, modelName string, enabled *bool) error {
	config, err := s.store.GetByID(id)
	if err != nil {
		return err
	}
	if name != "" {
		config.Name = name
	}
	if tag != "" {
		config.Tag = tag
	}
	if apiURL != "" {
		config.APIURL = apiURL
	}
	if token != "" {
		config.Token = token
	}
	if modelName != "" {
		config.Model = modelName
	}
	if enabled != nil {
		config.Enabled = *enabled
	}
	config.UpdatedAt = util.Now()
	return s.store.Update(*config)
}

func (s *APIService) Delete(id string) error {
	if err := s.store.Delete(id); err != nil {
		return err
	}
	return s.checkStore.DeleteByAPIID(id)
}
