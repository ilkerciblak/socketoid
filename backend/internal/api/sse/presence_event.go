package sse

import (
	"encoding/json"
	"fmt"
)

type PresencePayload struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
}

type SSEEvent struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

const (
	USER_JOINED string = "user.joined"
	USER_LEFT   string = "user.left"
)

func UserJoinedEvent(payload PresencePayload) SSEEvent {
	return SSEEvent{
		Type:    USER_JOINED,
		Payload: payload,
	}
}

func UserLeftEvent(payload PresencePayload) SSEEvent {
	return SSEEvent{
		Type:    USER_LEFT,
		Payload: payload,
	}
}

func (e SSEEvent) ToTextStream() string {
	dataByte, err := json.Marshal(e.Payload)
	if err != nil {
		return fmt.Sprintf("event: error\ndata: %s", err.Error())

	}
	data := fmt.Sprintf("data: %s", dataByte)
	event := fmt.Sprintf("event: %s", e.Type)
	return event + "\n" + data + "\n\n"

}
