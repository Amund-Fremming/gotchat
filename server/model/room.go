package model

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/amund-fremming/common/enum"
	"github.com/amund-fremming/common/model"
	"github.com/gorilla/websocket"
)

type Room struct {
	Mu      sync.RWMutex
	Name    string
	clients map[string]*websocket.Conn
	Connect chan *Client
	Leave   chan *Client
	Chat    chan *model.ChatMessage
}

func NewRoom(name string, c *websocket.Conn) Room {
	return Room{
		Mu:      sync.RWMutex{},
		Name:    name,
		clients: make(map[string]*websocket.Conn),
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
	slog.Info(fmt.Sprintf("ROOM [%s]: Created", r.Name))

ROOM:
	for {
		select {
		case client := <-r.Connect:
			r.SetClient(client.Name, client.Conn)
			slog.Info(fmt.Sprintf("ROOM [%s]: Client %s connected", r.Name, client.Name))
			slog.Info(fmt.Sprintf("ROOM [%s]: Connected clients: %d", r.Name, len(r.clients)))

			r.Chat <- &model.ChatMessage{Sender: "SERVER", Content: client.Name + " connected to the room..."}

		case client := <-r.Leave:
			client.Conn.Close()
			r.RemoveClient(client.Name)
			slog.Info(fmt.Sprintf("ROOM [%s]: Client %s left the room", r.Name, client.Name))

		case message := <-r.Chat:
			slog.Debug(fmt.Sprintf("ROOM [%s]: Client %s sendt a message", r.Name, message.Sender))
			envelope := model.NewEnvelope(enum.ChatMessage, message)
			r.Mu.RLock()
			for name, conn := range r.clients {
				err := conn.WriteJSON(envelope)
				if err != nil {
					slog.Error(fmt.Sprintf("ROOM [%s]: Client %s failed to write too room.", r.Name, message.Sender))
					slog.Error("Closing connection to client sender", "sender", message.Sender)
					r.RemoveClient(name)
					conn.Close()
				}
			}
			r.Mu.RUnlock()
		}

		if r.Empty() {
			slog.Info(fmt.Sprintf("ROOM [%s]: Closing room, no more clients connected", r.Name))
			break ROOM
		}
	}
}
