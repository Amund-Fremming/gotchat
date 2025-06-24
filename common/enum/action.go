package enum

type Action int

const (
	Help Action = iota
	Connect
	Create
	Status
	Exit
	Send
)

var actionNames = map[Action]string{
	Help:    "help",
	Connect: "connect",
	Create:  "create",
	Status:  "status",
	Exit:    "exit",
	Send:    "send",
}

func (a Action) String() string {
	return actionNames[a]
}
