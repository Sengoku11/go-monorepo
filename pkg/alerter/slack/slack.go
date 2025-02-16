// Package slack implements alerter for sending messages to messengers.
package slack

import (
	"context"
	"fmt"
	"os"

	"github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/slack-go/slack"
)

// Alerter is a Slack implementation.
type Alerter struct {
	client *slack.Client
}

// New Alerter.
func New() *Alerter {
	client := slack.New(os.Getenv("SLACK_API_TOKEN"))

	return &Alerter{
		client: client,
	}
}

// Alert event.
func (a *Alerter) Alert(ctx context.Context, event alerter.Event) error {
	// TODO: add fields
	text := event.Message

	_, _, err := a.client.PostMessageContext(ctx, event.Options.Channel, slack.MsgOptionText(text, false))
	if err != nil {
		return fmt.Errorf("slack post message: %w", err)
	}

	return nil
}

var _ alerter.Alerter = (*Alerter)(nil)
