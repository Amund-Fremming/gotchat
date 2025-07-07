package ws

import (
	"fmt"
	"log/slog"
	"server/model"
	"strconv"
	"strings"

	"github.com/amund-fremming/common/enum"
	common "github.com/amund-fremming/common/model"
	"github.com/gorilla/websocket"
)

func sendServerError(content string, conn *websocket.Conn) {
	slog.Error(strings.ToLower(content))
	serverError := common.ServerError{View: enum.Lobby, Content: content}
	envelope := common.NewEnvelope(enum.ServerError, &serverError)

	err := conn.WriteJSON(envelope)
	if err != nil {
		slog.Error("Failed to write json to client", "address", conn.LocalAddr())
		conn.Close()
	}
}

func handleConnect(wrapper *model.ConnectionWrapper) {
	cmd, conn := wrapper.UnWrap()
	room, roomExists := state.TryGetRoom(cmd.RoomName)

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
}

func handleCreate(wrapper *model.ConnectionWrapper) {
	cmd, conn := wrapper.UnWrap()
	_, roomExists := state.TryGetRoom(cmd.RoomName)

	if roomExists {
		sendServerError("Room name is already in use", conn)
		return
	}

	newRoom := model.NewRoom(cmd.RoomName, conn)
	state.AddRoom(cmd.RoomName, &newRoom)

	go newRoom.Run()

	client := model.Client{Name: cmd.ClientName, Conn: conn}
	newRoom.Connect <- &client

	slog.Info(fmt.Sprintf("Client %s created room %s", cmd.ClientName, cmd.RoomName))
}

func handleSend(wrapper *model.ConnectionWrapper) {
	cmd, conn := wrapper.UnWrap()
	room, roomExists := state.TryGetRoom(cmd.RoomName)

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
}

func handleLeave(wrapper *model.ConnectionWrapper) {
	cmd, conn := wrapper.UnWrap()
	room, roomExists := state.TryGetRoom(cmd.RoomName)

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
	err := conn.WriteJSON(envelope)
	if err != nil {
		slog.Error("Failed write to socket", "client", cmd.ClientName)
		slog.Info("Closing failed connection", "client", cmd.ClientName)
		conn.Close()
		return
	}

	chatMessage := common.ChatMessage{Sender: "SERVER", Content: cmd.ClientName + " left the room..."}
	room.Chat <- &chatMessage

	fmt.Println("[ROOM] Client disconnected")
	slog.Info("Client left the room", "client", cmd.ClientName, "room", cmd.RoomName)
}

func handleExit(wrapper *model.ConnectionWrapper) {
	cmd, conn := wrapper.UnWrap()
	room, exists := state.TryGetRoom(cmd.RoomName)
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
