package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for now (TODO: restrict in production)
		return true
	},
}

// HandleWebSocket handles WebSocket upgrade requests
func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket: Failed to upgrade connection: %v", err)
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan *Event, 256),
	}

	client.hub.register <- client

	// Start pumps in separate goroutines
	go client.writePump()
	go client.readPump()
}
