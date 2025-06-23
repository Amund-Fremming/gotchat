package ws

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

var broadcast = make(chan []byte)
var clients = make(map[*websocket.Conn]bool)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to create socket.")
	}
	defer conn.Close()

	clients[conn] = true

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			slog.Info("Client disconnected.")
			break
		}

		fmt.Println("[SERVER] Recieved message:", string(msg))

		// REMOVE THIS
		if string(msg) == "Ping" {
			msg = []byte("Pong")
		}

		broadcast <- msg
	}
}

func HandleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				delete(clients, client)
				client.Close()
				slog.Error("Lost connection to client.")
			}

			fmt.Println("[SERVER] Broadcasting messages to clients")
		}
	}
}
