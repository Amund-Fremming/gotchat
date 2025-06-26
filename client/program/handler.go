package program

import (
	"client/cmd"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"

	"github.com/amund-fremming/common/enum"
	"github.com/amund-fremming/common/model"
	"github.com/gorilla/websocket"
)

var state AppState = NewAppState()

func ConnectToServer() {
	url := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/chat"}
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		slog.Error(err.Error())
	}
	defer conn.Close()

	state.Conn = conn
	log.Println("[Server] is connected")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			slog.Error(err.Error())
		}

		var message model.Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			slog.Error(err.Error())
			break
		}

		cmd.DisplayMessage(&message)
	}
}

func InputReader() {
	fmt.Println("[Client] starting input reader")

	for {
		command, err := cmd.GetCommand()
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		switch command.Action {
		case enum.Help:
			cmd.DisplayCommands()
		case enum.Exit:
			{
				fmt.Println("[Client] shutting down input reader")
				os.Exit(0)
			}
		default:
			state.Broadcast <- &command
		}
	}
}

func CommandDispatcher() {
	fmt.Println("[Client] starting command dispatcher")
	for {
		command := <-state.Broadcast
		err := state.Conn.WriteJSON(command)
		fmt.Println("[DEBUG] sending to server")
		if err != nil {
			slog.Error(err.Error())
			break
		}
	}

	fmt.Println("[Client] shutting down command dispatcher")
}
