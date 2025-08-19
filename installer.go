package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	client "github.com/auh-xda/magnesia/helpers/client"
	console "github.com/auh-xda/magnesia/helpers/console"
	"github.com/auh-xda/magnesia/services"
	"github.com/common-nighthawk/go-figure"
)

func (magnesia Magnesia) Install() {
	console.SetColor("yellow")
	myFigure := figure.NewFigure("Magnesia", "", true)
	myFigure.Print()
	console.ResetColor()

	console.Info("Installing...")
	config, error := authenticateServer(magnesia)

	if error != nil {
		console.Error("Autentication Failed")
		return
	}

	if error := createConfigFile(config); error != nil {
		console.Error(error.Error())
		return
	}

	go magnesia.Intercept()
	go magnesia.ProcessList()
	go services.GetServiceList()

	console.Info("Waiting for go routines .... ")

	time.Sleep(20 * time.Second)
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
