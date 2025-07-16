package ws

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"server/model"

	"github.com/amund-fremming/gotchat-common/enum"
	common "github.com/amund-fremming/gotchat-common/model"
	"github.com/gorilla/websocket"
)

var commandBroadcast = make(chan *model.ConnectionWrapper)
var state = model.NewAppState()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ClientDispatcher(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to upgrade a incomming connection")
		conn.Close()
		return
	}

	slog.Info("Client connected to the server", "address", conn.LocalAddr())
	go commandReader(conn)
}

func commandReader(conn *websocket.Conn) {
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			slog.Error("Failed to read message from client", "address", conn.LocalAddr())
			conn.Close()
			break
		}

		var cmd common.Command
		err = json.Unmarshal(msg, &cmd)
		if err != nil {
			slog.Error("Failed to unmarshal command", "address", conn.LocalAddr())
			conn.Close()
			break
		}

		wrapper := model.ConnectionWrapper{Item: cmd, Conn: conn}
		commandBroadcast <- &wrapper
	}
}

func CommandDispatcher() {
	for {
		wrapper := <-commandBroadcast
		cmd := wrapper.Item

		switch cmd.Action {
		case enum.Connect:
			handleConnect(wrapper)
		case enum.Create:
			handleCreate(wrapper)
		case enum.Send:
			handleSend(wrapper)
		case enum.Leave:
			handleLeave(wrapper)
		case enum.Exit:
			handleExit(wrapper)
		case enum.Rooms:
			handleRooms(wrapper)
		}
	}
}
