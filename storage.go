package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
)

type RedisClient struct {
	conn redis.Conn
	redis.PubSubConn
}

func getMessageForIpad(hub *hub) {
	host := "127.0.0.1:6379"
	c, _ := redis.Dial("tcp", host)

	psc := redis.PubSubConn{c}
	psc.Subscribe("example")
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("REDIS SUBCRIPTION: Channel: %s; message: %s\n", v.Channel, v.Data)

			var dat map[string]interface{}
			if err := json.Unmarshal(v.Data, &dat); err != nil {
				panic(err)
			}
			ikey1 := dat["ipadKey"].(string)
			bodyMessage := dat["body"].(string)

			mes := serverMessage{body: []byte(bodyMessage), ikey:ikey1}
			h.broadcast <- mes
		case error:
			//return v
		}
	}
}