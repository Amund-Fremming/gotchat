package models

import (
	"github.com/gorilla/websocket"
)

type Room struct {
	Join      chan []byte
	Leave     chan []byte
	Broadcast chan []byte
}

type RoomInfo struct {
	Name string
	Conn *websocket.Conn
}
