package main

import (
	"log"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
)

// read messages from redis and send it to ipad
func getMessageFromServer(hub *hub) {
	host := REDIS_ADDR
	c, _ := redis.Dial("tcp", host)

	psc := redis.PubSubConn{c}
	psc.Subscribe(REDIS_CHANEL_TO_IPADS)
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			log.Printf("REDIS SUBCRIPTION: Channel: %s; message: %s\n", v.Channel, v.Data)

			var dat map[string]interface{}
			if err := json.Unmarshal(v.Data, &dat); err != nil {
				panic(err)
			}
			ikey := dat["ipadKey"].(string)
			bodyMessage := dat["body"].(string)

			mes := serverMessage{body: []byte(bodyMessage), ikey:ikey}
			webSocketHub.messageToClients <- mes
		case error:
		}
	}
}

//save message to redis
func sendMessageToServer(message []byte) {
	host := REDIS_ADDR
	c, _ := redis.Dial("tcp", host)
	c.Do("PUBLISH", REDIS_CHANEL_TO_SERVER, message)
}