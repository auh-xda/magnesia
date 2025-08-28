package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/auh-xda/magnesia/client"
	"github.com/auh-xda/magnesia/config"
	"github.com/auh-xda/magnesia/console"
	"github.com/auh-xda/magnesia/installer"
	"github.com/auh-xda/magnesia/interceptor"
	"github.com/common-nighthawk/go-figure"
)

func (magnesia Magnesia) Install() {
	console.SetColor("yellow")
	myFigure := figure.NewFigure("Magnesia", "", true)
	myFigure.Print()
	console.ResetColor()

	if magnesia.Installed() {
		console.Warn("Removing existing installation")
		installer.Uninstall()
	}

	console.Info("Installing...")
	config, err := authenticateServer(magnesia)

	if err != nil {
		console.Error("Authentication Failed")
		return
	}

	if err := createConfigFile(config); err != nil {
		console.Error(err.Error())
		return
	}

	err = installer.CreateService()

	if err != nil {
		console.Error(err.Error())
		return
	}

	magnesia.Intercept()
	magnesia.ProcessList()
	interceptor.GetServices()
	interceptor.InstalledSoftwareList()
}

func createConfigFile(cfg Config) error {
	configDir := config.Dir()

	console.Info("Generating Magnesia configurations")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	configFile := filepath.Join(configDir, "config.json")

	jsonData, err := json.MarshalIndent(cfg, "", "  ")

	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	console.Success("Config file updated")

	return nil
}

func authenticateServer(Agent Magnesia) (Config, error) {

	config := Config{}

	authPayload := AuthRequest{
		AuthToken:    Agent.AuthToken,
		ClientID:     Agent.ClientID,
		ClientSecret: Agent.ClientSecret,
		ApiKey:       Agent.ApiKey,
	}

	console.Info("Authenticating with server...")

	response, err := client.Post(authEndpoint, authPayload)

	if err != nil {
		console.Error("Some error occured while authenticating")
		return config, err
	}

	var Auth AuthResponse

	err = json.Unmarshal(response.Body(), &Auth)

	if err != nil {
		console.Error("Unmarshal error")
		return config, fmt.Errorf("error unmarshaling JSON: %s", err)
	}

	if !Auth.Success {
		console.Error(Auth.Message)
		return config, fmt.Errorf("%s", Auth.Message)
	}

	console.Success(Auth.Message)

	return Auth.Config, nil
}
