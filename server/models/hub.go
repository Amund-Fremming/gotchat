package models

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Name       string
	Connection websocket.Conn
}

type Server struct {
	Mu        sync.RWMutex
	Name      string
	Clients   []Client
	Broadcast chan string
}

type Room struct {
	Join      chan []byte
	Leave     chan []byte
	Broadcast chan []byte
}

type RoomInfo struct {
	Name string
	Conn *websocket.Conn
}

func NewServer(name string) Server {
	return Server{
		Mu:        sync.RWMutex{},
		Name:      name,
		Clients:   make([]Client, 20),
		Broadcast: make(chan string),
	}
}
