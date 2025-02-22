// Package alertexample demonstrates alerter.Alerter capabilities.
package alertexample

import (
	"context"

	"github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/Sengoku11/go-monorepo/pkg/alerter/channels"
	"github.com/Sengoku11/go-monorepo/pkg/alerter/slack"
	"github.com/Sengoku11/go-monorepo/pkg/logger"
)

// Run example.
func Run(ctx context.Context, log logger.Logger) {
	slackAlerter, err := slack.New()
	if err != nil {
		log.Panic("cannot create slack alerter", "error", err)
	}

	log.AddHook(slackAlerter)

	payload1 := map[string]any{"key1": "value1"}
	payload2 := map[string]any{"key2": "value2"}

	log.Alert(ctx, "this is alert", alerter.DefaultOpts(channels.TEST), payload1)
	log.Alert(ctx, "this is alert", alerter.DefaultOpts(channels.TEST), payload2)
	log.Warn("only one alert sent: default 1h rate limit for identical messages, even with different payloads")
}
