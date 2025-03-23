// Package logger provides a customizable logging solution by wrapping
// one of the most popular logging libraries. It allows you to integrate
// structured logging with additional logic tailored to your application's needs.
package logger

import (
	"time"
)

// LogMethod is alias to Logger methods.
type LogMethod = func(message string, args ...any)

// BasicLogger represents core logging methods.
type BasicLogger interface {
	// Info logs a message at the info level. Optional key/value pairs can be provided for structured logging.
	//
	// Example:
	//
	//	log.Info("hello", "someKey", "string value").
	Info(message string, args ...any)

	// Warn logs a message at the warn level.
	Warn(message string, args ...any)

	// Error logs a message at the error level.
	Error(message string, args ...any)

	// Fatal logs a message at the fatal level and terminates the app with os.Exit(1)
	Fatal(message string, args ...any)

	// Panic logs a message at the panic level and then panics, which stops the ordinary flow of a goroutine.
	Panic(message string, args ...any)

	// Debug logs a message at the debug level.
	Debug(message string, args ...any)
}

// Controls that change the BasicLogger state.
type Controls interface {
	// EnableDebugMode to allow BasicLogger.Debug messages.
	EnableDebugMode()

	// DisableDebugMode to prevent BasicLogger.Debug messages.
	DisableDebugMode()
}

// Options for logging.
type Options interface {
	// WithRateLim applies rate limiting based on the message hash.
	WithRateLim(cooldown time.Duration, cb LogMethod) LogMethod

	// WithCallStack adds call stack to the log.
	//	Example: "stack": [{"func":"main.someFunc","file":"/path/to/file.go","line":42},...]
	WithCallStack(cb LogMethod) LogMethod
}

// Logger represents an extended interface of logging options and methods.
type Logger interface {
	BasicLogger
	Controls
	Options
}
