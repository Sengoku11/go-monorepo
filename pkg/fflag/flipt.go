package fflag

import (
	"fmt"

	"github.com/Sengoku11/go-monorepo/pkg/logger"
	flipt "github.com/open-feature/go-sdk-contrib/providers/flipt/pkg/provider"
	"github.com/open-feature/go-sdk/openfeature"
)

// NewFlipt returns wrapped openfeature.Client with flipt.io provider.
func NewFlipt(url string, log logger.Logger, namespace string) (*Client, error) {
	provider := flipt.NewProvider(
		flipt.WithAddress(url),
		flipt.ForNamespace(namespace),
	)

	if err := openfeature.SetProviderAndWait(provider); err != nil {
		return nil, fmt.Errorf("failed to initialize provider: %w", err)
	}

	client := openfeature.NewClient(namespace)

	return &Client{client, log}, nil
}
