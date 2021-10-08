package ws

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func GetRandomStri(l int) string {

	rand.Seed(time.Now().UnixNano())
	result := make([]byte, l)

	for i := 0; i < l; i++ {
		rand := rand.Intn(62)
		if rand < 10 {
			result[i] = byte(48 + rand)
		} else if rand < 36 {
			result[i] = byte(55 + rand)
		} else {
			result[i] = byte(61 + rand)
		}
	}
	return string(result)
}

func getSocket(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	upgrader := websocket.Upgrader{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 5 * time.Second,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf("Upgrade Error: %s", err.Error())
		return nil
	}

	socket.SetReadLimit(8192)

	return socket
}
