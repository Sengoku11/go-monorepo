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
	alert.Alert(ctx, alerter.Event{Chan: alertchan.Test, Msg: "this is alert", Args: map[string]any{"key1": "value1"}})
	alert.Alert(ctx, alerter.Event{Chan: alertchan.Test, Msg: "this is alert", Args: nil})

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
