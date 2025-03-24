// Package slack implements alerter for sending messages to messengers.
package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/Sengoku11/go-monorepo/pkg/alerter/alertchan"
	"github.com/Sengoku11/go-monorepo/pkg/logger"
	"github.com/Sengoku11/go-monorepo/pkg/ratelim"
	slacksdk "github.com/slack-go/slack"
)

// MessengerPrefix contains a prefix used to find Slack channels from env.
const MessengerPrefix = "SLACK"

// Client interface for decoupling and testing (mocks) purposes.
type Client interface {
	PostMessageContext(ctx context.Context, channelID string, options ...slacksdk.MsgOption) (string, string, error)
}

// NewClient returns slack Client.
func NewClient(token string) *slacksdk.Client {
	client := slacksdk.New(token)

	return client
}

// Alerter is a Slack implementation of alerter.Alerter.
type Alerter struct {
	client  Client
	log     logger.Basic
	msgByTS *ratelim.MessageByTS
}

// New instance of Alerter.
func New(token string, opts ...Option) *Alerter {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	if options.client == nil {
		options.client = NewClient(token)
	}

	return &Alerter{
		log:     options.logger,
		client:  options.client,
		msgByTS: ratelim.NewMap(),
	}
}

// WithRateLimit given event to the Slack channel.
func (a *Alerter) WithRateLimit(
	ctx context.Context,
	cooldown time.Duration,
	event alerter.Event,
	opts alerter.Options,
) error {
	hashed, err := ratelim.Hash(event.Msg)
	if err != nil {
		return fmt.Errorf("alerter ratelimit: %w", err)
	}

	if ts, exists := a.msgByTS.Get(hashed); exists {
		if !ratelim.IsPastDue(ts, cooldown) {
			return nil
		}
	}

	a.msgByTS.AddNow(hashed)

	if err := a.Alert(ctx, event, opts); err != nil {
		a.msgByTS.Remove(hashed) // retry next time

		return err
	}

	return nil
}

// Alert given event to the Slack channel.
func (a *Alerter) Alert(ctx context.Context, event alerter.Event, _ alerter.Options) error {
	channel, err := alertchan.FromEnv(MessengerPrefix, event.Chan)
	if err != nil {
		return fmt.Errorf(`alert channel is undefined: %w`, err)
	}

	payloadBytes, err := json.Marshal(event.Args)
	if err != nil {
		return fmt.Errorf("alerter payload marshal: %w", err)
	}

	payloadString := string(payloadBytes)
	msg := event.Msg + "\n" + payloadString

	_, _, err = a.client.PostMessageContext(ctx, channel, slacksdk.MsgOptionText(msg, false))
	if err != nil {
		return fmt.Errorf("slack post message: %w", err)
	}

	return nil
}

// HandleError if failed to send an alert.
func (a *Alerter) HandleError(err error, event alerter.Event) {
	a.log.Error(err.Error(), "message", event.Msg, event.Args)
}

var _ alerter.Alerter = (*Alerter)(nil)
