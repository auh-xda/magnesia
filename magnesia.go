package main

import (
	"encoding/json"
	"os"

	"github.com/auh-xda/magnesia/console"
)

func (magnesia Magnesia) Installed() bool {
	_, err := magnesia.ParseConfig()

	return nil == err
}

func (magnesia Magnesia) Info() {
	config, _ := magnesia.ParseConfig()

	console.Table(config)
}

func (Magnesia) ParseConfig() (Config, error) {
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
