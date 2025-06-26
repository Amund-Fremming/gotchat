package models

import (
	"fmt"

	"github.com/amund-fremming/common/model"
	"github.com/gorilla/websocket"
)

type Room struct {
	Clients   map[string]*websocket.Conn
	Connect   chan *Client
	Leave     chan string
	Broadcast chan model.Message
}

func NewRoom(n string, c *websocket.Conn) Room {
	return Room{
		Clients:   map[string]*websocket.Conn{n: c},
		Connect:   make(chan *Client, 10),
		Leave:     make(chan string, 10),
		Broadcast: make(chan model.Message, 10),
	}
}

func (r *Room) Run() {
	fmt.Println("[Room] Created and started")
	for {
		select {

		case client := <-r.Connect:
			r.Clients[client.Name] = client.Conn
			fmt.Println("[Room] client connected. Clients connected: ", len(r.Clients))

			r.Broadcast <- model.Message{
				Name: "SERVER",
				Body: client.Name + " joined the room...",
			}

		case name := <-r.Leave:
			delete(r.Clients, name)

		case message := <-r.Broadcast:
			for name, conn := range r.Clients {
				err := conn.WriteJSON(message)
				if err != nil {
					fmt.Println("[ERROR] failed to broadcast message")
					conn.Close()
					delete(r.Clients, name)
				}
			}
		}
	}
}
