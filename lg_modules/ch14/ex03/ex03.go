package main

import (
	"context"
	"fmt"
	"net/http"
)

type levelKey int

const (
	_ levelKey = iota
	key
)

type Level string

const (
	Debug Level = "debug"
	Info  Level = "info"
)

func ContextWithLevel(ctx context.Context, level Level) context.Context {
	return context.WithValue(ctx, key, level)
}

func LevelFromContext(ctx context.Context) (Level, bool) {
	level, ok := ctx.Value(key).(Level)
	return level, ok
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		level := req.URL.Query().Get("log_level")
		ctx := ContextWithLevel(req.Context(), Level(level))
		req = req.WithContext(ctx)
		h.ServeHTTP(rw, req)
	})
}
func Log(ctx context.Context, level Level, message string) {
	var inLevel Level
	inLevel, ok := LevelFromContext(ctx)
	if !ok {
		return
	}
	if level == Debug && inLevel == Debug {
		fmt.Println(message)
	}
	if level == Info && (inLevel == Debug || inLevel == Info) {
		fmt.Println(message)
	}
}
