package model

import (
	"encoding/json"
	"fmt"

	"github.com/amund-fremming/common/enum"
)

type Envelope struct {
	Type    enum.Type       `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type PayloadStruct interface {
	ClientState | ChatMessage | ServerError | Command
}

func NewEnvelope[T PayloadStruct](envelopeType enum.Type, payloadStruct T) Envelope {
	rawPayload, err := json.Marshal(payloadStruct)
	if err != nil {
		fmt.Println("[ERROR] Failed to marshal payload in envelope")
		rawPayload = []byte{}
	}

	envelope := Envelope{Type: envelopeType, Payload: rawPayload}
	return envelope
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
