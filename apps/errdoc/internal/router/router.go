// Package router provides configured and ready-to-use chi.Mux router.
package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Sengoku11/go-monorepo/pkg/errcode"
	"github.com/Sengoku11/go-monorepo/pkg/logger"
	"github.com/Sengoku11/go-monorepo/pkg/middleware"
	"github.com/go-chi/chi/v5"
	chimware "github.com/go-chi/chi/v5/middleware"
)

// New router.
func New(log logger.Logger) *chi.Mux {
	mux := chi.NewRouter()
	codeToName := sliceToMap(errcode.AllCodes)

	mux.Use(chimware.RequestID)
	mux.Use(middleware.LogRequest(log))

	mux.Get("/errors", getAllErrors(codeToName))
	mux.Get("/errors/{code}", getErrorByCode(codeToName))

	return mux
}

func sliceToMap(errCodes []errcode.Code) map[int]string {
	codeToName := make(map[int]string, len(errCodes))
	for _, code := range errCodes {
		codeToName[int(code)] = code.String()
	}

	return codeToName
}

func getAllErrors(codeToName map[int]string) http.HandlerFunc {
	bytes, err := json.Marshal(codeToName)
	if err != nil {
		panic(fmt.Sprintf("marshal codeToName map: %v", err))
	}

	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		_, _ = w.Write(bytes)
	}
}

func getErrorByCode(codeToName map[int]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		codeParam := chi.URLParam(r, "code")

		code, err := strconv.Atoi(codeParam)
		if err != nil {
			http.Error(w, "invalid code", http.StatusBadRequest)

			return
		}

		name, exist := codeToName[code]
		if !exist {
			http.Error(w, fmt.Sprintf("code %d not found", code), http.StatusNotFound)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"name": name})
	}
}
