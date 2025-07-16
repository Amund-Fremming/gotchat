package main

import (
	"log/slog"
	"net/http"
	"server/api"
	"server/ws"

	"github.com/amund-fremming/gotchat-common/config"
)

func main() {
	config, err := config.Load()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.SetLogLoggerLevel(config.LogLevel)
	slog.Debug("Environemnt loaded", "env", config.Env)

	http.HandleFunc("/health", api.Health)
	http.HandleFunc("/chat", ws.ClientDispatcher)

	slog.Info("Server started")
	go ws.CommandDispatcher()

	http.ListenAndServe(":"+config.Port, nil)
}
