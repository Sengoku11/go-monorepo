// Leveraging the monorepo's full capabilities, this server exposes the complete set of error codes
// defined in errcode.AllCodes, essential for third-party consumers to integrate and debug effectively.
package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/Sengoku11/go-monorepo/apps/errdoc/internal/router"
	"github.com/Sengoku11/go-monorepo/pkg/bootstrap"
)

func main() {
	ctx, cancel, log := bootstrap.Default()
	defer cancel(nil)

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3003"
	}

	mux := router.New(log)

	//nolint:mnd,exhaustruct
	srv := http.Server{
		Addr:              fmt.Sprintf("%s:%s", host, port),
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
		BaseContext:       func(_ net.Listener) context.Context { return ctx },
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error(fmt.Sprintf("server ListenAndServe error: %v", err))
		} else {
			log.Info("server closed")
		}
	}()

	<-ctx.Done()
	log.Info("terminating...")

	//nolint:mnd
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error(fmt.Sprintf("error during shutdown: %v", err))
	} else {
		log.Info("shutdown complete")
	}
}
