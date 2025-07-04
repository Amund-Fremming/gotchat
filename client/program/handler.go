package program

import (
	"client/cmd"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/amund-fremming/common/config"
	"github.com/amund-fremming/common/enum"
	"github.com/amund-fremming/common/model"
	"github.com/gorilla/websocket"
)

var state AppState = NewAppState()

func ConnectToServer(cfg *config.Config) {
	url := url.URL{
		Scheme: cfg.SocketScheme,
		Host:   cfg.URL + ":" + cfg.Port,
		Path:   "/chat",
	}

	slog.Debug("Url constructed", "scheme", url.Scheme, "host", url.Host, "path", url.Path)

	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	state.Conn = conn
	fmt.Println("[SERVER] Connected")
}

func ServerReader() {
	for {
		_, bytes, err := state.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err.Error())
		}

		failedUnmarshallingMessage := "Failed to read message from the server. Shutting down.."

		var envelope model.Envelope
		err = json.Unmarshal(bytes, &envelope)
		if err != nil {
			fmt.Println("[ERROR]", failedUnmarshallingMessage)
			state.Conn.Close()
			break
		}

		switch envelope.Type {
		case enum.ChatMessage:
			var msg model.ChatMessage
			err := json.Unmarshal(envelope.Payload, &msg)
			if err != nil {
				fmt.Println("[ERROR]", failedUnmarshallingMessage)
				state.Conn.Close()
				break
			}
			cmd.DisplayMessage(&msg)

		case enum.ServerError:
			var error model.ServerError
			err := json.Unmarshal(envelope.Payload, &error)
			if err != nil {
				fmt.Println("[ERROR]", failedUnmarshallingMessage)
				state.Conn.Close()
				break
			}
			state.View = error.View
			cmd.DisplayError(error.Content)

		case enum.ClientState:
			var clientState model.ClientState
			err := json.Unmarshal(envelope.Payload, &clientState)
			if err != nil {
				fmt.Println("[ERROR]", failedUnmarshallingMessage)
				state.Conn.Close()
				break
			}
			state.Merge(&clientState)

		case enum.RoomsData:
			var data model.RoomData
			err := json.Unmarshal(envelope.Payload, &data)
			if err != nil {
				fmt.Println("[ERROR]", failedUnmarshallingMessage)
				state.Conn.Close()
				break
			}
			fmt.Println(data.Content)
		}
	}
}

// TODO: This is straight ugly, fix it
func CommandReader() {
	for {
		command, err := cmd.GetCommand(state.ClientName, state.RoomName)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		canExecute := state.CanExecuteCommand(&command)
		if !canExecute {
			fmt.Println("[ERROR] Cant execute this command in current context")
			continue
		}

		switch command.Action {
		case enum.Help:
			cmd.DisplayCommands()

		case enum.Leave:
			state.View = enum.Lobby

		case enum.Exit:
			state.Conn.Close()
			os.Exit(0)

		case enum.Connect, enum.Create:
			if state.IsConnected() {
				fmt.Println("[ERROR] Leave the current room before creating a new")
				continue
			}

			state.ClientName = command.ClientName
			state.RoomName = command.RoomName
			state.View = enum.Room
		}

		state.Broadcast <- &command
	}
}

func CommandDispatcher() {
	for {
		command := <-state.Broadcast
		err := state.Conn.WriteJSON(command)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
	}
}
