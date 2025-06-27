package main

import (
	"fmt"
	"net/http"
	"server/api"
	"server/ws"
)

func main() {
	http.HandleFunc("/health", api.Health)
	http.HandleFunc("/chat", ws.ClientDispatcher)

	fmt.Println("[SERVER] Started")

	go ws.CommandDispatcher()
	http.ListenAndServe(":8080", nil)
}
