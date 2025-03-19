// Package fflag provides an OpenFeature compatible client and extends it capabilities.
package fflag

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/Sengoku11/go-monorepo/pkg/logger"
	flipt "github.com/open-feature/go-sdk-contrib/providers/flipt/pkg/provider"
	"github.com/open-feature/go-sdk/openfeature"
)

var (
	// ErrHostUndefined is returned when the FEATURE_TOGGLE_HOST environment variable is not set.
	ErrHostUndefined = errors.New("FEATURE_TOGGLE_HOST env is undefined")

	// ErrPortUndefined is returned when the FEATURE_TOGGLE_GRPC_PORT environment variable is not set.
	ErrPortUndefined = errors.New("FEATURE_TOGGLE_GRPC_PORT env is undefined")
)

// Client is wrapper over OpenFeature's client.
type Client struct {
	*openfeature.Client
	log logger.Logger
}

// NewFlipt returns wrapped openfeature.Client with flipt.io provider.
func NewFlipt(log logger.Logger, namespace string) (*Client, error) {
	host := os.Getenv("FEATURE_TOGGLE_HOST")
	if host == "" {
		return nil, ErrHostUndefined
	}

	port := os.Getenv("FEATURE_TOGGLE_GRPC_PORT")
	if port == "" {
		return nil, ErrPortUndefined
	}

	provider := flipt.NewProvider(
		flipt.WithAddress(fmt.Sprintf("%s:%s", host, port)),
		flipt.ForNamespace(namespace),
	)

	if err := openfeature.SetProviderAndWait(provider); err != nil {
		return nil, fmt.Errorf("failed to initialize provider: %w", err)
	}

	client := openfeature.NewClient(namespace)

	return &Client{client, log}, nil
}

// BooleanFlag holds the parameters for polling a boolean flag.
type BooleanFlag struct {
	Name         string
	DefaultValue bool
	EvalCtx      openfeature.EvaluationContext
	Options      []openfeature.Option
}

// PollBoolean continuously polls a boolean flag and runs the callback when its value changes.
func (c *Client) PollBoolean(ctx context.Context, ticker *time.Ticker, flag BooleanFlag, callback func(bool)) {
	defer ticker.Stop()

	val, err := c.BooleanValue(ctx, flag.Name, flag.DefaultValue, flag.EvalCtx, flag.Options...)
	if err != nil {
		c.log.Alert(ctx, "cannot poll the flag", alerter.DefaultOpts("FLAG_ERROR"), nil)
	}

	for {
		select {
		case <-ticker.C:
			newVal, err := c.BooleanValue(ctx, flag.Name, flag.DefaultValue, flag.EvalCtx, flag.Options...)
			if err != nil {
				c.log.Alert(ctx, "cannot poll the flag", alerter.DefaultOpts("FLAG_ERROR"), nil)

				continue
			}

			if newVal != val {
				val = newVal
				callback(newVal)
			}
		case <-ctx.Done():
			return
		}
	}
}
