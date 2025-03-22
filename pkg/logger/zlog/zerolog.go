// Package zlog implements logger.Logger interface via zerolog logger.
package zlog

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/Sengoku11/go-monorepo/pkg/logger"
	"github.com/Sengoku11/go-monorepo/pkg/ratelim"
	"github.com/rs/zerolog"
)

// Logger is an implementation of logger.Logger interface using zerolog package.
// It leverages zero-allocation for efficient, structured logging.
type Logger struct {
	logger  *zerolog.Logger
	hooks   []alerter.Alerter
	debug   *atomic.Bool
	msgByTS *ratelim.MessageByTS
}

// New creates a new instance of Logger.
func New() *Logger {
	consoleWriter := zerolog.NewConsoleWriter()
	consoleWriter.TimeFormat = time.DateTime
	log := zerolog.New(consoleWriter).With().Timestamp().Logger()

	return &Logger{
		logger:  &log,
		debug:   new(atomic.Bool),
		hooks:   []alerter.Alerter{},
		msgByTS: ratelim.NewMap(),
	}
}

// Info logs a message at the info level.
func (l *Logger) Info(message string, args ...any) {
	l.logger.Info().Fields(args).Msg(message)
}

// Warn logs a message at the warn level.
func (l *Logger) Warn(message string, args ...any) {
	l.logger.Warn().Fields(args).Msg(message)
}

// Error logs a message at the error level.
func (l *Logger) Error(message string, args ...any) {
	l.logger.Error().Fields(args).Msg(message)
}

// Fatal logs a message at the fatal level and terminates the app with os.Exit(1).
func (l *Logger) Fatal(message string, args ...any) {
	l.logger.Fatal().Fields(args).Msg(message)
}

// Panic logs a message at the panic level and then panics, which stops the ordinary flow of a goroutine.
func (l *Logger) Panic(message string, args ...any) {
	l.logger.Panic().Fields(args).Msg(message)
}

// Debug logs a message at the info level if env DEBUG is set to true.
func (l *Logger) Debug(message string, args ...any) {
	if !l.debug.Load() {
		return
	}

	l.logger.Debug().Fields(args).Msg(message)
}

// EnableDebugMode to log debug messages.
func (l *Logger) EnableDebugMode() {
	l.debug.Store(true)
}

// DisableDebugMode to avoid spamming debug messages.
func (l *Logger) DisableDebugMode() {
	l.debug.Store(false)
}

// WithRateLim applies rate limiting to a logging callback based on the message hash.
func (l *Logger) WithRateLim(cooldown time.Duration, callback logger.LogMethod) logger.LogMethod {
	return func(message string, args ...any) {
		hashed, err := ratelim.Hash(message)
		if err != nil {
			l.Error("failed to hash message in WithRateLim", "error", err.Error(), "message", message)

			// Calling without limiting
			callback(message, args...)

			return
		}

		if ts, exists := l.msgByTS.Get(hashed); exists {
			if !ratelim.IsPastDue(ts, cooldown) {
				return
			}
		}

		l.msgByTS.AddNow(hashed)
		callback(message, args...)
	}
}

// WithCallStack adds a call stack to the log.
func (l *Logger) WithCallStack(callback logger.LogMethod) logger.LogMethod {
	return func(message string, args ...any) {
		if pc, _, line, ok := runtime.Caller(1); ok {
			funcNameFull := runtime.FuncForPC(pc).Name()

			args = append(args, "stack", fmt.Sprintf("%ss#%d", funcNameFull, line))
		}

		callback(message, args...)
	}
}

var _ logger.Logger = (*Logger)(nil)
