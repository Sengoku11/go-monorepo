// Package alerter provides Client implementations for sending data to messaging platforms
package alerter

import (
	"context"
	"sync"

	"github.com/Sengoku11/go-monorepo/pkg/alerter/alertchan"
)

// Event encapsulates alert data.
type Event struct {
	Chan alertchan.Channel
	Msg  string
	Args map[string]any
}

// Alerter defines the interface of alerter used in Client hooks.
type Alerter interface {
	// Alert processes the alert event within the provided context.
	Alert(ctx context.Context, event Event) error

	// HandleError if failed to send an alert.
	HandleError(err error, event Event)
}

// The Client is used to send alert to all added Alerter.
type Client struct {
	hooks []Alerter
}

// New Client instance.
func New(hooks ...Alerter) *Client {
	return &Client{
		hooks: hooks,
	}
}

// Alert the message to hooked Alerter with AlertTimeout to complete.
func (a *Client) Alert(ctx context.Context, event Event, opts ...Option) {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, options.Timeout())
	defer cancel()

	var wg sync.WaitGroup
	for _, hook := range a.hooks {
		wg.Add(1)

		go func(hook Alerter) {
			defer wg.Done()

			if err := hook.Alert(timeoutCtx, event); err != nil {
				hook.HandleError(err, event)
			}
		}(hook)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
	case <-done:
	}
}
