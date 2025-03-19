// Example application that you can use as a playground
package main

import (
	"sync"

	"github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/Sengoku11/go-monorepo/pkg/alerter/channels"
	"github.com/Sengoku11/go-monorepo/pkg/alerter/slack"
	"github.com/Sengoku11/go-monorepo/pkg/bootstrap"
)

func main() {
	ctx, cancel, log := bootstrap.Default()
	defer cancel(nil)

	if slackAlerter, err := slack.New(); err == nil {
		log.AddHook(slackAlerter)
	} else {
		log.Panic("cannot create slack alerter", "error", err)
	}

	log.Alert(ctx, "this is alert", alerter.DefaultOpts(channels.TEST), map[string]any{"key1": "value1"})
	log.Alert(ctx, "this is alert", alerter.DefaultOpts(channels.TEST), map[string]any{"key2": "value2"})
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
