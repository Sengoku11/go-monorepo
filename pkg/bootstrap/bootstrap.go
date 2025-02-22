// Package bootstrap provides common initialization logic for services.
package bootstrap

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/Sengoku11/go-monorepo/pkg/loadenv"
	"github.com/Sengoku11/go-monorepo/pkg/logger"
)

// Default creates a context that is canceled on SIGINT or SIGTERM, loads environment variables
// if running locally, and returns a logger instance. It panics if environment loading fails
// in a local environment.
//
//nolint:ireturn
func Default() (context.Context, context.CancelFunc, logger.Logger) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	if err := loadenv.Local(); err != nil {
		panic("cannot start locally: " + err.Error())
	}

	log := logger.NewZerologLogger()

	return ctx, cancel, log
}
