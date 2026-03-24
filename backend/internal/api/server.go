package api

import (
	"context"
	"encoding/json"
	"fmt"
	"ilkerciblak/socketoid/internal/api/sse"
	"ilkerciblak/socketoid/internal/api/ws"
	"ilkerciblak/socketoid/internal/services/board"
	"net/http"
	"time"
)

type server struct {
	server *http.Server
}

func Server(address string, idleTo int) (*server, error) {
	mux := http.NewServeMux()
	sseHub := sse.Hub()
	wsRouter := ws.NewRouter()

	wsHub := ws.NewHub()
	registerWsHandlers(
		wsRouter,
		wsHub,
		board.RegisterBoardHandlers)
	websocket := ws.WebSocket(address, wsHub, wsRouter)

	go func() {
		sseHub.Run(context.Background())
	}()
	go func() {
		wsHub.Run(context.Background())
	}()
	mux.HandleFunc("/", greet)
	mux.HandleFunc("/health", healthCheck)
	mux.HandleFunc("/events", sse.SseHandler(sseHub))
	mux.HandleFunc("/ws", websocket.Upgrade)

	return &server{
		server: &http.Server{
			Addr:        address,
			IdleTimeout: time.Second * time.Duration(idleTo),
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

func registerWsHandlers(
	router *ws.Router,
	hub *ws.Hub,
	handlerRegisterFuncs ...func(router *ws.Router, hub *ws.Hub),
) {
	for _, f := range handlerRegisterFuncs {
		f(router, hub)
	}
}
