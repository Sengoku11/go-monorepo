// Package slack implements alerter for sending messages to messengers.
package slack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/Sengoku11/go-monorepo/pkg/ratelim"
	"github.com/slack-go/slack"
)

var errTokenNotFound = errors.New("SLACK_API_TOKEN environment variable not set")

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

// AlertWithRateLimit given event to the Slack channel.
func (a *Alerter) AlertWithRateLimit(
	ctx context.Context,
	cooldown time.Duration,
	channel alerter.Channel,
	message string,
	payload map[string]any,
) error {
	hashed, err := ratelim.Hash(message)
	if err != nil {
		return fmt.Errorf("alerter ratelimit: %w", err)
	}

	if ts, exists := a.msgByTS.Get(hashed); exists {
		if !ratelim.IsPastDue(ts, cooldown) {
			return nil
		}
	}

	a.msgByTS.AddNow(hashed)

	if err := a.Alert(ctx, channel, message, payload); err != nil {
		a.msgByTS.Remove(hashed) // retry next time

		return err
	}

	return nil
}

// Alert given event to the Slack channel.
func (a *Alerter) Alert(ctx context.Context, suffix alerter.Channel, message string, payload map[string]any) error {
	channel, err := alerter.FromENV("SLACK", suffix)
	if err != nil {
		return fmt.Errorf(`SLACK_* env is not defined: %w`, err)
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("alerter payload marshal: %w", err)
	}

	payloadString := string(payloadBytes)
	msg := message + "\n" + payloadString

	_, _, err = a.client.PostMessageContext(ctx, channel, slack.MsgOptionText(msg, false))
	if err != nil {
		return fmt.Errorf("slack post message: %w", err)
	}

	return nil
}

var _ alerter.Alerter = (*Alerter)(nil)
