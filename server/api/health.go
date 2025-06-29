package api

import (
	"log/slog"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	slog.Info("Client requested health-check", "health", "OK")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{ "status": "OK"}`))
}
