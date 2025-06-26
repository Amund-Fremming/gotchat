package ws

import (
	"encoding/json"
	"fmt"
	"server/models"

	"github.com/amund-fremming/common/enum"
	"github.com/amund-fremming/common/model"
	"github.com/gorilla/websocket"
)

func handleError(content string, conn *websocket.Conn) {
	fmt.Println("[SERVER] " + content)

	serverError := model.ServerError{View: enum.Lobby, Content: content}
	rawServerError, _ := json.Marshal(serverError)
	envelope := model.Envelope{Type: enum.ServerError, Payload: rawServerError}

	err := conn.WriteJSON(envelope)
	if err != nil {
		fmt.Println("[MAJOR_ERROR] Failed to write json to client")
		conn.Close()
	}
}

func handleConnect(wrapper *models.ConnectionWrapper) {
	cmd, conn := wrapper.UnWrap()
	room, roomExists := state.GetRoom(cmd.RoomName)

	if !roomExists {
		handleError("Room does not exist", conn)
		return
	}

	_, ok := room.Clients[cmd.ClientName]
	if ok {
		handleError("Username is already in use", conn)
		return
	}

	client := models.Client{Name: cmd.ClientName, Conn: conn}
	room.Connect <- &client

	fmt.Println("[CLIENT] Connected")
}

func handleCreate(state *models.AppState, wrapper *models.ConnectionWrapper) {
	cmd, conn := wrapper.UnWrap()
	_, roomExists := state.GetRoom(cmd.RoomName)

	if roomExists {
		handleError("Room name is already in use", conn)
		return
	}

	newRoom := models.NewRoom(cmd.ClientName, conn)
	state.AddRoom(cmd.RoomName, &newRoom)

	go newRoom.Run()

	client := models.Client{Name: cmd.ClientName, Conn: conn}
	newRoom.Connect <- &client

	fmt.Println("[CLIENT] Created room")
}

func handleSend(wrapper *models.ConnectionWrapper) {
	cmd, conn := wrapper.UnWrap()
	room, roomExists := state.GetRoom(cmd.RoomName)

	if !roomExists {
		handleError("Room does not exist", conn)
		return
	}

	_, ok := room.Clients[cmd.ClientName]
	if !ok {
		handleError("You are not connected to this room.", conn)
		return
	}

	message := model.ChatMessage{Sender: cmd.ClientName, Content: cmd.Message}
	room.Broadcast <- message

	fmt.Println("[ROOM] Client sendt a message")
}

func handleLeave(wrapper *models.ConnectionWrapper) {
	cmd, conn := wrapper.UnWrap()
	room, roomExists := state.GetRoom(cmd.RoomName)

	if !roomExists {
		handleError("Cannot leave non-existing room", conn)
		return
	}

	conn, ok := room.Clients[cmd.ClientName]
	if !ok {
		handleError("Cannot leave non entered room", conn)
		return
	}

	delete(room.Clients, cmd.ClientName)
	message := model.ChatMessage{Sender: "SERVER", Content: "You left the room"}
	rawMessage, _ := json.Marshal(message)

	envelope := model.Envelope{Type: enum.ChatMessage, Payload: rawMessage}
	conn.WriteJSON(envelope)

	roomMessage := model.ChatMessage{Sender: "SERVER", Content: cmd.ClientName + " left the room..."}
	room.Broadcast <- roomMessage

	fmt.Println("[ROOM] Client disconnected")
}
