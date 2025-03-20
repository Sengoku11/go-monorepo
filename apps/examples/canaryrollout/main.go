// Example of how you can use fflag.Client to implement a canary rollout.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Sengoku11/go-monorepo/apps/example/canaryrollout/internal/featureflag"
	"github.com/Sengoku11/go-monorepo/pkg/bootstrap"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/open-feature/go-sdk/openfeature"
)

func main() {
	ctx, cancel, log := bootstrap.Default()
	defer cancel(nil)

	// You can configure environments in ./internal/featureflag dir.
	environment := os.Getenv("ENVIRONMENT")
	if environment != "prod" && environment != "uat" {
		environment = "prod"
	}

	// Initialize a feature flag client.
	client, err := featureflag.NewClient(ctx, log, environment)
	if err != nil {
		log.Panic(err.Error())
	}

	// To test, update ./internal/featureflag/goff-flags.yaml, and try these commands:
	//	curl -H "X-USER-ID: admin" http://localhost:3031
	//	curl http://localhost:3031

	mux := chi.NewMux()
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		var resp string

		// Pass a static userID for deterministic results, like "X-USER-ID: any-user".
		userID := r.Header.Get("X-User-Id")
		if userID == "" {
			// Simulating a user.
			id, _ := uuid.NewV4()
			userID = id.String()
		}

		// Evaluate the "new_feature" flag.
		evalCtx := openfeature.NewEvaluationContext(userID, nil)

		featEnabled, err := client.BooleanValue(ctx, "new_feature", false, evalCtx)
		if err != nil {
			http.Error(w, fmt.Sprintf("unexpected error: %v", err), http.StatusInternalServerError)

			return
		}

		if featEnabled {
			resp = tryNewFeature(userID)
		} else {
			resp = defaultFunc(userID)
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

func tryNewFeature(userID string) string {
	return fmt.Sprintf("hey %s! wanna see our new feature?", userID)
}

func defaultFunc(userID string) string {
	return userID + ", we have nothing new for you"
}
