package model

import (
	"github.com/amund-fremming/common/enum"
)

type Command struct {
	Action     enum.Action `json:"action"`
	Message    string      `json:"message"`
	RoomName   string      `json:"roomname"`
	ClientName string      `json:"clientname"`
}

func NewCommand(action enum.Action) Command {
	return Command{
		Action:     action,
		Message:    "",
		RoomName:   "",
		ClientName: "",
	}
}
