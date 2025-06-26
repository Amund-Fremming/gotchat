package enum

type Type int

const (
	ClientState Type = iota
	ServerError
	ChatMessage
	Command
)

var typeName = map[Type]string{
	ClientState: "state",
	ServerError: "servererror",
	ChatMessage: "chatmessage",
	Command:     "command",
}

func (t Type) String() string {
	return typeName[t]
}
