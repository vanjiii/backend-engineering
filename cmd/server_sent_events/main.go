package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "welcome\n")
	})

	r.Get("/stream", stream)

	slog.Error(http.ListenAndServe(":5432", r).Error())
}

func stream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// You may need this locally for CORS requests
	w.Header().Set("Access-Control-Allow-Origin", "*")

	clientGone := r.Context().Done()

	rc := http.NewResponseController(w)
	t := time.NewTicker(time.Second)
	defer t.Stop()

	i := 0
	for {
		i++

		select {
		case <-clientGone:
			fmt.Fprintf(w, "client disconnected\n")
			return
		case <-t.C:
			slog.Info("iter", "i", i)

			// NOTE: the data should be in this format:
			// `data: <msg>\n\n`
			fmt.Fprintf(w, "data: iter: %v\n\n", i)
			rc.Flush()
		}
	}
}
