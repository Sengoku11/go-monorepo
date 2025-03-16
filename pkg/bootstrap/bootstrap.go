// Package bootstrap provides common initialization logic for services.
package bootstrap

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Sengoku11/go-monorepo/pkg/loadenv"
	"github.com/Sengoku11/go-monorepo/pkg/logger"
)

// Default creates a context that is canceled on SIGINT or SIGTERM, loads environment variables
// if running locally, and returns a logger instance. It panics if environment loading fails
// in a local environment.
//
//nolint:ireturn
func Default() (context.Context, context.CancelCauseFunc, logger.Logger) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	if err := loadenv.Local(); err != nil {
		stop()
		panic("cannot start locally: " + err.Error())
	}

	log := logger.NewZerologLogger()
	ctxWithCause, cancelCause := context.WithCancelCause(ctx)

	// Set debug mode on if enabled
	if debug, err := strconv.ParseBool(os.Getenv("DEBUG")); err != nil {
		log.Warn("DEBUG environment variable is not set, defaulting to false")
	} else if debug {
		log.EnableDebugMode()
	}

	log.Info("starting the application")

	return ctxWithCause, cancelCause, log
}
