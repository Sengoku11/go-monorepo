// Package fflag provides an OpenFeature compatible client and extends its capabilities.
package fflag

import (
	"context"
	"fmt"
	"time"

	"github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/Sengoku11/go-monorepo/pkg/logger"
	"github.com/open-feature/go-sdk/openfeature"
)

// Client is a wrapper over OpenFeature's client.
type Client struct {
	*openfeature.Client
	log logger.Logger
}

// BooleanFlag holds the parameters for polling a boolean flag.
type BooleanFlag struct {
	Name         string
	DefaultValue bool
	EvalCtx      openfeature.EvaluationContext
	Options      []openfeature.Option
}

// WatchBoolFlag continuously polls a boolean flag and runs the callback when its value changes.
// Not all providers in OpenFeature support event subscription on flag change.
func (c *Client) WatchBoolFlag(ctx context.Context, ticker *time.Ticker, flag BooleanFlag, callback func(bool)) {
	defer ticker.Stop()

	currentVal, err := c.fetchBooleanValue(ctx, flag)
	if err != nil {
		currentVal = flag.DefaultValue
	}

	for {
		select {
		case <-ticker.C:
			if newVal, err := c.fetchBooleanValue(ctx, flag); err == nil && newVal != currentVal {
				currentVal = newVal

				callback(newVal)
			}
		case <-ctx.Done():
			return
		}
	}
}

// fetchBooleanValue wraps the BooleanValue call, logs errors, and returns the current flag value.
func (c *Client) fetchBooleanValue(ctx context.Context, flag BooleanFlag) (bool, error) {
	val, err := c.BooleanValue(ctx, flag.Name, flag.DefaultValue, flag.EvalCtx, flag.Options...)
	if err != nil {
		msg := fmt.Sprintf("cannot fetch %s flag", flag.Name)

		c.log.Alert(ctx, msg, alerter.DefaultOpts("FLAG_ERROR"), map[string]any{"error": err.Error()})
		c.log.Error(msg, "error", err.Error())

		return flag.DefaultValue, fmt.Errorf("%s: %w", msg, err)
	}

	return val, nil
}
