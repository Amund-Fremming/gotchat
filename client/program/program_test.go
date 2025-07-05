package program

import (
	"testing"

	"github.com/amund-fremming/common/enum"
	"github.com/amund-fremming/common/model"
)

func TestMerge(t *testing.T) {
	state := NewAppState()
	clientState := model.ClientState{
		View:       enum.Room,
		RoomName:   "Room one",
		ClientName: "Client one",
	}

	state.Merge(&clientState)

	if state.View != clientState.View || state.RoomName != clientState.RoomName || state.ClientName != clientState.ClientName {
		t.Errorf("Failed to merge states, values differ")
	}
}

func TestCanExecuteLobbyCommand(t *testing.T) {
	state := AppState{View: enum.Lobby}

	lobbyCommands := [5]enum.Action{
		enum.Connect,
		enum.Exit,
		enum.Help,
		enum.Rooms,
		enum.Help,
	}

	for _, action := range lobbyCommands {
		cmd := model.NewCommand(action)
		if !state.CanExecuteCommand(&cmd) {
			t.Errorf("Should be able to execute command with action: %s, when in Lobby", action)
		}
	}

}

func TestCanExecuteRoomCommand(t *testing.T) {
	state := AppState{View: enum.Room}

	roomCommands := [3]enum.Action{
		enum.Send,
		enum.Leave,
		enum.Help,
	}

	for _, action := range roomCommands {
		cmd := model.NewCommand(action)
		if !state.CanExecuteCommand(&cmd) {
			t.Errorf("Should be able to execute command with action: %s, when in Room", action)
		}
	}
}
