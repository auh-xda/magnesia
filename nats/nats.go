package nats

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/auh-xda/magnesia/config"
	"github.com/auh-xda/magnesia/console"
	"github.com/nats-io/nats.go"
)

const natsWsEndpoint = "nats://192.168.3.53:4222"

type Websocket struct {
	MagnesiaUid      string `json:"magnesia_uuid"`
	MagnesiaClientId string `json:"magnesia_client_id"`
	MagnesiaType     string `json:"magnesia_type"`
	MagnesiaPayload  any    `json:"magnesia_payload"`
}

func SendData(payload any, payloadType string) {
	console.Info("Establishing connection with NATS")

	cfg, err := config.ParseConfig()
	if err != nil {
		console.Error("Error parsing config: " + err.Error())
		return
	}

	// compare with state & get only changed values
	// changed, err := GetChangedValue(payload)
	// if err != nil {
	// 	console.Error("Error getting changed values: " + err.Error())
	// 	return
	// }

	// // skip publishing if nothing changed
	// if changed == nil {
	// 	console.Info("No changes detected, skipping publish")
	// 	return
	// }

	ws := Websocket{
		MagnesiaUid:      cfg.UUID,
		MagnesiaClientId: cfg.ClientID,
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

	subject := cfg.Channel
	console.Info("Publishing to subject: " + subject)

	if err := nc.Publish(subject, data); err != nil {
		console.Error("Error publishing: " + err.Error())
		return
	}

	if err := nc.Flush(); err != nil {
		console.Error("Error flushing: " + err.Error())
		return
	}

	console.Success(fmt.Sprintf("Sent message to NATS on subject %s", subject))

	time.Sleep(4 * time.Second)
}

// ---------------- State Handling ----------------

func GetChangedValue(payload any) (any, error) {
	state, _ := LoadStateData() // ignore if missing

	changed, err := DeepDiff(state, payload)
	if err != nil {
		return nil, err
	}

	// Always save new payload as the latest state
	if err := SaveStateData(payload); err != nil {
		return nil, err
	}

	return changed, nil
}

func LoadStateData() (any, error) {
	c := config.State()
	content, err := os.ReadFile(c)
	if err != nil {
		// no file yet
		return nil, err
	}

	var state any
	if err := json.Unmarshal(content, &state); err != nil {
		return nil, fmt.Errorf("invalid JSON in state file: %w", err)
	}
	return state, nil
}

func SaveStateData(payload any) error {
	c := config.State()

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(c, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// ---------------- Diff Logic ----------------

func DeepDiff(oldData, newData any) (any, error) {
	oldNorm, err := normalize(oldData)
	if err != nil {
		return nil, err
	}
	newNorm, err := normalize(newData)
	if err != nil {
		return nil, err
	}
	return diffRecursive(oldNorm, newNorm), nil
}

func normalize(v any) (any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var out any
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func diffRecursive(old, new any) any {
	switch oldTyped := old.(type) {
	case map[string]any:
		newTyped, ok := new.(map[string]any)
		if !ok {
			return new
		}
		changed := make(map[string]any)
		for k, newVal := range newTyped {
			if oldVal, exists := oldTyped[k]; !exists {
				changed[k] = newVal
			} else {
				if d := diffRecursive(oldVal, newVal); d != nil {
					changed[k] = d
				}
			}
		}
		if len(changed) == 0 {
			return nil
		}
		return changed

	case []any:
		newTyped, ok := new.([]any)
		if !ok || !reflect.DeepEqual(oldTyped, newTyped) {
			return new
		}
		return nil

	default:
		if !reflect.DeepEqual(old, new) {
			return new
		}
		return nil
	}
}
