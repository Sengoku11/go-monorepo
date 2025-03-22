// Package alerter provides Client implementations for sending data to messaging platforms
package alerter

import (
	"context"
	"sync"
	"time"

	"github.com/Sengoku11/go-monorepo/pkg/alerter/alertchan"
)

// AlertTimeout defines how long to wait for alert to be sent.
const AlertTimeout = time.Second * 2

// Alerter defines the interface of alerter used in Client hooks.
type Alerter interface {
	// Alert processes the alert event within the provided context.
	Alert(ctx context.Context, channel alertchan.Channel, message string, payload map[string]any) error

	// HandleError if failed to send an alert.
	HandleError(err error, message string, payload map[string]any)
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
func (a *Client) Alert(ctx context.Context, channel alertchan.Channel, message string, payload map[string]any) {
	timeoutCtx, cancel := context.WithTimeout(ctx, AlertTimeout)
	defer cancel()

	var wg sync.WaitGroup
	for _, hook := range a.hooks {
		wg.Add(1)

		go func(hook Alerter) {
			defer wg.Done()

			if err := hook.Alert(timeoutCtx, channel, message, payload); err != nil {
				hook.HandleError(err, message, payload)
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
