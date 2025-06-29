package main

import (
	"log/slog"
	"net/http"
	"server/api"
	"server/config"
	"server/ws"
)

func main() {
	config.LoadEnv()
	env := config.GetEnv()

	switch env {
	case config.Production:
		slog.SetLogLoggerLevel(slog.LevelError)
	case config.Development:
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	slog.Debug("Environemnt loaded", "env", env)

	http.HandleFunc("/health", api.Health)
	http.HandleFunc("/chat", ws.ClientDispatcher)

	slog.Info("Server started")
	go ws.CommandDispatcher()
	http.ListenAndServe(":8080", nil)
}
