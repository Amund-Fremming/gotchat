package program

import (
	"github.com/amund-fremming/common/enum"
	"github.com/amund-fremming/common/model"
	"github.com/gorilla/websocket"
)

type AppState struct {
	Broadcast chan *model.Command
	Conn      *websocket.Conn
	AppView   enum.View
}

func NewAppState() AppState {
	return AppState{
		Broadcast: make(chan *model.Command),
		Conn:      &websocket.Conn{},
		AppView:   enum.Lobby,
	}
}
