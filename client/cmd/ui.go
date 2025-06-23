package cmd

import (
	"log"

	"github.com/amund-fremming/common"
)

func DisplayWelcomeMessage() {
	log.Println(`
Welcome to go tchat!
    - You are now in the lobby
    - Start by running the "/help" command`)
}

func DisplayCommands() {
	log.Println(`
Lobby commands:
    /help             Displays available commands in you context
    /connect <name>   Connect to a room with name "<name>"
    /create  <name>   Creates a room with name "<name>"
    /status           Displays all available rooms with a counter
    /exit             Disconnects the client and shuts down the app

Room commands:
    /help             Displays available commands in you context
    /exit             Exits the room back to the lobby
    <message>         Send a message by typing a "<message>" and hit enter
	`)
}

func DisplayMessage(msg *common.Message) {
	log.Println("[" + msg.Name + "]" + msg.Text)
}
