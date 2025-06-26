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
var rooms = make(map[string]*models.Room)

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

	fmt.Println("[Connected] Client connected")
	go commandReader(conn)
}

func commandReader(conn *websocket.Conn) {
	defer func() {
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("[ERROR] Failed to read message from client, closing connection")
			break
		}

		var cmd model.Command
		err = json.Unmarshal(msg, &cmd)
		if err != nil {
			slog.Error("Failed to parse command")
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
		room, roomExists := rooms[cmd.RoomName]

		fmt.Println("[Command Dispatcher] dispatching command:", cmd.Action)

		switch cmd.Action {
		case enum.Connect:
			if !roomExists {
				fmt.Println("[DEBUG] tried to send a message to non existing room")
				wrapper.Conn.WriteJSON(model.Message{
					Name: "SERVER",
					Body: "Room does not exist",
				})
				break
			}

			_, ok := room.Clients[cmd.ClientName]
			if ok {
				fmt.Println("[DEBUG] tried to join a room with an existing name")
				wrapper.Conn.WriteJSON(model.Message{
					Name: "SERVER",
					Body: "Name is already in use",
				})
				break
			}

			client := models.Client{
				Name: cmd.ClientName,
				Conn: wrapper.Conn,
			}

			room.Connect <- &client

		case enum.Create:
			if roomExists {
				fmt.Println("[DEBUG] tried to create a room with an existing name")
				wrapper.Conn.WriteJSON(model.Message{
					Name: "SERVER",
					Body: "Name is already in use",
				})
				break
			}

			room := models.NewRoom(cmd.ClientName, wrapper.Conn)
			rooms[cmd.RoomName] = &room

			go room.Run()

			client := models.Client{
				Name: cmd.ClientName,
				Conn: wrapper.Conn,
			}

			room.Connect <- &client

			fmt.Println("[Room] was created")

		case enum.Send:
			if !roomExists {
				fmt.Println("[DEBUG] tried to send a message to non existing room")
				wrapper.Conn.WriteJSON(model.Message{
					Name: "SERVER",
					Body: "Room does not exist",
				})
				break
			}

			_, ok := room.Clients[cmd.ClientName]
			if !ok {
				err := wrapper.Conn.WriteJSON(model.Message{
					Name: "SERVER",
					Body: "Unauthorized",
				})

				if err != nil {
					fmt.Println("[ERROR] failed to write json to client")
					wrapper.Conn.Close()
				}
			}

			message := model.Message{
				Name: cmd.ClientName,
				Body: cmd.Message,
			}

			room.Broadcast <- message

			fmt.Println("[Room] client sendt a message")

		case enum.Leave:
			_, ok := room.Clients[cmd.ClientName]
			if ok {
				delete(room.Clients, cmd.ClientName)
			}

			fmt.Println("[Room] client disconnected")
		}
	}
}
