// Package logger provides a customizable logging solution by wrapping
// one of the most popular logging libraries. It allows you to integrate
// structured logging with additional logic tailored to your application's needs.
package logger

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/rs/zerolog"
)

const alertTimeout = time.Second * 2

// Error that is logged at the error level when the alert hook fails.
var (
	ErrSendAlert   = errors.New("failed to send alert")
	ErrSendTimeout = errors.New("failed to send in time")
)

// Logger defines the minimal logging interface. It can be extended with more methods as needed.
type Logger interface {
	// AddHook for sending alerts to messengers.
	AddHook(hook alerter.Alerter)

	// Info logs a message at the info level. Optional key/value pairs can be provided for structured logging.
	//
	// Example:
	//
	//	log.Info("hello", "someKey", "string value").
	Info(message string, args ...any)

	// Debug logs a message at the info level if env "DEBUG" is set to true.
	Debug(message string, args ...any)

	// Warn logs a message at the warn level.
	Warn(message string, args ...any)

	// Error logs a message at the error level.
	Error(message string, args ...any)

	// Fatal logs a message at the fatal level and terminates the app with os.Exit(1)
	Fatal(message string, args ...any)

	// Panic logs a message at the panic level and then panics, which stops the ordinary flow of a goroutine.
	Panic(message string, args ...any)

	// Alert the message to connected hooks.
	Alert(ctx context.Context, message string, options alerter.Options, payload map[string]any)

	// EnableDebugMode to log debug messages.
	EnableDebugMode()

	// DisableDebugMode to avoid spamming debug messages.
	DisableDebugMode()
}

// ZerologLogger is an implementation of Logger interface using zerolog package.
// It leverages zero-allocation for efficient, structured logging.
type ZerologLogger struct {
	logger *zerolog.Logger
	hooks  []alerter.Alerter
	debug  bool
}

// NewZerologLogger creates a new instance of ZerologLogger.
func NewZerologLogger() *ZerologLogger {
	consoleWriter := zerolog.NewConsoleWriter()
	consoleWriter.TimeFormat = time.DateTime
	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()

	return &ZerologLogger{
		logger: &logger,
		debug:  false,
		hooks:  []alerter.Alerter{},
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

// AddHook for sending alerts to messengers.
func (l *ZerologLogger) AddHook(hook alerter.Alerter) {
	l.hooks = append(l.hooks, hook)
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

var _ Logger = (*ZerologLogger)(nil)
