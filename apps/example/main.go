// Example application that you can use as a playground
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/Sengoku11/go-monorepo/pkg/alerter/channels"
	"github.com/Sengoku11/go-monorepo/pkg/alerter/slack"
	"github.com/Sengoku11/go-monorepo/pkg/loadenv"
	"github.com/Sengoku11/go-monorepo/pkg/logger"
)

var errTest = errors.New("this is an error")

func main() {
	ctx := context.Background()

	if err := loadenv.Local(); err != nil {
		panic("cannot start locally: " + err.Error())
	}

	var log logger.Logger = logger.NewZerologLogger()

	log.Info("starting app")

	log.Debug("this message appears when env DEBUG is true", "debug", os.Getenv("DEBUG"))

	err2 := fmt.Errorf("error2 says hello: %w", errTest)
	err3 := fmt.Errorf("error3 wraps everything up: %w", err2)
	log.Error("intentionally failed", "error", err3)

	log.AddHook(slack.New())
	log.AlertInfo(ctx, "this is alert", alerter.DefaultOpts(channels.TEST), "someData", 1)
}
