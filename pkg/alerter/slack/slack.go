// Package slack implements alerter for sending messages to messengers.
package slack

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/slack-go/slack"
)

var (
	errChannelNotFound = errors.New("channel not found")
	errTokenNotFound   = errors.New("SLACK_API_TOKEN environment variable not set")
)

// Alerter is a Slack implementation of alerter.Alerter interface.
type Alerter struct {
	client *slack.Client
}

// New instance of Alerter with "SLACK_API_TOKEN" env.
func New() (*Alerter, error) {
	token := os.Getenv("SLACK_API_TOKEN")
	if token == "" {
		return nil, errTokenNotFound
	}

	client := slack.New(os.Getenv("SLACK_API_TOKEN"))

	return &Alerter{
		client: client,
	}, nil
}

// Alert given event to the Slack channel.
func (a *Alerter) Alert(ctx context.Context, event alerter.Event) error {
	text := event.Message

	channel := getChannel(event.Options.ChannelSuffix)
	if channel == "" {
		return fmt.Errorf(`env SLACK_%s not defined: %w`, event.Options.ChannelSuffix, errChannelNotFound)
	}

	_, _, err := a.client.PostMessageContext(ctx, channel, slack.MsgOptionText(text, false))
	if err != nil {
		return fmt.Errorf("slack post message: %w", err)
	}

	return nil
}

func getChannel(envSuffix string) string {
	return os.Getenv("SLACK_" + envSuffix)
}

var _ alerter.Alerter = (*Alerter)(nil)
