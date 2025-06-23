package main

import (
	"log/slog"
	"net/http"
	"server/api"
	"server/ws"
)

func main() {
	http.HandleFunc("/health", api.Health)
	http.HandleFunc("/chat", ws.HandleWebSocket)

	go ws.HandleMessages()

	slog.Info("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
