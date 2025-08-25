package nats

import "C"
import (
	"encoding/json"
	"fmt"

	"github.com/auh-xda/magnesia/config"
	"github.com/auh-xda/magnesia/console"
	"github.com/nats-io/nats.go"
)

const natsWsEndpoint = "nats://127.0.0.1:4222"

type Websocket struct {
	MagnesiaUid      string `json:"magnesia_uuid"`
	MagnesiaClientId string `json:"magnesia_client_id"`
	MagnesiaType     string `json:"magnesia_type"`
	MagnesiaPayload  any    `json:"magnesia_payload"`
}

func SendData(payload any, payloadType string) {
	console.Info("Establishing connection with NATS")

	config, err := config.ParseConfig()
	if err != nil {
		console.Error("Error parsing config: " + err.Error())
		return
	}

	ws := Websocket{
		MagnesiaUid:      config.UUID,
		MagnesiaClientId: config.ClientID,
		MagnesiaPayload:  payload,
		MagnesiaType:     payloadType,
	}

	nc, err := nats.Connect(natsWsEndpoint)
	if err != nil {
		console.Error("Error connecting to NATS: " + err.Error())
		return
	}
	defer nc.Close()

	data, err := json.Marshal(ws)
	if err != nil {
		console.Error("Error marshaling data: " + err.Error())
		return
	}

	subject := config.Channel
	console.Log("Publishing to subject: " + subject)
	console.Log(string(data))

	err = nc.Publish(subject, data)
	if err != nil {
		console.Error("Error publishing: " + err.Error())
		return
	}

	err = nc.Flush()
	if err != nil {
		console.Error("Error flushing: " + err.Error())
		return
	}

	console.Success(fmt.Sprintf("Sent message to NATS on subject %s", subject))
}
