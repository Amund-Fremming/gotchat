package models

import (
	"github.com/amund-fremming/common/model"
	"github.com/gorilla/websocket"
)

/*type Generic interface {
	model.Command | string
}

type ConnectionWrapper[T Generic] struct {
	item T
	conn *websocket.Conn
}

func (w *ConnectionWrapper[T]) UnWrap() (T, *websocket.Conn) {
	return w.item, w.conn
}*/

type ConnectionWrapper struct {
	Item model.Command
	Conn *websocket.Conn
}

func (w *ConnectionWrapper) UnWrap() (model.Command, *websocket.Conn) {
	return w.Item, w.Conn
}
