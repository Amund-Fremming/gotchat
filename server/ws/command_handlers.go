package ws

import (
	"fmt"
	"server/models"

	"github.com/amund-fremming/common/model"
)

func handleConnect(wrapper *models.ConnectionWrapper[model.Command]) {
	conn := wrapper.Conn
	cmd := wrapper.Item
	room, roomExists := state.GetRoom(cmd.RoomName)

	if !roomExists {
		fmt.Println("[CLIENT] Tried to send a message to non existing room")
		conn.WriteJSON(model.Message{
			Sender: "SERVER",
			Body:   "Room does not exist",
		})
		return
	}

	_, ok := room.Clients[cmd.ClientName]
	if ok {
		fmt.Println("[CLIENT] Tried to join a room with an existing name")
		conn.WriteJSON(model.Message{
			Sender: "SERVER",
			Body:   "Name is already in use",
		})
		return
	}

	client := models.Client{
		Name: cmd.ClientName,
		Conn: wrapper.Conn,
	}

	room.Connect <- &client
}

func handleCreate(wrapper *models.ConnectionWrapper[model.Command]) {
	conn := wrapper.Conn
	cmd := wrapper.Item
	room, roomExists := state.GetRoom(cmd.RoomName)

	if roomExists {
		fmt.Println("[CLIENT] Tried to create a room with an existing name")
		conn.WriteJSON(model.Message{
			Sender: "SERVER",
			Body:   "Name is already in use",
		})
		return
	}

	newRoom := models.NewRoom(cmd.ClientName, wrapper.Conn)
	state.AddRoom(cmd.RoomName, &newRoom)

	go room.Run()

	client := models.Client{
		Name: cmd.ClientName,
		Conn: wrapper.Conn,
	}

	room.Connect <- &client
	fmt.Println("[ROOM] Created")
}

func handleSend(wrapper *models.ConnectionWrapper[model.Command]) {
	conn := wrapper.Conn
	cmd := wrapper.Item
	room, roomExists := state.GetRoom(cmd.RoomName)

	if !roomExists {
		fmt.Println("[CLIENT] Tried to send a message to non existing room")
		conn.WriteJSON(model.Message{
			Sender: "SERVER",
			Body:   "Room does not exist",
		})
		return
	}

	_, ok := room.Clients[cmd.ClientName]
	if !ok {
		err := wrapper.Conn.WriteJSON(model.Message{
			Sender: "SERVER",
			Body:   "You are not connected to room `" + cmd.RoomName + "`.",
		})

		if err != nil {
			fmt.Println("[ERROR] Failed to write json to client")
			wrapper.Conn.Close()
		}
	}

	message := model.Message{
		Sender: cmd.ClientName,
		Body:   cmd.Message,
	}

	room.Broadcast <- message
	fmt.Println("[ROOM] Client sendt a message")
}

func handleLeave(wrapper *models.ConnectionWrapper[model.Command]) {
	conn := wrapper.Conn
	cmd := wrapper.Item
	room, roomExists := state.GetRoom(cmd.RoomName)

	if !roomExists {
		fmt.Println("[CLIENT] Tried to leave non existing room")
		conn.WriteJSON(model.Message{
			Sender: "SERVER",
			Body:   "Cannot leave non existing room",
		})
		return
	}

	conn, ok := room.Clients[cmd.ClientName]
	if ok {
		delete(room.Clients, cmd.ClientName)
		conn.WriteJSON(model.Message{
			Sender: "SERVER",
			Body:   "You left the room",
		})
	}

	message := model.Message{
		Sender: "SERVER",
		Body:   cmd.ClientName + " left the room...",
	}

	room.Broadcast <- message
	fmt.Println("[ROOM] Client disconnected")
}
