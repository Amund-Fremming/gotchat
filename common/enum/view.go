package enum

type View int

const (
	Lobby View = iota
	Hub
)

var viewNames = map[View]string{
	Lobby: "lobby",
	Hub:   "hub",
}

func (v View) String() string {
	return viewNames[v]
}
