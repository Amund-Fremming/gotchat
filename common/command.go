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
	Action Action
	Arg    string
}

func NewCommand(action Action) Command {
	return Command{
		Action: action,
		Arg:    "",
	}
}
