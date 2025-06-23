package main

import (
	"client/handler"
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

	log.Println("[Server] Server is healthy")
	log.Println("[Client] Creating conneciton to the server")

	err = handler.ConnectWebSocket()
	if err != nil {
		log.Println(err.Error())
		return
	}
}
