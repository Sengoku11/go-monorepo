// Package slack implements alerter for sending messages to messengers.
package slack

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/Sengoku11/go-monorepo/pkg/alerter/ratelim"
	"github.com/slack-go/slack"
)

var (
	errChannelNotFound = errors.New("channel not found")
	errTokenNotFound   = errors.New("SLACK_API_TOKEN environment variable not set")
)

// Alerter is a Slack implementation of alerter.Alerter interface.
type Alerter struct {
	client  *slack.Client
	msgByTS *ratelim.MessageByTS
}

// New instance of Alerter with "SLACK_API_TOKEN" env.
func New() (*Alerter, error) {
	token := os.Getenv("SLACK_API_TOKEN")
	if token == "" {
		return nil, errTokenNotFound
	}

	client := slack.New(os.Getenv("SLACK_API_TOKEN"))
	msgByTS := ratelim.NewMap()

	return &Alerter{
		client:  client,
		msgByTS: msgByTS,
	}, nil
}

// Alert given event to the Slack channel.
func (a *Alerter) Alert(ctx context.Context, event alerter.Event) error {
	hashed, err := ratelim.Hash(event.Message)
	if err != nil {
		return fmt.Errorf("alerter ratelimit: %w", err)
	}

	if ts, exists := a.msgByTS.Get(hashed); exists {
		if !ratelim.IsPastDue(ts, event.Options.RateLimit) {
			return nil
		}
	}

	a.msgByTS.AddNow(hashed)

	channel := getChannel(event.Options.ChannelSuffix)
	if channel == "" {
		return fmt.Errorf(`env SLACK_%s not defined: %w`, event.Options.ChannelSuffix, errChannelNotFound)
	}

	_, _, err = a.client.PostMessageContext(ctx, channel, slack.MsgOptionText(event.Message, false))
	if err != nil {
		a.msgByTS.Remove(hashed) // retry next time

		return fmt.Errorf("slack post message: %w", err)
	}

	return nil
}

func getChannel(envSuffix string) string {
	return os.Getenv("SLACK_" + envSuffix)
}

var _ alerter.Alerter = (*Alerter)(nil)
