package ws

import (
	"context"
	"fmt"
	"sync"
)

type hub struct {
	mu          sync.RWMutex
	register    chan *client
	disconnect  chan *client
	broadcast   chan []byte
	connections map[string]*client
}

func Hub() *hub {
	return &hub{
		mu:          sync.RWMutex{},
		connections: make(map[string]*client, 512),
		register:    make(chan *client, 512),
		disconnect:  make(chan *client, 512),
		broadcast:   make(chan []byte, 512),
	}
}

func (h *hub) registerClient(c *client) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, exists := h.connections[c.ID]; exists {
		return fmt.Errorf("client connection already exists")
	}

	h.connections[c.ID] = c

	return nil
}

func (h *hub) disconnectClient(c *client) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, exists := h.connections[c.ID]; !exists {
		return fmt.Errorf("client connection does not exists")
	}

	delete(h.connections, c.ID)
	return nil

}
func (h *hub) broadcastPayload(payload []byte) error {
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

func (h *hub) Run(ctx context.Context) {
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

		case client := <-h.disconnect:
			if err := h.disconnectClient(client); err != nil {
				fmt.Println(err)
			}
		case msg := <-h.broadcast:
			if err := h.broadcastPayload(msg); err != nil {
				fmt.Println(err)
			}
		}
	}

}
