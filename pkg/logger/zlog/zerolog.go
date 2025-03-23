// Package zlog implements logger.Logger interface via zerolog logger.
package zlog

import (
	"runtime"
	"strconv"
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
func New(options ...func(w *zerolog.ConsoleWriter)) *Logger {
	consoleWriter := zerolog.NewConsoleWriter(options...)
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

// WithRateLimit applies rate limiting to a logging callback based on the message hash.
func (l *Logger) WithRateLimit(cooldown time.Duration, callback logger.LogMethod) logger.LogMethod {
	return func(message string, args ...any) {
		hashed, err := ratelim.Hash(message)
		if err != nil {
			l.Error("failed to hash message in WithRateLimit", "error", err.Error(), "message", message)

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

// WithCallStack adds a call stack to the log. Each object in the slice has the keys: "func", "file", and "line".
//
//	Example: "stack": [{"func":"main.someFunc","file":"/path/to/file.go","line":42},...]
func (l *Logger) WithCallStack(callback logger.LogMethod) logger.LogMethod {
	return func(message string, args ...any) {
		const (
			skipFrames  = 2
			optimalSize = 10
		)

		pcs := make([]uintptr, optimalSize)
		n := runtime.Callers(skipFrames, pcs)
		frames := runtime.CallersFrames(pcs[:n])

		stack := make([]map[string]string, 0, optimalSize)

		for {
			frame, more := frames.Next()
			stack = append(stack, map[string]string{
				"file": frame.File,
				"func": frame.Function,
				"line": strconv.Itoa(frame.Line),
			})

			if !more {
				break
			}
		}

		args = append(args, "stack", stack)

		callback(message, args...)
	}
}

var (
	_ logger.Logger   = (*Logger)(nil)
	_ logger.Basic    = (*Logger)(nil)
	_ logger.Advanced = (*Logger)(nil)
	_ logger.Options  = (*Logger)(nil)
	_ logger.Controls = (*Logger)(nil)
)
