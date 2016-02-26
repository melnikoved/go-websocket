package main

import (
	"log"
	"net/http"
)


func main() {
	//Run web socket hub
	go webSocketHub.run()

	//Get messages from redis and send it to the client
	go getMessageFromServer(&webSocketHub)

	//Listen HTTP to join clients
	http.HandleFunc("/ws", serveWs)
	err := http.ListenAndServe(WEBSOCKET_ADDR, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
