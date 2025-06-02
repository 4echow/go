package main

import (
	"encoding/json"
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

type TimeData struct {
	DayOfWeek  string `json:"day_of_week"`
	DayOfMonth int    `json:"day_of_month"`
	Month      string `json:"month"`
	Year       int    `json:"year"`
	Hour       int    `json:"hour"`
	Minute     int    `json:"minute"`
	Second     int    `json:"second"`
}

func toJsonTime(t time.Time) string {
	JsonStruct := TimeData{
		DayOfWeek:  t.Weekday().String(),
		DayOfMonth: t.Day(),
		Month:      t.Month().String(),
		Year:       t.Year(),
		Hour:       t.Hour(),
		Minute:     t.Minute(),
		Second:     t.Second(),
	}
	JsonOut, _ := json.Marshal(JsonStruct)
	return string(JsonOut)
}

func toTextTime(t time.Time) string {
	rfc3339time := t.Format(time.RFC3339)
	return rfc3339time
}

func main() {
	options := &slog.HandlerOptions{Level: slog.LevelInfo}
	handler := slog.NewJSONHandler(os.Stderr, options)
	mySlog := slog.New(handler)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now()
		accept := r.Header.Get("Accept")
		switch {
		case accept == "application/json":
			out := toJsonTime(currentTime)
			w.WriteHeader(http.StatusOK)
			w.Write(fmt.Appendf(nil, "current time JSON: %s\n", out))
		case accept == "text/html":
			out := toTextTime(currentTime)
			w.WriteHeader(http.StatusOK)
			w.Write(fmt.Appendf(nil, "current time: %s\n", out))
		default:
			out := toTextTime(currentTime)
			w.WriteHeader(http.StatusOK)
			w.Write(fmt.Appendf(nil, "current time: %s\n", out))
		}
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
