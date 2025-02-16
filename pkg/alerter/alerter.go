// Package alerter provides Alerter implementations for sending data to messaging platforms
package alerter

import (
	"context"
	"time"
)

// Alerter defines the interface for alerting hooks.
type Alerter interface {
	// Alert processes the alert event within the provided context.
	// The context can be used for cancellation and timeout propagation.
	Alert(ctx context.Context, event Event) error
}

// Event encapsulates information for processing an alert.
type Event struct {
	Message    string
	Level      LogLevel
	Options    Options
	StructLogs []any
}

// LogLevel represents the type of logging alerts.
type LogLevel int

// LogLevel enum.
const (
	InfoLevel LogLevel = iota
	WarnLevel
	ErrorLevel
	FatalLevel
)

// Options contains additional configuration options for alerts.
type Options struct {
	RateLimit time.Duration // The minimum duration between alert messages (e.g., time.Hour or time.Minute).
	Channel   string        // The channel or chat name where the alert message should be pushed.
}

// DefaultOpts returns option with 1-hour rate limit.
func DefaultOpts(channel string) Options {
	return Options{
		RateLimit: time.Hour,
		Channel:   channel,
	}
}
