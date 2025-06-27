package enum

type Action int

const (
	Help Action = iota
	Connect
	Create
	Rooms
	Leave
	Exit
	Send
)

var actionNames = map[Action]string{
	Help:    "help",
	Connect: "connect",
	Create:  "create",
	Rooms:   "rooms",
	Leave:   "leave",
	Exit:    "exit",
	Send:    "send",
}

func (a Action) String() string {
	return actionNames[a]
}
