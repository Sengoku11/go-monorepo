// Package middleware provides various decorators for http.Handler.
package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sengoku11/go-monorepo/pkg/logger"
	chimware "github.com/go-chi/chi/v5/middleware"
)

// LogRequest returns a middleware that logs HTTP requests upon completion.
func LogRequest(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			defer func() {
				duration := time.Since(start).Milliseconds()
				msg, fields := buildLogEntry(r, duration)
				log.Info(msg, fields...)
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// buildLogEntry constructs a log message and fields for a request.
func buildLogEntry(r *http.Request, duration int64) (string, []any) {
	reqID := chimware.GetReqID(r.Context())

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	// Build the URL string.
	url := fmt.Sprintf("%s://%s%s %s", scheme, r.Host, r.RequestURI, r.Proto)
	msg := fmt.Sprintf("%s %s", r.Method, url)

	// Return the message and structured fields.
	return msg, []any{
		"method", r.Method,
		"path", r.URL.Path,
		"query", r.URL.RawQuery,
		"request_id", reqID,
		"duration_ms", duration,
	}
}
