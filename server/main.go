package main

import (
	"log/slog"
	"net/http"
	"server/api"
	"server/ws"
)

func main() {
	http.HandleFunc("/health", api.Health)
	http.HandleFunc("/chat", ws.ClientHandler)

	slog.Info("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
