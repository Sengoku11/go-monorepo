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

// ErrChannelNotFound -- you probably haven't created SLACK_ env for a given prefix.
var ErrChannelNotFound = errors.New("channel not found")

// Alerter is a Slack implementation of alerter.Alerter interface.
type Alerter struct {
	client *slack.Client
}

// New constructs Alerter with "SLACK_API_TOKEN" env.
func New() *Alerter {
	client := slack.New(os.Getenv("SLACK_API_TOKEN"))

	return &Alerter{
		client: client,
	}
}

// Alert given event to the Slack channel.
func (a *Alerter) Alert(ctx context.Context, event alerter.Event) error {
	text := event.Message

	channel := getChannel(event.Options.ChannelSuffix)
	if channel == "" {
		return ErrChannelNotFound
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
