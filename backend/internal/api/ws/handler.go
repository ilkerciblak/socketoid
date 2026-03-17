package ws

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

type websocket struct {
	address string
	h       *hub
}

func WebSocket(address string, h *hub) *websocket {

	return &websocket{
		address: address,
		h:       h,
	}
}

const magicString string = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

func (ws *websocket) Upgrade(w http.ResponseWriter, r *http.Request) {

	hijacker, k := w.(http.Hijacker)
	if !k {
		http.Error(w, "websocket connection is not supported", http.StatusInternalServerError)
		return
	}

	conn, buffRW, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client := NewClient(conn, buffRW)

	// Read Sec-WebSocket-Key from request headers
	key := r.Header.Get("Sec-WebSocket-Key")

	// Generate Sec-WebSocket-Accept header value
	acceptHeader := generateAcceptKey(key)
	data := handshakeResponse(acceptHeader)
	buffRW.Write(data)

	buffRW.Flush()

	ws.h.register <- client

	go client.readPump(ws.h)

	go client.writePump(ws.h)
}

func generateAcceptKey(clientKey string) string {
	// Concatenate with magic string
	clientKey = clientKey + magicString
	// Hash the result using SHA-1
	hasher := sha1.New()
	hasher.Write([]byte(clientKey))
	// Encoding and creating Sec-WebSocket-Accept Header
	hashed := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	return hashed
}

func handshakeResponse(acceptHeader string) []byte {
	lines := []string{
		fmt.Sprintf("HTTP/1.1 %d %s", http.StatusSwitchingProtocols, http.StatusText(http.StatusSwitchingProtocols)),
		fmt.Sprintf("Sec-WebSocket-Accept: %s", acceptHeader),
		"Upgrade: websocket",
		"Connection: Upgrade",
		"",
	}

	return []byte(strings.Join(lines, "\r\n"))
}
