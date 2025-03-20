// Example of how you can use fflag.Client to implement a kill switch.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Sengoku11/go-monorepo/apps/example/killswitch/internal/featureflag"
	"github.com/Sengoku11/go-monorepo/pkg/bootstrap"
	"github.com/go-chi/chi/v5"
	"github.com/open-feature/go-sdk/openfeature"
)

func main() {
	ctx, cancel, log := bootstrap.Default()
	defer cancel(nil)

	// Initialize a feature flag client
	client, err := featureflag.NewClient(ctx, log)
	if err != nil {
		log.Panic(err.Error())
	}

	// Kill switch implementation using feature flags
	enabled, err := featureflag.NewKillSwitch(ctx, client, log)
	if err != nil {
		log.Panic(err.Error())
	}

	// Now testing with a showcase API
	// Open http://localhost:3031 in your browser and try to change goff-flags.yaml configs in internal dir.
	mux := chi.NewMux()
	mux.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		var resp string
		if !enabled.Load() {
			resp = "service is disabled"
		} else {
			resp = "enabled"
		}

		_ = json.NewEncoder(w).Encode(resp)
	})

	srv := http.Server{Addr: ":3031", Handler: mux} //nolint
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(fmt.Sprintf("listen: %s\n", err))
		}
	}()

	<-ctx.Done()
	_ = srv.Shutdown(ctx)
	openfeature.Shutdown()
}
