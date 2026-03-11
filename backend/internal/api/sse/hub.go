package sse

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type hub struct {
	mu          *sync.Mutex
	connections map[string]Client
	register    chan *Client
	disconnect  chan *Client
	broadcast   chan string
}

func Hub() *hub {
	return &hub{
		mu:          &sync.Mutex{},
		connections: make(map[string]Client, 512),
		register:    make(chan *Client, 512),
		disconnect:  make(chan *Client, 512),
		broadcast:   make(chan string, 512),
	}
}

type Client struct {
	ID       string      `json:"id"`
	Username string      `json:"user_name"`
	Channel  chan string `json:""`
}

func NewClient(name string) Client {
	id := uuid.New().String()
	channel := make(chan string, 1)

	return Client{
		ID:       id,
		Username: name,
		Channel:  channel,
	}
}

func (h *hub) registerClient(client Client) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, exists := h.connections[client.ID]; exists {
		return fmt.Errorf("client connection already exists")
	}

	h.connections[client.ID] = client
	return nil
}

func (h *hub) disconnectClient(client_id string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, exists := h.connections[client_id]; !exists {
		return fmt.Errorf("client connection does not exists")
	}
	clientCh := h.connections[client_id].Channel
	close(clientCh)
	delete(h.connections, client_id)

	return nil
}

func (h *hub) broadcastMessage(message string) error {
	h.mu.Lock()
	channels := make([]chan string, 0, len(h.connections))

	for _, client := range h.connections {
		channels = append(channels, client.Channel)
	}

	h.mu.Unlock()

	for _, channel := range channels {
		channel <- message
	}

	return nil
}

func (h *hub) Run(ctx context.Context) {
	fmt.Println("sse server running")

	for {
		select {
		case <-ctx.Done():
			return
		case new_client := <-h.register:
			if err := h.registerClient(*new_client); err != nil {
				fmt.Println(err)
			}
			payload := PresencePayload{
				Name:   new_client.Username,
				UserId: new_client.ID,
			}

			event := UserJoinedEvent(payload)
			new_client.Channel <- h.ConnectionListEvent(new_client.ID)

			h.broadcastMessage(event.ToTextStream())

		case client := <-h.disconnect:
			h.disconnectClient(client.ID)
			payload := PresencePayload{
				Name:   client.Username,
				UserId: client.ID,
			}
			event := UserLeftEvent(payload)
			h.broadcastMessage(event.ToTextStream())

		case msg := <-h.broadcast:
			h.broadcastMessage(msg)
		}
	}
}

func (h *hub) ConnectionListEvent(except string) string {
	h.mu.Lock()
	var userList []PresencePayload
	for _, client := range h.connections {
		if !strings.EqualFold(except, client.ID) {
			userList = append(userList, PresencePayload{
				UserId: client.ID,
				Name:   client.Username,
			})
		}
	}
	h.mu.Unlock()

	dataByte, _ := json.Marshal(userList)
	return fmt.Sprintf("event: presence.init\ndata: %s\n\n", string(dataByte))

}
