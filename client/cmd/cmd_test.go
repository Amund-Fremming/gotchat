package cmd

import (
	"testing"

	"github.com/amund-fremming/common/enum"
)

func TestGetCommand_Successful(t *testing.T) {
	inputs := map[string]enum.Action{
		"this is a chat message": enum.Send,
		"/help":                  enum.Help,
		"/rooms":                 enum.Rooms,
		"/exit":                  enum.Exit,
		"/leave":                 enum.Leave,
		"/connect john room":     enum.Connect,
		"/create john room":      enum.Create,
	}

	for key, value := range inputs {
		cmd, err := GetCommand(key, "", "")
		if err != nil {
			t.Errorf("GetCommand failed with argument %s \n\t %s", value, err.Error())
			return
		}

		if cmd.Action != value {
			t.Errorf("expected=/help, recieved=%s", cmd.Action)
			return
		}
	}

	input := "/help"
	cmd, err := GetCommand(input, "", "")
	if err != nil {
		t.Error("GetCommand failed")
		return
	}

	if cmd.Action != enum.Help {
		t.Errorf("expected=/help, recieved=%s", cmd.Action)
		return
	}
}

func TestGetCommand_ShouldFail(t *testing.T) {
	inputs := map[string]enum.Action{
		"/connect": enum.Connect,
		"/create":  enum.Create,
	}

	for k := range inputs {
		_, err := GetCommand(k, "", "")
		if err == nil {
			t.Errorf("These commands should fail due to missing arguments")
		}
	}
}
