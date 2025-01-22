package middleware

import (
	"log/slog"
	"net/http"
)

func Func(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("request made",
			"method", r.Method,
			"url", r.URL.Path,
		)

		fn(w, r)
	}
}

func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("request made",
			"method", r.Method,
			"url", r.URL.Path,
		)

		next.ServeHTTP(w, r)
	})
}
