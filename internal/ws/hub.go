package ws

import (
	"sync"
)

type Hub struct {
	Clients    map[string]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	Mutex      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Mutex.Lock()
			h.Clients[client.IDValue] = client
			h.Mutex.Unlock()
		case client := <-h.Unregister:
			h.Mutex.Lock()
			if _, ok := h.Clients[client.IDValue]; ok {
				delete(h.Clients, client.IDValue)
				close(client.SendChan)
			}
			h.Mutex.Unlock()
		case message := <-h.Broadcast:
			h.Mutex.RLock()
			for _, client := range h.Clients {
				select {
				case client.SendChan <- message:
				default:
					close(client.SendChan)
					delete(h.Clients, client.IDValue)
				}
			}
			h.Mutex.RUnlock()
		}
	}
}
