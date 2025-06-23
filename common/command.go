package common

type Action int

const (
	Help Action = iota
	Connect
	Create
	Status
	Exit
	Send
)

type Command struct {
	Action Action `json:"action"`
	Arg    string `json:"arg"`
}

func NewCommand(action Action) Command {
	return Command{
		Action: action,
		Arg:    "",
	}
}
