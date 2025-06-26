package program

import (
	"github.com/amund-fremming/common/enum"
	"github.com/amund-fremming/common/model"
	"github.com/gorilla/websocket"
)

type AppState struct {
	Broadcast  chan *model.Command
	Conn       *websocket.Conn
	View       enum.View
	ClientName string
	RoomName   string
}

func NewAppState() AppState {
	return AppState{
		Broadcast:  make(chan *model.Command),
		Conn:       &websocket.Conn{},
		View:       enum.Lobby,
		ClientName: "",
		RoomName:   "",
	}
}

func (s *AppState) Clear() {
	s.View = enum.Lobby
	s.ClientName = ""
	s.RoomName = ""
}

func (s *AppState) IsConnected() bool {
	return s.View == enum.Room
}

func (s *AppState) Merge(cs *model.ClientState) {
	s.View = cs.View
	s.ClientName = cs.ClientName
	s.RoomName = cs.RoomName
}
