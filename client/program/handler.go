package program

import (
	"client/cmd"
	"encoding/json"
	"fmt"
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

	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		cmd.DisplayErrorMessage(err.Error())
	}

	state.Conn = conn
	cmd.DisplayServerMessage("Connected")
}

func ServerReader() {
	for {
		_, bytes, err := state.Conn.ReadMessage()
		if err != nil {
			cmd.DisplayErrorMessage(err.Error())
		}

		failedUnmarshallingMessage := "Failed to read message from the server. Shutting down.."

		var envelope model.Envelope
		err = json.Unmarshal(bytes, &envelope)
		if err != nil {
			cmd.DisplayErrorMessage(failedUnmarshallingMessage)
			state.Conn.Close()
			break
		}

		switch envelope.Type {
		case enum.ChatMessage:
			var msg model.ChatMessage
			err := json.Unmarshal(envelope.Payload, &msg)
			if err != nil {
				cmd.DisplayErrorMessage(failedUnmarshallingMessage)
				state.Conn.Close()
				break
			}

			if state.ClientName == msg.Sender {
				break
			}

			cmd.DisplayMessage(&msg)

		case enum.ServerError:
			var error model.ServerError
			err := json.Unmarshal(envelope.Payload, &error)
			if err != nil {
				cmd.DisplayErrorMessage(failedUnmarshallingMessage)
				state.Conn.Close()
				break
			}
			state.View = error.View
			cmd.DisplayServerMessage(error.Content)

		case enum.ClientState:
			var clientState model.ClientState
			err := json.Unmarshal(envelope.Payload, &clientState)
			if err != nil {
				cmd.DisplayErrorMessage(failedUnmarshallingMessage)
				state.Conn.Close()
				break
			}
			state.Merge(&clientState)

		case enum.RoomsData:
			var data model.RoomData
			err := json.Unmarshal(envelope.Payload, &data)
			if err != nil {
				cmd.DisplayErrorMessage(failedUnmarshallingMessage)
				state.Conn.Close()
				break
			}
			fmt.Println(data.Content)
		}
	}
}

func CommandReader() {
	for {
		input := cmd.ReadInput()
		command, err := cmd.GetCommand(input, state.ClientName, state.RoomName)
		if err != nil {
			cmd.DisplayErrorMessage(err.Error())
			continue
		}

		canExecute := state.CanExecuteCommand(&command)
		if !canExecute {
			cmd.DisplayErrorMessage("Cant execute this command in current context")
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
				cmd.DisplayErrorMessage("Leave the current room before creating a new")
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
			cmd.DisplayErrorMessage(err.Error())
			break
		}
	}
}
