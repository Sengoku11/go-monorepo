// Package alerter provides Alerter implementations for sending data to messaging platforms
package alerter

import (
	"context"
	"time"
)

// Alerter defines the interface used in alerting hooks.
type Alerter interface {
	// Alert processes the alert event within the provided context.
	// The context can be used for cancellation and timeout propagation.
	Alert(ctx context.Context, event Event) error
}

// Event info for processing an alert.
type Event struct {
	Message string
	Payload map[string]any
	Channel Channel // The suffix for the channel's environment variable, e.g., "ERROR".
	Options Options
}

// Options for alerts.
type Options struct {
	RateLimit time.Duration // The minimum duration between alert messages (e.g., time.Hour or time.Minute).
	// ChannelSuffix alertchannel.Channel
}

// DefaultOpts returns option with 1-hour rate limit.
func DefaultOpts() Options {
	return Options{
		RateLimit: time.Hour,
		// ChannelSuffix: channel,
	}
}
