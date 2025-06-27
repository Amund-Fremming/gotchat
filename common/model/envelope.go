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
	ClientState | ChatMessage | ServerError | Command | RoomData
}

func NewEnvelope[T PayloadStruct](envelopeType enum.Type, payloadStruct *T) Envelope {
	rawPayload, err := json.Marshal(payloadStruct)
	if err != nil {
		fmt.Println("[ERROR] Failed to marshal payload in envelope")
		rawPayload = []byte{}
	}

	envelope := Envelope{Type: envelopeType, Payload: rawPayload}
	return envelope
}
