package ws

import "sync"

type hub struct {
	mu          sync.RWMutex
	register    chan *client
	disconnect  chan *client
	connections map[string]*client
}
