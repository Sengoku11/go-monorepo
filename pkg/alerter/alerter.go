// Package alerter provides Client implementations for sending data to messaging platforms
package alerter

import (
	"context"
	"errors"
	"sync"

	"github.com/Sengoku11/go-monorepo/pkg/alerter/alertchan"
	"github.com/Sengoku11/go-monorepo/pkg/logger"
)

// Error that is logged at the error level when the alert hook fails.
var (
	ErrSendAlert   = errors.New("failed to send alert")
	ErrSendTimeout = errors.New("failed to send in time")
)

// Alerter defines the interface used in alerting hooks.
type Alerter interface {
	// Alert processes the alert event within the provided context.
	// The context can be used for cancellation and timeout propagation.
	Alert(ctx context.Context, channel alertchan.Channel, message string, payload map[string]any) error

	// HandleError if failed to send an alert.
	HandleError(err error, message string, payload map[string]any)
}

// The Client is used to send alert to all added Alerter.
type Client struct {
	mu    sync.RWMutex
	hooks []Alerter
}

// New Client instance.
func New() *Client {
	return &Client{
		mu:    sync.RWMutex{},
		hooks: []Alerter{},
	}
}

// AddHook appends a new Alerter to alerting list.
func (a *Client) AddHook(hook Alerter) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.hooks = append(a.hooks, hook)
}

// Alert the message to connected hooks.
func (a *Client) Alert(ctx context.Context, channel alertchan.Channel, message string, payload map[string]any) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	timeoutCtx, cancel := context.WithTimeout(ctx, logger.AlertTimeout)
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
