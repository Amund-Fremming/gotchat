package model

import (
	"testing"

	"github.com/gorilla/websocket"
)

func TestHandleRoom(t *testing.T) {
	roomName := "RoomOne"
	state := NewAppState()
	room := NewRoom(roomName, &websocket.Conn{})

	if _, exists := state.TryGetRoom(roomName); exists {
		t.Errorf("Room with name %s should not exist in the state", roomName)
	}

	state.AddRoom(roomName, &room)
	value, exists := state.TryGetRoom(roomName)
	if !exists {
		t.Errorf("Failed to get or add room with name %s", roomName)
	}

	if value.Name != roomName {
		t.Errorf("Failed to get the correct room. Recieved name: %s, expected name: %s", value.Name, roomName)
	}

	state.RemoveRoom(roomName)
	if _, exists := state.TryGetRoom(roomName); exists {
		t.Errorf("Failed to remove room with name %s", roomName)
	}
}
