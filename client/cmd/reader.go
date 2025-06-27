package cmd

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/amund-fremming/common/enum"
	"github.com/amund-fremming/common/model"
)

type Command = model.Command

func ReadInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	return input
}

// TODO: This is straight ugly, fix it
func GetCommand(clientName string, roomName string) (Command, error) {
	var input string = ReadInput()
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
	case "/connect":
		if len(verbs) < 3 {
			return Command{}, errors.New("[ERROR] This command required two arguments")
		}
		return model.Command{
			Action:     enum.Connect,
			ClientName: verbs[1],
			RoomName:   verbs[2],
		}, nil
	case "/create":
		if len(verbs) < 3 {
			return Command{}, errors.New("[ERROR] This command required two arguments")
		}
		return model.Command{
			Action:     enum.Create,
			ClientName: verbs[1],
			RoomName:   verbs[2],
		}, nil

	}

	return Command{}, errors.New("[ERROR] Invalid command")
}
