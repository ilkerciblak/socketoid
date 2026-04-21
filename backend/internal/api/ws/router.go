package ws

import (
	"encoding/json"
	"fmt"
)

type HandlerFunc func(client *Client, payload json.RawMessage) error

type Router struct {
	routes map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes: map[string]HandlerFunc{},
	}
}

func (r *Router) Register(eventType string, handler HandlerFunc) error {
	if _, exists := r.routes[eventType]; exists {
		return fmt.Errorf("event handler already registered as: %v", eventType)
	}

	r.routes[eventType] = handler

	return nil

}

func (r *Router) Route(client *Client, event Event) error {
	handler, exists := r.routes[event.Type]
	if !exists {
		// client<-unknown event
		return fmt.Errorf("event handler not registered: %v", event.Type)
	}

	if err := handler(client, event.Payload); err != nil {
		return err
	}

	return nil
}
