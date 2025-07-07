package model

import (
	"github.com/amund-fremming/common/enum"
)

type ClientState struct {
	View       enum.View `json:"view"`
	RoomName   string    `json:"roomname"`
	ClientName string    `json:"clientname"`
	Prompt     string    `json:"prompt"`
}

type RoomData struct {
	Content string `json:"content"`
}

type ChatMessage struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

type ServerError struct {
	View    enum.View `json:"view"`
	Content string    `json:"message"`
}

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
