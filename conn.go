package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
// Time allowed to write a message to the peer.
	WRITE_WAIT = 10 * time.Second

// Time allowed to read the next pong message from the peer.
	PONG_WAIT = 60 * time.Second

// Send pings to peer with this period. Must be less than pongWait.
	PING_PERIOD = (PONG_WAIT * 9) / 10

// Maximum message size allowed from peer.
	MAX_MESSAGE_SIZE = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	ikey string

	// The websocket connection.
	ws   *websocket.Conn

	// Buffered channel of outbound messages.
	send chan serverMessage
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(PING_PERIOD)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case serverMessage, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, serverMessage.body); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}


// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		webSocketHub.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(MAX_MESSAGE_SIZE)
	c.ws.SetReadDeadline(time.Now().Add(PONG_WAIT))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(PONG_WAIT)); return nil
	})
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		log.Printf("message: %v", message)

		webSocketHub.messageToServer <- message
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
	ikey := r.Header.Get("X-USER")

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{send: make(chan serverMessage, 256), ws: ws, ikey: ikey}
	webSocketHub.register <- c
	go c.writePump()
	c.readPump()
}


