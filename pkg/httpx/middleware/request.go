package middleware

import (
	"log/slog"
	"net/http"
)

func Request(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("request made",
			"method", r.Method,
			"url", r.URL.Path,
		)

		fn(w, r)
	}
}
