package common

type View int

const (
	Lobby View = iota
	Hub
)

type Message struct {
	Name string
	Text string
}
