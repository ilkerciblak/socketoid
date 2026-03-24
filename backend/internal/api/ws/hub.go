package ws

import (
	"context"
	"fmt"
	"sync"
)

type Hub struct {
	mu          sync.RWMutex
	register    chan *Client
	disconnect  chan *Client
	Broadcast   chan []byte
	connections map[string]*Client
}

func NewHub() *Hub {
	return &Hub{
		mu:          sync.RWMutex{},
		connections: make(map[string]*Client, 512),
		register:    make(chan *Client, 512),
		disconnect:  make(chan *Client, 512),
		Broadcast:   make(chan []byte, 512),
	}
}

func (h *Hub) registerClient(c *Client) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, exists := h.connections[c.ID]; exists {
		return fmt.Errorf("Client connection already exists")
	}

	h.connections[c.ID] = c

	return nil
}

func (h *Hub) disconnectClient(c *Client) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, exists := h.connections[c.ID]; !exists {
		return fmt.Errorf("Client connection does not exists")
	}

	delete(h.connections, c.ID)
	return nil

}
func (h *Hub) broadcastPayload(payload []byte) error {
	h.mu.Lock()
	channels := make([]chan []byte, 0, len(h.connections))
	for id := range h.connections {
		channels = append(channels, h.connections[id].Channel)
	}

	h.mu.Unlock()

	for _, channel := range channels {
		channel <- payload
	}
	return nil

}

func (h *Hub) Run(ctx context.Context) {
	fmt.Println("ws server running")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("ws context done")
			return
		case newClient := <-h.register:
			if err := h.registerClient(newClient); err != nil {
				fmt.Println(err)
			}

		case Client := <-h.disconnect:
			if err := h.disconnectClient(Client); err != nil {
				fmt.Println(err)
			}
		case msg := <-h.Broadcast:
			if err := h.broadcastPayload(msg); err != nil {
				fmt.Println(err)
			}
		}
	}

}
