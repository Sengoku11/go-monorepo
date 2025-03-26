package fflag

import "github.com/open-feature/go-sdk/openfeature"

// MockClient defines an interface that embeds OpenFeature client.
// This interface exists to facilitate automatic mock generation and simplify testing.
type MockClient interface {
	openfeature.IClient
}

// MockProvider defines an interface that embeds OpenFeature provider.
// This interface exists to facilitate automatic mock generation and simplify testing.
type MockProvider interface {
	openfeature.FeatureProvider
}

// WithMockClient instantiates with a mock client.
func WithMockClient(client MockClient) *Client {
	return &Client{client}
}
