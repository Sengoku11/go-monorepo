// Package alerter provides Alerter implementations for sending data to messaging platforms
package alerter

import (
	"context"
)

// Alerter defines the interface used in alerting hooks.
type Alerter interface {
	// Alert processes the alert event within the provided context.
	// The context can be used for cancellation and timeout propagation.
	Alert(ctx context.Context, channel Channel, message string, payload map[string]any) error
}
