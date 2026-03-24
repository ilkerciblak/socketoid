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

type errorRespond struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	RefEvent string `json:"ref,omitempty"`
}

func ErrorEvent(eventType, msg string) *event {
	payload := errorRespond{
		Code:     4001,
		Message:  fmt.Sprintf("%v", msg),
		RefEvent: fmt.Sprintf("%v", eventType),
	}

	rawMessage, err := json.Marshal(&payload)
	if err != nil {
		return &event{
			Type:    "error",
			Payload: json.RawMessage{},
		}

	}

	return &event{
		Type:    "error",
		Payload: rawMessage,
	}
}

func UnkownEventRespond(eventType string) *event {
	payload := errorRespond{
		Code:    1008,
		Message: fmt.Sprintf("unknown event type: %s", eventType),
	}

	rawMessage, err := json.Marshal(&payload)
	if err != nil {
		return &event{
			Type:    "error",
			Payload: json.RawMessage{},
		}
	}

	return &event{
		Type:    "error",
		Payload: rawMessage,
	}
}
