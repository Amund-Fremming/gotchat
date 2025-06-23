package models

import "sync"

type AppState struct {
	Mu       sync.RWMutex
	Channels map[string]chan string
}

func NewAppState() AppState {
	return AppState{
		Mu:       sync.RWMutex{},
		Channels: make(map[string]chan string),
	}
}
