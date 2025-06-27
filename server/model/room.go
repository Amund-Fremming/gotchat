package model

import (
	"fmt"
	"sync"

	"github.com/amund-fremming/common/enum"
	"github.com/amund-fremming/common/model"
	"github.com/gorilla/websocket"
)

type Room struct {
	Mu      sync.RWMutex
	clients map[string]*websocket.Conn
	Connect chan *Client
	Leave   chan *Client
	Chat    chan *model.ChatMessage
}

func NewRoom(n string, c *websocket.Conn) Room {
	return Room{
		Mu:      sync.RWMutex{},
		clients: map[string]*websocket.Conn{n: c},
		Connect: make(chan *Client, 100),
		Leave:   make(chan *Client, 100),
		Chat:    make(chan *model.ChatMessage, 100),
	}
}

func (r *Room) SetClient(name string, conn *websocket.Conn) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.clients[name] = conn
}

func (r *Room) TryGetClient(name string) (*websocket.Conn, bool) {
	r.Mu.RLock()
	defer r.Mu.RUnlock()
	value, ok := r.clients[name]
	return value, ok
}

func (r *Room) RemoveClient(name string) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	delete(r.clients, name)
}

func (r *Room) Empty() bool {
	r.Mu.RLock()
	defer r.Mu.RUnlock()
	return len(r.clients) == 0
}

func (r *Room) Count() int {
	r.Mu.RLock()
	defer r.Mu.RUnlock()
	return len(r.clients)
}

func (r *Room) Run() {
	fmt.Println("[ROOM] Created")

ROOM:
	for {
		select {
		case client := <-r.Connect:
			r.SetClient(client.Name, client.Conn)
			fmt.Println("[ROOM] Client connected. Size:", len(r.clients))

			r.Chat <- &model.ChatMessage{Sender: "SERVER", Content: client.Name + " connected to the room..."}

		case client := <-r.Leave:
			client.Conn.Close()
			r.RemoveClient(client.Name)
			fmt.Println("[SERVER] " + client.Name + " left the room..")

		case message := <-r.Chat:
			envelope := model.NewEnvelope(enum.ChatMessage, message)
			r.Mu.RLock()
			for name, conn := range r.clients {
				err := conn.WriteJSON(envelope)
				if err != nil {
					fmt.Println("[ERROR] Failed to broadcast message")
					r.RemoveClient(name)
					conn.Close()
				}
			}
			r.Mu.RUnlock()
		}

		if r.Empty() {
			fmt.Println("[ROOM] Room is closing itself")
			break ROOM
		}
	}
}
