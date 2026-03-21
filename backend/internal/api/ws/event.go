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

const (
	EventCardCreated = "board.card.created"
	EventCardMoved   = "board.card.moved"
	EventCardDeleted = "board.card.deleted"
	EventCardUpdated = "board.card.updated"
)

type CardCreatedPayload struct {
	CardID string `json:"card_id"`
	Title  string `json:"title"`
	Column string `json:"column"`
}

type CardMovedPayload struct {
	CardID string `json:"card_id"`
	Column string `json:"column"`
}

type CardDeletedPayload struct {
	CardID string `json:"card_id"`
}

type CardUpdatePayload struct {
	CardID string `json:"card_id"`
	Title  string `json:"title"`
	Column string `json:"column"`
}

type errorRespond struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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
