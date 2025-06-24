package model

import "github.com/amund-fremming/common/enum"

type Command struct {
	Action enum.Action `json:"action"`
	Arg    string      `json:"arg"`
}

func NewCommand(action enum.Action) Command {
	return Command{
		Action: action,
		Arg:    "",
	}
}
