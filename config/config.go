package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Version  string `json:"version"`
	UUID     string `json:"uuid"`
	Momentum string `json:"server"`
	Interval string `json:"interval"`
	Channel  string `json:"channel"`
	ClientID string `json:"client_id"`
}

func ParseConfig() (Config, error) {
	jsonData, err := os.ReadFile("/magnesia/config.json")

	var config Config

	if err != nil {
		return config, err
	}

	if err := json.Unmarshal(jsonData, &config); err != nil {
		return config, err
	}

	return config, nil
}
