package main

import (
	"log"

)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	connections map[string]*connection

	// Inbound messages from the connections.
	broadcast chan serverMessage

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

type serverMessage struct {
	body []byte
	ikey string
}

var h = hub{
	broadcast:   make(chan serverMessage),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[string]*connection),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			log.Print("WEBSOCKET: Register with key ", c.ikey)
			h.connections[c.ikey] = c
		case c := <-h.unregister:
			log.Print("WEBSOCKET: Unregister with key ", c.ikey)
			if _, ok := h.connections[c.ikey]; ok {
				delete(h.connections, c.ikey)
				close(c.send)
			}
		case m := <-h.broadcast:
			log.Printf("WEBSOCKET: Message for key: %s with data: %s", m.ikey, m.body)
			if c, ok := h.connections[m.ikey]; ok {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(h.connections, c.ikey)
				}
			}

		}
	}
}