// Package fflag provides an OpenFeature compatible client and extends its capabilities.
package fflag

import (
	"context"
	"fmt"
	"time"

	"github.com/open-feature/go-sdk/openfeature"
)

// Client is a wrapper over OpenFeature's client.
type Client struct {
	openfeature.IClient
}

// New creates and returns a Client instance using a unique identifier for this client (namespace).
//
// For OpenFeature-compatible providers see https://github.com/open-feature/go-sdk-contrib/tree/main/providers.
func New(namespace string, provider openfeature.FeatureProvider) (*Client, error) {
	if err := openfeature.SetProviderAndWait(provider); err != nil {
		return nil, fmt.Errorf("failed to initialize provider: %w", err)
	}

	client := openfeature.NewClient(namespace)

	return &Client{client}, nil
}

// BooleanFlag holds the parameters for polling a boolean flag.
type BooleanFlag struct {
	Name         string
	DefaultValue bool
	EvalCtx      openfeature.EvaluationContext
	Options      []openfeature.Option
}

// BoolCallback function to be called on bool flag update.
type BoolCallback func(bool, error)

// WatchBoolFlag continuously polls a boolean flag and runs the callback when its value changes.
// Not all providers in OpenFeature support event subscription on flag change.
//
// Helpful if you want to implement a kill-switch.
func (c *Client) WatchBoolFlag(ctx context.Context, flag BooleanFlag, ticker *time.Ticker, callback BoolCallback) {
	defer ticker.Stop()

	// currentVal, err := c.fetchBooleanValue(ctx, flag)
	// if err != nil {
	//	currentVal = flag.DefaultValue
	//
	//	callback(flag.DefaultValue, err)
	//}

	currentVal := flag.DefaultValue

	for {
		select {
		case <-ticker.C:
			newVal, err := c.fetchBooleanValue(ctx, flag)
			if err != nil {
				callback(newVal, err)
			} else if newVal != currentVal {
				currentVal = newVal

				callback(newVal, nil)
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

		return flag.DefaultValue, fmt.Errorf("%s: %w", msg, err)
	}

	return val, nil
}
