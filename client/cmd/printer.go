package cmd

import (
	"fmt"
	"strings"

	"github.com/amund-fremming/common/model"
)

func DisplayWelcomeMessage() {
	fmt.Println(`
Welcome to go tchat!
    - You are now in the lobby
    - Start by running the "/help" command
    `)
}

func DisplayCommands() {
	fmt.Println(`
Lobby commands:
    /help                            Displays available commands in you context
    /connect <username> <room_name>  Connects a user to a room
    /create  <name>                  Creates a room with name "<name>"
    /status                          Displays all available rooms with a counter
    /exit                            Disconnects the client and shuts down the app

Room commands:
    /help                            Displays available commands in you context
    /leave                           Exits the room back to the lobby
    <message>                        Send a message by typing a "<message>" and hit enter
	`)
}

func DisplayMessage(msg *model.ChatMessage) {
	fmt.Println(strings.ToLower("[" + msg.Sender + "] " + msg.Content))
}

func DisplayError(content string) {
	fmt.Println("[SERVER] " + content)
}

func DisplayServerMessage(content string) {
	fmt.Println("[SERVER]", content)
}

func DisplayErrorMessage(content string) {
	fmt.Println("[ERROR]", content)
}

type CursorAction int

const (
	MoveTwoUp CursorAction = iota
	MoveTwoDown
	MoveLinestart
	DeleteAllRight
	Newline
)

func handleCursor(a CursorAction) {
	switch a {
	case MoveTwoUp:
		fmt.Print("\033[A]")
		fmt.Print("\033[A]")
	case MoveTwoDown:
		fmt.Print("\033[B]")
		fmt.Print("\033[B]")
	case MoveLinestart:
		fmt.Print("\r")
	case DeleteAllRight:
		fmt.Print("\033[K]")
	case Newline:
		fmt.Println()
	}
}
