package model

import "time"

type APIConfig struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Tag       string    `json:"tag"`
	APIURL    string    `json:"api_url"`
	Token     string    `json:"token"`
	Model     string    `json:"model"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
