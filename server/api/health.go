package api

import (
	"fmt"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[CLIENT] Requested health-check")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{ "status": "OK"}`))
}
