package http

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"
)

type contextKey string

const RequestIDKey contextKey = "requestID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			bytes := make([]byte, 8)
			_, _ = rand.Read(bytes)
			reqID = hex.EncodeToString(bytes)
		}

		w.Header().Set("X-Request-ID", reqID)

		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		reqID, _ := r.Context().Value(RequestIDKey).(string)

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		log.Printf("[%s] %s %s | Completed in %v", reqID, r.Method, r.URL.Path, duration)
	})
}
