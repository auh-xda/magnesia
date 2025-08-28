package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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
	path := Path()

	jsonData, err := os.ReadFile(path)

	var config Config

	if err != nil {
		return config, err
	}

	if err := json.Unmarshal(jsonData, &config); err != nil {
		return config, err
	}

	return config, nil
}

func State() string {
	dir := Dir()
	switch runtime.GOOS {
	case "windows":
		// Prefer ProgramData
		programData := os.Getenv("ProgramData")
		if programData == "" {
			programData = `C:\ProgramData`
		}
		return filepath.Join(dir, "state.json")

	default:
		// Linux / Unix
		return fmt.Sprintf("%s/state.json", dir)
	}
}

func Path() string {
	dir := Dir()

	switch runtime.GOOS {
	case "windows":
		// Prefer ProgramData
		programData := os.Getenv("ProgramData")
		if programData == "" {
			programData = `C:\ProgramData`
		}
		return filepath.Join(dir, "config.json")

	default:
		// Linux / Unix
		return fmt.Sprintf("%s/config.json", dir)
	}
}

func Dir() string {
	switch runtime.GOOS {
	case "windows":
		// Prefer ProgramData
		programData := os.Getenv("ProgramData")
		if programData == "" {
			programData = `C:\ProgramData`
		}
		return filepath.Join(programData, "Magnesia")

	case "darwin":
		// macOS
		return "/Library/Application Support/Magnesia"

	default:
		// Linux / Unix
		return "/etc/magnesia"
	}
}
