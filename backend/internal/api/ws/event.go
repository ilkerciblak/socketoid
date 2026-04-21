package ws

import (
	"encoding/json"
	"fmt"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func (e Event) Marshal() (payload []byte, err error) {
	payload, err = json.Marshal(&e)
	if err != nil {
		return nil, err
	}

	return
}

func UnmarshallEvent(payload []byte) (*Event, error) {
	var message Event
	if err := json.Unmarshal(payload, &message); err != nil {
		return nil, fmt.Errorf("failed Event unmarshalling: %v", err)
	}

	return &message, nil
}

type errorRespond struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	RefEvent string `json:"ref,omitempty"`
}

func ErrorEvent(eventType, msg string) *Event {
	payload := errorRespond{
		Code:     4001,
		Message:  fmt.Sprintf("%v", msg),
		RefEvent: fmt.Sprintf("%v", eventType),
	}

	rawMessage, err := json.Marshal(&payload)
	if err != nil {
		return &Event{
			Type:    "error",
			Payload: json.RawMessage{},
		}

	}

	return &Event{
		Type:    "error",
		Payload: rawMessage,
	}
}

func UnkownEventRespond(eventType string) *Event {
	payload := errorRespond{
		Code:    1008,
		Message: fmt.Sprintf("unknown Event type: %s", eventType),
	}

	rawMessage, err := json.Marshal(&payload)
	if err != nil {
		return &Event{
			Type:    "error",
			Payload: json.RawMessage{},
		}
	}

	return &Event{
		Type:    "error",
		Payload: rawMessage,
	}
}

const (
	UserJoinedChannel  string = "user.joined"
)
