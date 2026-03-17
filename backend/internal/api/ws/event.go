package ws

import (
	"encoding/json"
	"fmt"
)

type event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func (e event) Marshal() (payload []byte, err error) {
	payload, err = json.Marshal(&e)
	if err != nil {
		return nil, err
	}

	return
}

func UnmarshallEvent(payload []byte) (*event, error) {
	var message event
	if err := json.Unmarshal(payload, &message); err != nil {
		return nil, fmt.Errorf("failed event unmarshalling: %v", err)
	}

	return &message, nil
}
