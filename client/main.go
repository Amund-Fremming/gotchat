package main

import (
	"client/program"
	"log"
	"net/http"
)

const ServerUrlBase = "http://localhost:8080"

func main() {
	log.SetFlags(0)

	response, err := http.Get(ServerUrlBase + "/health")
	if err != nil || response.Status != "200 OK" {
		log.Println("[Error] The sever is currently unavailable")
		return
	}
	log.Println("[Server] is healthy")

	go program.ConnectToServer()
	go program.InputReader()
	go program.CommandDispatcher()

	select {}
}
