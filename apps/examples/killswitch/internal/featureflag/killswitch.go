package featureflag

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/Sengoku11/go-monorepo/pkg/fflag"
	"github.com/Sengoku11/go-monorepo/pkg/logger"
	"github.com/open-feature/go-sdk/openfeature"
)

// NewKillSwitch returns a pointer to atomic boolean, indicating true if service is enabled and false if not.
func NewKillSwitch(ctx context.Context, client *fflag.Client, log logger.Logger) (*atomic.Bool, error) {
	enabled := new(atomic.Bool)

	flag := fflag.BooleanFlag{
		Name:         "enabled",
		DefaultValue: false,
		EvalCtx:      openfeature.NewEvaluationContext("doesn't matter here", nil),
		Options:      nil,
	}

	initValue, err := client.BooleanValue(ctx, flag.Name, flag.DefaultValue, flag.EvalCtx, flag.Options...)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch killswitch flag: %w", err)
	}

	enabled.Store(initValue)

	go client.WatchBoolFlag(ctx, flag, time.NewTicker(time.Second), func(val bool, err error) {
		if err != nil {
			logErr := log.WithRateLimit(10*time.Minute, log.Error) //nolint:mnd
			logErr("error during watch", "error", err.Error())
		} else {
			enabled.Store(val)

			log.Info("kill switch triggered", "enabled", val)
		}
	})

	return enabled, nil
}
