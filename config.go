package main

import (
	"encoding/json"
	"os"
)

func (config Config) Parse() (Config, error) {
	jsonData, err := os.ReadFile("/magnesia/config.json")

	if err != nil {
		return config, err
	}

	if err := json.Unmarshal(jsonData, &config); err != nil {
		return config, err
	}

	return config, nil
}
