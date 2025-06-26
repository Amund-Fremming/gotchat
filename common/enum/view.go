package enum

type View int

const (
	Lobby View = iota
	Room
)

var viewNames = map[View]string{
	Lobby: "lobby",
	Room:  "room",
}

func (v View) String() string {
	return viewNames[v]
}
