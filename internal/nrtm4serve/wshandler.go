package nrtm4serve

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(hub *Hub) func(http.ResponseWriter, *http.Request) {
	//messageBuffer := service.NewRingBuffer[service.LogMessage](1000)
	return func(w http.ResponseWriter, r *http.Request) {
		var wg sync.WaitGroup
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error("Error upgrading", "error", err)
			return
		}
		defer conn.Close()

		client := &Client{ID: "logs", hub: hub, conn: conn, send: make(chan message, 256)}
		hub.register <- client
		wg.Add(2)
		go client.writePump(&wg)
		go client.readPump(&wg)
		wg.Wait()
	}
}
