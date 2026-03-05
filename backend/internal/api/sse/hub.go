package sse

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type hub struct {
	mu          *sync.Mutex
	connections map[string]chan string
	register    chan *Client
	disconnect  chan *Client
	broadcast   chan string
}

func Hub() *hub {
	return &hub{
		mu:          &sync.Mutex{},
		connections: map[string]chan string{},
		register:    make(chan *Client, 512),
		disconnect:  make(chan *Client, 512),
		broadcast:   make(chan string, 512),
	}
}

type Client struct {
	ID      string
	Channel chan string
}

func NewClient() Client {
	id := uuid.New().String()
	channel := make(chan string)

	return Client{
		ID:      id,
		Channel: channel,
	}
}

func (h *hub) registerClient(ctx context.Context, client Client) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, exists := h.connections[client.ID]; exists {
		return fmt.Errorf("client connection already exists")
	}

	h.connections[client.ID] = client.Channel
	return nil
}

func (h *hub) disconnectClient(ctx context.Context, client_id string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, exists := h.connections[client_id]; !exists {
		return fmt.Errorf("client connection does not exists")
	}

	delete(h.connections, client_id)

	return nil
}

func (h *hub) broadcastMessage(ctx context.Context, message string) error {
	h.mu.Lock()
	channels := make([]chan string, 0, len(h.connections))

	for _, channel := range h.connections {
		channels = append(channels, channel)
	}

	h.mu.Unlock()

	for _, channel := range h.connections {
		channel <- message
	}

	return nil
}

func (h *hub) Run(ctx context.Context) {
	fmt.Println("sse server running")

	for {
		select {
		case new_client := <-h.register:
			if err := h.registerClient(ctx, *new_client); err != nil {
				fmt.Println(err)
			}
		case client := <-h.disconnect:
			h.disconnectClient(ctx, client.ID)
		case msg := <-h.broadcast:
			h.broadcastMessage(ctx, msg)
		}
	}
}
