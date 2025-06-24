package main

import (
	"fmt"
	"net/http"
	"server/api"
	"server/ws"
)

func main() {
	http.HandleFunc("/health", api.Health)
	http.HandleFunc("/chat", ws.ClientHandler)

	fmt.Println("[SERVER] awaken")
	go ws.CommandRouter()
	http.ListenAndServe(":8080", nil)
}
