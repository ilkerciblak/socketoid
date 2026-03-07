package sse

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
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

		t := time.NewTicker(time.Duration(1) * time.Second)
		defer t.Stop()

		for {
			select {
			case <-r.Context().Done():
				h.disconnect <- &client
				return

			case <-t.C:
				event := fmt.Sprintf(`event:"user-joined"\ndata:{"user-id":%s}\n\n`, uuid.New().String())
				_, err := fmt.Fprintf(
					w,
					event,
				)
				if err != nil {
					http.Error(w, "msg failed", http.StatusInternalServerError)
					return
				}
				flusher.Flush()

			case msg := <-client.Channel:
				_, err := fmt.Fprintf(
					w,
					"event:channel-msg\ndata:%s\n\n",
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
