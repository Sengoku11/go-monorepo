// Package router provides chi.Mux.
package router

import (
	"net/http"

	"github.com/Sengoku11/go-monorepo/pkg/logger"
	"github.com/Sengoku11/go-monorepo/pkg/middleware"
	"github.com/go-chi/chi/v5"
	chimware "github.com/go-chi/chi/v5/middleware"
)

// New router.
func New(log logger.Logger) *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(chimware.RequestID)
	mux.Use(chimware.Recoverer)
	mux.Use(middleware.LogRequest(log))

	mux.HandleFunc("/hello", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("Hello World"))
		if err != nil {
			log.Error("failed to say hello")
		}
	})

	return mux
}
