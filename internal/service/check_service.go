package service

import (
	"api-status-check/internal/client"
	"api-status-check/internal/model"
	"api-status-check/internal/storage"
	"api-status-check/internal/util"
	"log"
)

type CheckService struct {
	apiStore   *storage.APIStore
	checkStore *storage.CheckStore
}

func NewCheckService() *CheckService {
	return &CheckService{
		apiStore:   storage.NewAPIStore(),
		checkStore: storage.NewCheckStore(),
	}
}

func (s *CheckService) RunCheck() {
	configs, err := s.apiStore.GetAll()
	if err != nil {
		log.Printf("Failed to get API configs: %v", err)
		return
	}

	for _, cfg := range configs {
		if !cfg.Enabled {
			continue
		}

		result := client.CheckAPI(cfg.APIURL, cfg.Token, cfg.Model)

		record := model.CheckRecord{
			ID:              util.GenerateUUID(),
			APIID:           cfg.ID,
			Name:            cfg.Name,
			Tag:             cfg.Tag,
			APIURL:          cfg.APIURL,
			Model:           cfg.Model,
			Available:       result.Available,
			LatencyMs:       result.LatencyMs,
			CheckedAt:       util.Now(),
			StatusCode:      result.StatusCode,
			ErrorMessage:    result.ErrorMessage,
			ResponsePreview: result.ResponsePreview,
		}

		if err := s.checkStore.AppendHistory(record); err != nil {
			log.Printf("Failed to append history: %v", err)
		}

		if err := s.checkStore.UpdateRecent100(record); err != nil {
			log.Printf("Failed to update recent 100: %v", err)
		}

		log.Printf("Checked API %s: available=%v, latency=%dms", cfg.Name, result.Available, result.LatencyMs)
	}
}

func (s *CheckService) GetHistory() ([]model.CheckRecord, error) {
	return s.checkStore.GetHistory()
}

func (s *CheckService) GetRecent100() ([]model.CheckRecord, error) {
	return s.checkStore.GetRecent100()
}
