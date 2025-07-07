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
	output := fmt.Sprintf("\r[%s] %s\n", strings.ToLower(msg.Sender), msg.Content)
	rl.Write([]byte(output))
	rl.Refresh()
}

func DisplayError(content string) {
	output := fmt.Sprintf("\r[SERVER] %s\n", content)
	rl.Write([]byte(output))
	rl.Refresh()
}

func DisplayServerMessage(content string) {
	output := fmt.Sprintf("\r[SERVER] %s\n", content)
	rl.Write([]byte(output))
	rl.Refresh()
}

func DisplayErrorMessage(content string) {
	output := fmt.Sprintf("\r[ERROR] %s\n", content)
	rl.Write([]byte(output))
	rl.Refresh()
}
