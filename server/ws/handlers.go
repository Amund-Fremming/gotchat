package ws

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/amund-fremming/common/enum"
	"github.com/amund-fremming/common/model"
	"github.com/gorilla/websocket"
)

var commandBroadcast = make(chan *model.Command)
var clients = make(map[*websocket.Conn]bool)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Handles all connecting clients
func ClientHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to upgrade connection.")
		conn.Close()
		return
	}
	// defer conn.Close()

	clients[conn] = true
	clients := strconv.Itoa(len(clients))
	fmt.Println("[Connected] Current clients:", clients)
	go clientCommandReader(conn)
}

// Handles commands for a single client
func clientCommandReader(conn *websocket.Conn) {
	defer func() {
		delete(clients, conn)
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("[DEBUG], ", err.Error())
			slog.Error("Failed to read message from client, closing connection")
			break
		}

		var cmd model.Command
		err = json.Unmarshal(msg, &cmd)
		if err != nil {
			slog.Error("Failed to parse command")
			break
		}

		fmt.Println("[DEBUG], recieved message from client: " + cmd.Action.String() + ":" + cmd.Arg)

		commandBroadcast <- &cmd
	}
}

// Dispatches all command to the right channels
func CommandDispatcher() {
	for {
		cmd := <-commandBroadcast
		fmt.Println("[CommandRouter] Routing command:", cmd.Action)

		switch cmd.Action {
		case enum.Connect:
		case enum.Create:
		case enum.Send:
		case enum.Leave:
		}
	}
}
