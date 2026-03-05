package api

import (
	"context"
	"encoding/json"
	"fmt"
	"ilkerciblak/socketoid/internal/api/sse"
	"net/http"
	"time"
)

type server struct {
	server *http.Server
}

func Server(address string, read_to, write_to, idle_to int) (*server, error) {
	mux := http.NewServeMux()
	sse_hub := sse.Hub()
	go func() {
		sse_hub.Run(context.TODO())
	}()
	mux.HandleFunc("/", greet)
	mux.HandleFunc("/health", healthCheck)
	mux.HandleFunc("/events", sse.SseHandler(sse_hub))

	return &server{
		server: &http.Server{
			Addr: address,
			IdleTimeout: time.Second * time.Duration(idle_to),
			Handler:     mux,
		},
	}, nil
}

func (s *server) Start(errchan chan<- error) {
	fmt.Println("server running on addr: ", s.server.Addr)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errchan <- err
		}
	}()
}

func (s *server) Shutdown(ctx context.Context) {
	fmt.Println("server is shutting down")
	c, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	s.server.Shutdown(c)
}

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	resp := struct {
		Status int    `json:"status"`
		Time   string `json:"time"`
	}{
		Status: http.StatusOK,
		Time:   time.Now().Format(time.RFC3339),
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}
