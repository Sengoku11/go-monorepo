// Package featureflag encapsulates feature flag's client initialization.
package featureflag

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Sengoku11/go-monorepo/pkg/fflag"
	"github.com/Sengoku11/go-monorepo/pkg/logger"
	gofeatureflaginprocess "github.com/open-feature/go-sdk-contrib/providers/go-feature-flag-in-process/pkg"
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag/retriever/fileretriever"
)

var errPathNotFound = errors.New("runtime.Caller(0) failed")

// NewClient returns an instantiated fflag.Client using GoFeatureFlag in-process provider.
//
// The available flags storages are:
//   - GitHub
//   - GitLab
//   - HTTP endpoint
//   - AWS S3
//   - Local file
//   - Google Cloud Storage
//   - Kubernetes ConfigMaps
//   - MongoDB
//   - Redis
//   - BitBucket
//   - AzBlobStorage
//   - Flipt, Unleash, and other cloud providers.
func NewClient(ctx context.Context, log logger.Logger, environment string) (*fflag.Client, error) {
	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		return nil, errPathNotFound
	}

	baseDir := filepath.Dir(currentFilePath)
	absPath := filepath.Join(baseDir, "goff-flags.yaml")

	//nolint:mnd,exhaustruct
	options := gofeatureflaginprocess.ProviderOptions{
		GOFeatureFlagConfig: &ffclient.Config{
			PollingInterval: 100 * time.Millisecond,
			Context:         ctx,
			Environment:     environment,
			Retriever: &fileretriever.Retriever{
				Path: absPath,
			},
		},
	}

	provider, err := gofeatureflaginprocess.NewProviderWithContext(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to start provider: %w", err)
	}

	client, err := fflag.New("example", provider, log)
	if err != nil {
		return nil, fmt.Errorf("failed to get fflag client: %w", err)
	}

	return client, nil
}
