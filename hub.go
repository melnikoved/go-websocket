package main

import (
	"log"
)

type hub struct {
	connections      map[string]*connection

	messageToClients chan serverMessage

	messageToServer  chan []byte

	register         chan *connection

	unregister       chan *connection
}

type serverMessage struct {
	body []byte
	ikey string
}

var webSocketHub = hub{
	messageToClients:   make(chan serverMessage),
	messageToServer: make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[string]*connection),
}

func (webSocketHub *hub) run() {
	for {
		select {
		case c := <-webSocketHub.register:
			log.Print("WEBSOCKET: Register with key ", c.ikey)
			webSocketHub.connections[c.ikey] = c
		case c := <-webSocketHub.unregister:
			log.Print("WEBSOCKET: Unregister with key ", c.ikey)
			if _, ok := webSocketHub.connections[c.ikey]; ok {
				delete(webSocketHub.connections, c.ikey)
				close(c.send)
			}
		case m := <-webSocketHub.messageToClients:
			log.Printf("WEBSOCKET OUTPUT: Message for key: %s with data: %s", m.ikey, m.body)
			if c, ok := webSocketHub.connections[m.ikey]; ok {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(webSocketHub.connections, c.ikey)
				}
			}
		case m := <-webSocketHub.messageToServer:
			log.Printf("WEBSOCKET INPUT: Message data: %s", m)
			sendMessageToServer(m);
		}
	}
}