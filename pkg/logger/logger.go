// Package logger provides a customizable logging solution by wrapping
// one of the most popular logging libraries. It allows you to integrate
// structured logging with additional logic tailored to your application's needs.
package logger

import (
	"time"

	"github.com/rs/zerolog"
)

// Logger defines the minimal logging interface. It can be extended with more methods as needed.
type Logger interface {
	// Info logs a message at the info level. Optional key/value pairs can be provided for structured logging.
	//
	// Example:
	//
	//	log.Info("hello", "someKey", "string value").
	Info(message string, args ...any)

	Fatal(message string, args ...any)
}

// ZerologLogger is an implementation of Logger interface using zerolog package.
// It leverages zero-allocation for efficient, structured logging.
type ZerologLogger struct {
	logger *zerolog.Logger
}

// NewZerologLogger creates a new instance of ZerologLogger.
func NewZerologLogger() *ZerologLogger {
	consoleWriter := zerolog.NewConsoleWriter()
	consoleWriter.TimeFormat = time.DateTime

	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()

	return &ZerologLogger{logger: &logger}
}

// Info logs a message at the info level.
func (l *ZerologLogger) Info(message string, args ...any) {
	l.logger.Info().Fields(args).Msg(message)
}

// Fatal logs a message at the fatal level and terminates the app with os.Exit(1).
func (l *ZerologLogger) Fatal(message string, args ...any) {
	l.logger.Fatal().Fields(args).Msg(message)
}

var _ Logger = (*ZerologLogger)(nil)
