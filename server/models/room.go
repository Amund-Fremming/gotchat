package models

import (
	"encoding/json"
	"fmt"

	"github.com/amund-fremming/common/enum"
	"github.com/amund-fremming/common/model"
	"github.com/gorilla/websocket"
)

type Room struct {
	Clients   map[string]*websocket.Conn
	Connect   chan *Client
	Leave     chan string
	Broadcast chan model.ChatMessage
}

func NewRoom(n string, c *websocket.Conn) Room {
	return Room{
		Clients:   map[string]*websocket.Conn{n: c},
		Connect:   make(chan *Client, 10),
		Leave:     make(chan string, 10),
		Broadcast: make(chan model.ChatMessage, 10),
	}
}

func (r *Room) Run() {
	fmt.Println("[ROOM] Created")
	for {
		select {
		case client := <-r.Connect:
			r.Clients[client.Name] = client.Conn
			fmt.Println("[ROOM] Client connected. Size:", len(r.Clients))

			r.Broadcast <- model.ChatMessage{Sender: "SERVER", Content: client.Name + " connected to the room..."}

		case name := <-r.Leave:
			delete(r.Clients, name)

		case message := <-r.Broadcast:
			rawMessage, _ := json.Marshal(message)
			envelope := model.Envelope{Type: enum.ChatMessage, Payload: rawMessage}

			for name, conn := range r.Clients {
				err := conn.WriteJSON(envelope)
				if err != nil {
					fmt.Println("[ERROR] Failed to broadcast message")
					conn.Close()
					delete(r.Clients, name)
				}
			}
		}
	}
}
