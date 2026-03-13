// Package ws
// Internal package for websocket connection configuration
package ws

import (
	"net"

	"github.com/google/uuid"
)

type client struct {
	ID         string
	Connection net.Conn
	Channel    chan []byte
}

func NewClient(conn net.Conn) *client {

	return &client{
		ID:         uuid.NewString(),
		Connection: conn,
		Channel:    make(chan []byte, 512),
	}
}
