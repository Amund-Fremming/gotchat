package main

import (
	"client/cmd"
	"client/program"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/amund-fremming/common/config"
)

const ServerUrlBase = "http://localhost:8080"

func main() {
	config, err := config.Load()
	if err != nil {
		fmt.Println("[CLIENT] Failed due to missing enviroment variables")
		return
	}

	const ServerUrlBase = ""
	slog.SetLogLoggerLevel(config.LogLevel)

	response, err := http.Get(ServerUrlBase + "/health")
	if err != nil || response.Status != "200 OK" {
		fmt.Println("[ERROR] The sever is currently unavailable")
		fmt.Println("[CLIENT] Shutting down..")
		return
	}
	fmt.Println("[SERVER] Healthy")

	program.ConnectToServer(&config)
	cmd.DisplayWelcomeMessage()

	go program.CommandReader()
	go program.ServerReader()
	go program.CommandDispatcher()

	select {}
}
