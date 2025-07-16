package model

import (
	"github.com/amund-fremming/common/model"
	"github.com/gorilla/websocket"
)

type ConnectionWrapper struct {
	Item model.Command
	Conn *websocket.Conn
}

func (w *ConnectionWrapper) UnWrap() (model.Command, *websocket.Conn) {
	return w.Item, w.Conn
}
