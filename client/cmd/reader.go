package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/amund-fremming/common"
)

type Command = common.Command

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("> ")
	input, _ := reader.ReadString('\n')
	return input
}

func GetCommand() (Command, error) {
	var input string = readInput()
	verbs := strings.Split(input, " ")

	if len(verbs) > 2 {
		return Command{}, errors.New("commands can only have one action and one argument")
	}

	if len(verbs) == 1 {
		runes := []rune(verbs[0])
		if runes[0] == '/' {
			return Command{}, errors.New("command was invalid")
		}

		command := Command{
			Action: common.Send,
			Arg:    verbs[1],
		}

		return command, nil
	}

	switch verbs[0] {
	case "/help":
		return common.NewCommand(common.Help), nil
	case "/connect":
		return common.NewCommand(common.Connect), nil
	case "/create":
		return common.NewCommand(common.Create), nil
	case "/status":
		return common.NewCommand(common.Status), nil
	case "/exit":
		return common.NewCommand(common.Exit), nil
	}

	return Command{}, errors.New("Command was invalid")
}
