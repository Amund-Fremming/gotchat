package program

import (
	"client/cmd"
	"encoding/json"
	"fmt"
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

		var envelope model.Envelope
		err = json.Unmarshal(bytes, &envelope)
		if err != nil {
			fmt.Println("[DEBUG]", err.Error())
			fmt.Println("[ERROR] Failed to unmarshal bytes (100)")
			break
		}

		switch envelope.Type {
		case enum.ChatMessage:
			var msg model.ChatMessage
			err := json.Unmarshal(envelope.Payload, &msg)
			if err != nil {
				fmt.Println("[ERROR] Failed to unmarshal bytes (101)")
				break
			}
			cmd.DisplayMessage(&msg)

		case enum.ServerError:
			var error model.ServerError
			err := json.Unmarshal(envelope.Payload, &error)
			if err != nil {
				fmt.Println("[ERROR] Failed to unmarshal bytes (102)")
				break
			}
			state.View = error.View
			cmd.DisplayError(error.Content)

		case enum.ClientState:
			var clientState model.ClientState
			err := json.Unmarshal(envelope.Payload, &clientState)
			if err != nil {
				fmt.Println("[ERROR] Failed to unmarshal bytes (103)")
				break
			}
			state.Merge(&clientState)
		}
	}
}

func CommandReader() {
	fmt.Println("[DEBUG] Starting input reader")

	for {
		command, err := cmd.GetCommand()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		switch command.Action {
		case enum.Help:
			cmd.DisplayCommands()

		case enum.Exit:
			// Handle leaving correct so server does not lag with conneciton open
			fmt.Println("[DEBUG] Shutting down input reader")
			os.Exit(0)

		case enum.Connect:
			if state.IsConnected() {
				fmt.Println("[ERROR] Leave the current room before connection to a new one")
				break
			}

			state.ClientName = command.ClientName
			state.RoomName = command.RoomName
			state.View = enum.Lobby

		default:
			state.Broadcast <- &command
		}
	}
}

func CommandDispatcher() {
	fmt.Println("[DEBUG] Starting command dispatcher")
	for {
		command := <-state.Broadcast
		err := state.Conn.WriteJSON(command)
		fmt.Println("[DEBUG] Sending to server")
		if err != nil {
			fmt.Println(err.Error())
			break
		}
	}

	fmt.Println("[DEBUG] Shutting down command dispatcher")
}
