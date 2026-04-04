package model

import "time"

type CheckRecord struct {
	ID              string    `json:"id"`
	APIID           string    `json:"api_id"`
	Name            string    `json:"name"`
	Tag             string    `json:"tag"`
	APIURL          string    `json:"api_url"`
	Model           string    `json:"model"`
	Available       bool      `json:"available"`
	LatencyMs       int64     `json:"latency_ms"`
	CheckedAt       time.Time `json:"checked_at"`
	StatusCode      int       `json:"status_code"`
	ErrorMessage    string    `json:"error_message"`
	ResponsePreview string    `json:"response_preview"`
}
