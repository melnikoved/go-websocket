WEB SOCKET server for project agenda
------------------------------------
Application works with 2 business way:
1. Read messages from REDIS and push it to ipads through web socket connection
2. Receive messages from ipad through web socket connection and save it to redis


INSTALL
=======
1. git clone http://rhodecode.shire.local/tatneft/agenda-web-socket-server ws_agenda
2. go get github.com/garyburd/redigo/redis
3. go get github.com/gorilla/websocket
4. go install {FOLDER}

