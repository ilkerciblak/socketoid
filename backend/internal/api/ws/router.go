package ws

import (
	"encoding/json"
	"fmt"
)

type HandlerFunc func(client *client, payload json.RawMessage) error

type router struct {
	routes map[string]HandlerFunc
}

func Router() *router {
	return &router{
		routes: map[string]HandlerFunc{},
	}
}

func (r *router) Register(eventType string, handler HandlerFunc) error {
	if _, exists := r.routes[eventType]; exists {
		return fmt.Errorf("event handler already registered as: %v", eventType)
	}

	r.routes[eventType] = handler

	return nil

}

func (r *router) Route(client *client, event event) error {
	handler, exists := r.routes[event.Type]
	if !exists {
		return fmt.Errorf("event handler not registered: %v", event.Type)
	}

	if err := handler(client, event.Payload); err != nil {
		return err
	}

	return nil
}
