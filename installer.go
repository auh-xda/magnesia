package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/auh-xda/magnesia/client"
	"github.com/auh-xda/magnesia/console"
	"github.com/auh-xda/magnesia/interceptor"
	"github.com/common-nighthawk/go-figure"
)

func (magnesia Magnesia) Install() {
	console.SetColor("yellow")
	myFigure := figure.NewFigure("Magnesia", "", true)
	myFigure.Print()
	console.ResetColor()

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

	go magnesia.Intercept()
	go magnesia.ProcessList()
	go interceptor.GetServices()

	console.Info("Waiting for go routines .... ")

	time.Sleep(10 * time.Second)
}

func createConfigFile(config Config) error {
	configDir := "/magnesia"

	console.Info("Generating Magnesia configurations")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	configFile := filepath.Join(configDir, "config.json")

	jsonData, err := json.MarshalIndent(config, "", "  ")

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
		return config, err
	}

	var Auth AuthResponse

	err = json.Unmarshal(response.Body(), &Auth)

	if err != nil {
		return config, fmt.Errorf("error unmarshaling JSON: %s", err)
	}

	if !Auth.Success {
		return config, fmt.Errorf("%s", Auth.Message)
	}

	console.Success("Authentication Sucessfull")

	return Auth.Config, nil
}
