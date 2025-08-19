package main

import (
	"encoding/json"
	"flag"
	"fmt"

	console "github.com/auh-xda/magnesia/helpers/console"
	"github.com/nats-io/nats.go"
)

const (
	authEndpoint = "/c/6df9-4d9b-42e5-b4d5"
	version      = "0.1.0"
)

func main() {
	action := flag.String("action", "install", "Action to perform: install, update, or remove the Magnesia agent")
	auth_token := flag.String("auth_token", "", "Authentication token provided by the server")
	api_key := flag.String("api_key", "", "API key for your account")
	client_id := flag.String("client_id", "", "Unique client identifier")
	client_secret := flag.String("client_secret", "", "Client secret used for secure authentication")

	flag.Parse()

	magnesia := Magnesia{
		AuthToken:    *auth_token,
		ApiKey:       *api_key,
		ClientID:     *client_id,
		ClientSecret: *client_secret,
	}

	switch *action {
	case "install":
		magnesia.Install()

	case "intercept":
		magnesia.Intercept()

	default:
		console.Error(fmt.Sprintf("Magnesia is not aware of this action (i.e %s)", *action))
	}
}

func (ws Websocket) SendData() {
	console.Info("Establishing Connection with ws server")
	config, err := Config{}.Parse()

	if err != nil {
		console.Error("Unable to parse the config")
	}

	ws.MagnesiaUid = config.UUID
	ws.MagnesiaChannel = config.Channel
	ws.MagnesiaSiteId = config.ClientID

	nc, err := nats.Connect("nats://127.0.0.1:4222")

	if err != nil {
		console.Error("Error establishing connection with Nats")
	}

	defer nc.Drain()

	data, _ := json.Marshal(ws)

	err = nc.Publish(ws.MagnesiaChannel, data)

	if err != nil {
		console.Error("Error publishing data on WS")
	}

	console.Success(fmt.Sprintf("Sent message to NATS on subject %s", ws.MagnesiaSiteId))
}
