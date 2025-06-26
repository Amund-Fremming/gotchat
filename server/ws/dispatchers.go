package ws

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"server/models"

	"github.com/amund-fremming/common/enum"
	"github.com/amund-fremming/common/model"
	"github.com/gorilla/websocket"
)

var commandBroadcast = make(chan *models.ConnectionWrapper[model.Command])
var state = models.NewAppState()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ClientDispatcher(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to upgrade connection.")
		conn.Close()
		return
	}

	fmt.Println("[CLIENT] Connected")
	go commandReader(conn)
}

func commandReader(conn *websocket.Conn) {
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("[ERROR] Failed to read message from client")
			break
		}

		var cmd model.Command
		err = json.Unmarshal(msg, &cmd)
		if err != nil {
			slog.Error("[ERROR] Failed to unmarshal command")
			break
		}

		commandBroadcast <- &models.ConnectionWrapper[model.Command]{
			Item: cmd,
			Conn: conn,
		}
	}
}

func CommandDispatcher() {
	for {
		wrapper := <-commandBroadcast
		cmd := wrapper.Item

		fmt.Println("[CMD_DISPATCHER] dispatching command:", cmd.Action)

		switch cmd.Action {
		case enum.Connect:
			handleConnect(wrapper)
			break
		case enum.Create:
			handleCreate(wrapper)
			break
		case enum.Send:
			handleSend(wrapper)
			break
		case enum.Leave:
			handleLeave(wrapper)
			break
		}
	}
}
