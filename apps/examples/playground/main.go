// Example application that you can use as a playground
package main

import (
	"sync"

	"github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/Sengoku11/go-monorepo/pkg/alerter/alertchan"
	"github.com/Sengoku11/go-monorepo/pkg/alerter/slack"
	"github.com/Sengoku11/go-monorepo/pkg/bootstrap"
)

func main() {
	ctx, cancel, log := bootstrap.Default()
	defer cancel(nil)

	slackAlerter, err := slack.New(log)
	if err != nil {
		log.Panic("cannot create slack alerter", "error", err)
	}

	alert := alerter.New(slackAlerter)

	alert.Alert(ctx, alertchan.Test, "this is alert", map[string]any{"key1": "value1"})
	alert.Alert(ctx, alertchan.Test, "this is alert", map[string]any{"key2": "value2"})

	log.Warn("only first alert sent: default 1h rate limit for identical messages, even with different payloads")

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		log.Info("emulating an ongoing background process; press ctrl+c to shutdown gracefully")
		<-ctx.Done()
	}()

	wg.Wait()
	log.Info("terminating app")
}
