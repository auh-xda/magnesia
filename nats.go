package main

import (
	"encoding/json"
	"fmt"

	console "github.com/auh-xda/magnesia/helpers/console"
	"github.com/nats-io/nats.go"
)

func (ws Websocket) SendData() {
	return
	console.Info("Establishing Connection with ws server")
	config, err := Config{}.Parse()
	ws.MagnesiaUid = config.UUID
	ws.MagnesiaChannel = config.Channel
	ws.MagnesiaSiteId = config.ClientID

	nc, err := nats.Connect(natsWsEndpoint)

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
