package sse

import (
	"fmt"
	"net/http"
	"time"
)

func SseHandler(h *hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headers(w)
		flusher, k := w.(http.Flusher)
		if !k {
			http.Error(w, "sse not supported", http.StatusInternalServerError)
			return
		}

		client := NewClient()
		h.register <- &client
		payload := PresencePayload{
			Name:   r.URL.Query().Get("name"),
			UserId: client.ID,
		}
		event := UserJoinedEvent(payload)

		h.broadcast <- event.ToTextStream()

		t := time.NewTicker(time.Duration(1) * time.Second)
		defer t.Stop()

		for {
			select {
			case <-r.Context().Done():
				h.disconnect <- &client
				event := UserLeftEvent(payload)
				h.broadcast <- event.ToTextStream()

				return

			case <-t.C:
				_, err := fmt.Fprintf(
					w,
					"data: keepalive\n\n",
				)
				if err != nil {
					http.Error(w, "msg failed", http.StatusInternalServerError)
					return
				}
				flusher.Flush()

			case msg := <-client.Channel:
				_, err := fmt.Fprintf(
					w,
					msg,
				)
				if err != nil {
					http.Error(w, "msg failed", http.StatusInternalServerError)
					return
				}
				flusher.Flush()
			}
		}

	}
}

func headers(w http.ResponseWriter) {
	w.Header().Set("content-type", "text/event-stream")
	w.Header().Set("cache-control", "no-cache")
	w.Header().Set("connection", "keep-alive")
	w.Header().Set("access-control-allow-origin", "*")
}
