package main

import (
	"context"
	"net/http"
	"time"
)

func ContextWithTimeout(milliseconds int) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			ctx, cancel := context.WithTimeout(ctx, time.Duration(milliseconds)*time.Millisecond)
			defer cancel()
			req = req.WithContext(ctx)
			h.ServeHTTP(rw, req)
		})
	}
}
