package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"backend-engineering/pkg/httpx/middleware"
)

func main() {
	http.HandleFunc("/get-updates", middleware.Request(getUpdates))

	slog.Info("server started", "port", "5432")

	http.ListenAndServe(":5432", nil)
}

func getUpdates(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"message": `{"name": "John", "age": 30}`,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
