package ws

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/amund-fremming/common"
	"github.com/gorilla/websocket"
)

var commandBroadcast = make(chan *common.Command)
var clients = make(map[*websocket.Conn]bool)

var connectChan = make(chan string)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Handles all connecting clients
func ClientHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to create socket.")
		conn.Close()
		return
	}
	// defer conn.Close()

	clients[conn] = true
	clients := strconv.Itoa(len(clients))
	fmt.Println("[Connected] Current clients:", clients)
	go commandHandler(conn)
}

// Handles commands for a single client
func commandHandler(conn *websocket.Conn) {
	defer func() {
		delete(clients, conn)
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			slog.Error("Failed to read message from client, closing connection")
			break
		}

		var cmd common.Command
		err = json.Unmarshal(msg, &cmd)
		if err != nil {
			slog.Error("Failed to parse command")
			break
		}

		commandBroadcast <- &cmd
	}
}

// Routes all command for the app to their handlers
func commandRouter() {
	for {
		cmd := <-commandBroadcast
		switch {
		case cmd.Action == common.Connect:
		case cmd.Action == common.Create:
		case cmd.Action == common.Exit:
		case cmd.Action == common.Send:
		}
	}
}
