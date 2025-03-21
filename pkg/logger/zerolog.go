package logger

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/Sengoku11/go-monorepo/pkg/ratelim"
	"github.com/rs/zerolog"
)

// ZerologLogger is an implementation of Logger interface using zerolog package.
// It leverages zero-allocation for efficient, structured logging.
type ZerologLogger struct {
	logger  *zerolog.Logger
	hooks   []alerter.Alerter
	debug   bool
	msgByTS *ratelim.MessageByTS
}

// NewZerologLogger creates a new instance of ZerologLogger.
func NewZerologLogger() *ZerologLogger {
	consoleWriter := zerolog.NewConsoleWriter()
	consoleWriter.TimeFormat = time.DateTime
	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()

	return &ZerologLogger{
		logger:  &logger,
		debug:   false,
		hooks:   []alerter.Alerter{},
		msgByTS: ratelim.NewMap(),
	}
}

// Info logs a message at the info level.
func (l *ZerologLogger) Info(message string, args ...any) {
	l.logger.Info().Fields(args).Msg(message)
}

// Debug logs a message at the info level if env DEBUG is set to true.
func (l *ZerologLogger) Debug(message string, args ...any) {
	if !l.debug {
		return
	}

	l.logger.Debug().Fields(args).Msg(message)
}

// Warn logs a message at the warn level.
func (l *ZerologLogger) Warn(message string, args ...any) {
	l.logger.Warn().Fields(args).Msg(message)
}

// Error logs a message at the error level.
func (l *ZerologLogger) Error(message string, args ...any) {
	l.logger.Error().Fields(args).Msg(message)
}

// Fatal logs a message at the fatal level and terminates the app with os.Exit(1).
func (l *ZerologLogger) Fatal(message string, args ...any) {
	l.logger.Fatal().Fields(args).Msg(message)
}

// Panic logs a message at the panic level and then panics, which stops the ordinary flow of a goroutine.
func (l *ZerologLogger) Panic(message string, args ...any) {
	l.logger.Panic().Fields(args).Msg(message)
}

// AddHook for sending alerts to messengers.
func (l *ZerologLogger) AddHook(hook alerter.Alerter) {
	l.hooks = append(l.hooks, hook)
}

// Alert the message to connected hooks.
func (l *ZerologLogger) Alert(ctx context.Context, message string, options alerter.Options, payload map[string]any) {
	event := alerter.Event{
		Message: message,
		Payload: payload,
		Options: options,
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, alertTimeout)
	defer cancel()

	var wg sync.WaitGroup
	for _, hook := range l.hooks {
		wg.Add(1)

		go func(h alerter.Alerter) {
			defer wg.Done()

			err := h.Alert(timeoutCtx, event)
			if errors.Is(err, context.DeadlineExceeded) {
				l.Error(
					ErrSendTimeout.Error(),
					"error", context.DeadlineExceeded.Error(),
					"channel", options.ChannelSuffix,
					"msg", message,
				)
			} else if err != nil {
				l.Error(
					ErrSendAlert.Error(),
					"error", err.Error(),
					"channel", options.ChannelSuffix,
					"msg", message,
				)
			}
		}(hook)
	}

	wgChan := make(chan any)
	go func() {
		wg.Wait()
		close(wgChan)
	}()

	select {
	case <-ctx.Done():
	case <-wgChan:
	}
}

// EnableDebugMode to log debug messages.
func (l *ZerologLogger) EnableDebugMode() {
	l.debug = true
}

// DisableDebugMode to avoid spamming debug messages.
func (l *ZerologLogger) DisableDebugMode() {
	l.debug = false
}

// WithRateLim applies rate limiting to a logging callback based on the message hash.
func (l *ZerologLogger) WithRateLim(cooldown time.Duration, callback LogMethod) LogMethod {
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
func (l *ZerologLogger) WithCallStack(callback LogMethod) LogMethod {
	return func(message string, args ...any) {
		if pc, _, line, ok := runtime.Caller(1); ok {
			funcNameFull := runtime.FuncForPC(pc).Name()

			args = append(args, "stack", fmt.Sprintf("%ss#%d", funcNameFull, line))
		}

		callback(message, args...)
	}
}

var _ Logger = (*ZerologLogger)(nil)
