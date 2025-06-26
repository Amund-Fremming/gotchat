package model

import (
	"encoding/json"

	"github.com/amund-fremming/common/enum"
)

type Envelope struct {
	Type    enum.Type       `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type ClientState struct {
	View       enum.View `json:"view"`
	RoomName   string    `json:"roomname"`
	ClientName string    `json:"clientname"`
}

type ChatMessage struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

type ServerError struct {
	View    enum.View `json:"view"`
	Content string    `json:"message"`
}
