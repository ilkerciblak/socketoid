// Package ws
// Internal package for websocket connection configuration
package ws

import (
	"bufio"
	"fmt"
	"net"
	"sync"

	"github.com/google/uuid"
)

type client struct {
	ID         string
	Connection net.Conn
	BuffRW     *bufio.ReadWriter
	Channel    chan []byte
	closeOnce  sync.Once
}

func NewClient(conn net.Conn, buffRW *bufio.ReadWriter) *client {

	return &client{
		ID:         uuid.NewString(),
		Connection: conn,
		BuffRW:     buffRW,
		Channel:    make(chan []byte, 512),
	}
}

// Client.readPump concurrently reads data frame and processes the incoming data
func (c *client) readPump(h *hub) {

	for {
		opcode, payload, err := ReadFrame(c.BuffRW)
		if err != nil {
			fmt.Printf("\n\nerr: %v", err)
			c.cleanUp(h)

			return
		}

		if opcode == opcodeClose {
			WriteCloseFrame(c.BuffRW)
			c.cleanUp(h)
			return
		}

		if opcode == opcodePing {
			if err := WritePongFrame(c.BuffRW); err != nil {
				c.cleanUp(h)
				return
			}
		}

		if opcode == opcodeUTF8Text {
			event, err := UnmarshallEvent(payload)
			if err != nil {
				WriteCloseFrame(c.BuffRW)
				c.cleanUp(h)
				return 
			}
			if err := h.router.Route(
				c,
				*event,
			); err != nil {
				errEvent := UnkownEventRespond(event.Type)
				data, e := errEvent.Marshal()
				if e != nil {
					WriteCloseFrame(c.BuffRW)
					return
				}
				WriteFrame(c.BuffRW, data)
				continue
			}

			data, _ := event.Marshal()
			WriteFrame(c.BuffRW, data)
		}

	}

}

func (c *client) writePump(h *hub) {
	for msg := range c.Channel {
		WriteFrame(c.BuffRW, msg)
	}
	c.cleanUp(h)

}

func (c *client) cleanUp(h *hub) {
	c.closeOnce.Do(func() {
		c.Connection.Close()
		close(c.Channel)
		h.disconnect <- c
	})
}

/*
[RFC 5.5.1 Close]()
If an endpoint receives a Close frame and did not previously send a
   Close frame, the endpoint MUST send a Close frame in response.  (When
   sending a Close frame in response, the endpoint typically echos the
   status code it received.)  It SHOULD do so as soon as practical.  An
   endpoint MAY delay sending a Close frame until its current message is
   sent (for instance, if the majority of a fragmented message is
   already sent, an endpoint MAY send the remaining fragments before
   sending a Close frame).  However, there is no guarantee that the
   endpoint that has already sent a Close frame will continue to process
   data.
*/
