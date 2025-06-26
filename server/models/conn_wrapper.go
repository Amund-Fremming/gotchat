package models

import (
	"github.com/amund-fremming/common/model"
	"github.com/gorilla/websocket"
)

type Generic interface {
	model.Command | string
}

type ConnectionWrapper[T Generic] struct {
	Item T
	Conn *websocket.Conn
}
