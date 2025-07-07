package cmd

import (
	"errors"
	"strings"

	"github.com/amund-fremming/common/enum"
	"github.com/amund-fremming/common/model"
	"github.com/chzyer/readline"
)

type Command = model.Command

var rl *readline.Instance

func InitReadline() error {
	var err error
	rl, err = readline.New("> ")
	return err
}

func SetPrompt(prompt string) {
	rl.SetPrompt(prompt)
	rl.Refresh()
}

func ReadInput() string {
	line, err := rl.Readline()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(line)
}

func GetCommand(input string, clientName string, roomName string) (Command, error) {
	isChatMessage := !strings.HasPrefix(input, "/")

	if isChatMessage {
		return model.Command{
			Action:     enum.Send,
			ClientName: clientName,
			RoomName:   roomName,
			Message:    input,
		}, nil
	}

	verbs := strings.Split(input, " ")

	switch verbs[0] {
	case "/help":
		return model.NewCommand(enum.Help), nil

	case "/rooms":
		return model.NewCommand(enum.Rooms), nil

	case "/exit":
		return model.NewCommand(enum.Exit), nil

	case "/leave":
		return model.Command{
			Action:     enum.Leave,
			ClientName: clientName,
			RoomName:   roomName,
		}, nil

	case "/connect", "/create":
		if len(verbs) < 3 {
			return Command{}, errors.New("[ERROR] This command required two arguments")
		}

		action := enum.Connect
		if verbs[0] == "/create" {
			action = enum.Create
		}

		return model.Command{
			Action:     action,
			ClientName: verbs[1],
			RoomName:   verbs[2],
		}, nil
	}

	return Command{}, errors.New("[ERROR] Invalid command")
}
