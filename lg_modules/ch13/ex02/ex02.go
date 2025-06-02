package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func RequestLogger(l *slog.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			h.ServeHTTP(w, r)
			l.Info("GET Request",
				"time", start.Format(time.RFC3339),
				"remote_addr", r.RemoteAddr,
			)
		})
	}
}

func main() {
	options := &slog.HandlerOptions{Level: slog.LevelInfo}
	handler := slog.NewJSONHandler(os.Stderr, options)
	mySlog := slog.New(handler)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now().Format(time.RFC3339)
		w.WriteHeader(http.StatusOK)
		w.Write(fmt.Appendf(nil, "current time: %s\n", currentTime))
	})
	wrappedMux := RequestLogger(mySlog)(mux)
	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      wrappedMux,
	}

	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}
