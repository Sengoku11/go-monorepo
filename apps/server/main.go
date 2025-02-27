// Example server to demonstrate the ability of monorepo to handle multiple applications.
package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/Sengoku11/go-monorepo/apps/server/internal/router"
	"github.com/Sengoku11/go-monorepo/pkg/bootstrap"
)

func main() {
	ctx, cancel, log := bootstrap.Default()
	defer cancel(nil)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	//nolint:mnd,exhaustruct
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router.New(log),
		ReadHeaderTimeout: 5 * time.Second,
		BaseContext:       func(_ net.Listener) context.Context { return ctx },
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Info("http server closed")
		}
	}()

	<-ctx.Done()
}
