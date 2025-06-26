package main

import (
	"client/program"
	"fmt"
	"net/http"
)

const ServerUrlBase = "http://localhost:8080"

func main() {
	response, err := http.Get(ServerUrlBase + "/health")
	if err != nil || response.Status != "200 OK" {
		fmt.Println("[ERROR] The sever is currently unavailable")
		fmt.Println("[CLIENT] Shutting down..")
		return
	}
	fmt.Println("[SERVER] Healthy")

	program.ConnectToServer()
	go program.CommandReader()
	go program.ServerReader()
	go program.CommandDispatcher()

	select {}
}
