// Example application that you can use as a playground
package main

import (
	"sync"

	"github.com/Sengoku11/go-monorepo/apps/example/internal/alertexample"
	"github.com/Sengoku11/go-monorepo/apps/example/internal/logexample"
	"github.com/Sengoku11/go-monorepo/pkg/bootstrap"
)

func main() {
	ctx, cancel, log := bootstrap.Default()
	defer cancel(nil)

	var wg sync.WaitGroup

	wg.Add(1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		logexample.Run(log)
		alertexample.Run(ctx, log)
	}()
	go func() {
		defer wg.Done()
		log.Info("emulating an ongoing background process; press ctrl+c to shutdown gracefully")
		<-ctx.Done()
	}()

	wg.Wait()
	log.Info("terminating app")
}
