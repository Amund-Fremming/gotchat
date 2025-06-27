package model

import "sync"

type AppState struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

func NewAppState() AppState {
	return AppState{
		rooms: make(map[string]*Room),
		mu:    sync.RWMutex{},
	}
}

func (state *AppState) AddRoom(name string, room *Room) {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.rooms[name] = room
}

func (state *AppState) RemoveRoom(name string) {
	state.mu.Lock()
	defer state.mu.Unlock()
	delete(state.rooms, name)
}

func (state *AppState) GetRoom(name string) (*Room, bool) {
	state.mu.RLock()
	defer state.mu.RUnlock()

	room, ok := state.rooms[name]
	return room, ok
}
