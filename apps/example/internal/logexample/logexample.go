// Package logexample demonstrates logger.Logger capabilities.
package logexample

import (
	"errors"
	"fmt"
	"os"

	"github.com/Sengoku11/go-monorepo/pkg/logger"
)

var errTest = errors.New("this is an error")

// Run example.
func Run(log logger.Logger) {
	log.Info("starting app")
	log.Debug("this message appears when env DEBUG is true", "debug", os.Getenv("DEBUG"))

	err2 := fmt.Errorf("error2 says hello: %w", errTest)
	err3 := fmt.Errorf("error3 wraps everything up: %w", err2)
	log.Error("intentionally failed", "error", err3)
}
