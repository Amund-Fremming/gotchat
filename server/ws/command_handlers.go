package ws

import (
	"fmt"
	"server/model"
	"strconv"
	"strings"

	"github.com/amund-fremming/common/enum"
	common "github.com/amund-fremming/common/model"
	"github.com/gorilla/websocket"
)

func sendServerError(content string, conn *websocket.Conn) {
	fmt.Println("[SERVER] " + content)

	serverError := common.ServerError{View: enum.Lobby, Content: content}
	envelope := common.NewEnvelope(enum.ServerError, &serverError)

	err := conn.WriteJSON(envelope)
	if err != nil {
		fmt.Println("[ERROR] Failed to write json to client")
		conn.Close()
	}
}

func handleConnect(wrapper *model.ConnectionWrapper) {
	cmd, conn := wrapper.UnWrap()
	room, roomExists := state.GetRoom(cmd.RoomName)

	if !roomExists {
		sendServerError("Room does not exist", conn)
		return
	}

	_, ok := room.TryGetClient(cmd.ClientName)
	if ok {
		sendServerError("Username is already in use", conn)
		return
	}

	client := model.Client{Name: cmd.ClientName, Conn: conn}
	room.Connect <- &client

	fmt.Println("[CLIENT] Connected")
}

func handleCreate(wrapper *model.ConnectionWrapper) {
	cmd, conn := wrapper.UnWrap()
	_, roomExists := state.GetRoom(cmd.RoomName)

	if roomExists {
		sendServerError("Room name is already in use", conn)
		return
	}

	newRoom := model.NewRoom(cmd.ClientName, conn)
	state.AddRoom(cmd.RoomName, &newRoom)

	go newRoom.Run()

	client := model.Client{Name: cmd.ClientName, Conn: conn}
	newRoom.Connect <- &client

	fmt.Println("[CLIENT] Created room")
}

func handleSend(wrapper *model.ConnectionWrapper) {
	cmd, conn := wrapper.UnWrap()
	room, roomExists := state.GetRoom(cmd.RoomName)

	if !roomExists {
		sendServerError("Room does not exist", conn)
		return
	}

	_, ok := room.TryGetClient(cmd.ClientName)
	if !ok {
		sendServerError("You are not connected to this room.", conn)
		return
	}

	message := common.ChatMessage{Sender: cmd.ClientName, Content: cmd.Message}
	room.Chat <- &message

	fmt.Println("[ROOM] Client sendt a message")
}

func handleLeave(wrapper *model.ConnectionWrapper) {
	cmd, conn := wrapper.UnWrap()
	room, roomExists := state.GetRoom(cmd.RoomName)

	if !roomExists {
		sendServerError("Cannot leave non-existing room", conn)
		return
	}

	_, ok := room.TryGetClient(cmd.ClientName)
	if !ok {
		sendServerError("Cannot leave non entered room", conn)
		return
	}

	room.RemoveClient(cmd.ClientName)
	if room.Empty() {
		state.RemoveRoom(cmd.RoomName)
	}

	message := common.ChatMessage{Sender: "SERVER", Content: "You left the room"}
	envelope := common.NewEnvelope(enum.ChatMessage, &message)
	conn.WriteJSON(envelope)

	chatMessage := common.ChatMessage{Sender: "SERVER", Content: cmd.ClientName + " left the room..."}
	room.Chat <- &chatMessage

	fmt.Println("[ROOM] Client disconnected")
}

func handleExit(wrapper *model.ConnectionWrapper) {
	cmd, conn := wrapper.UnWrap()
	room, exists := state.GetRoom(cmd.RoomName)
	if !exists {
		return
	}

	client := model.Client{Name: cmd.ClientName, Conn: conn}
	room.Leave <- &client
}

func handleRooms(wrapper *model.ConnectionWrapper) {
	_, conn := wrapper.UnWrap()
	var sb strings.Builder
	sb.WriteString("\nROOMS\n")

	i := 1
	for k, v := range state.GetRoomsUnsafe() {
		count := strconv.Itoa(v.Count())
		var nl string = ""
		if i%3 == 0 {
			nl = "\n"
		}
		i++

		v.Mu.RLock()
		sb.WriteString(fmt.Sprintf("[%-10s -> %-2s]  %s", k, count, nl))
		v.Mu.RUnlock()
	}

	roomsData := common.RoomData{Content: sb.String()}
	envelope := common.NewEnvelope(enum.RoomsData, &roomsData)

	conn.WriteJSON(envelope)
}
