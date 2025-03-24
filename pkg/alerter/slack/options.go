package slack

import (
	"github.com/Sengoku11/go-monorepo/pkg/logger"
	"github.com/Sengoku11/go-monorepo/pkg/logger/zlog"
)

// Options encapsulates Alerter options.
type Options struct {
	logger logger.Basic
	client Client
}

// DefaultOptions constructs default Options.
func DefaultOptions() Options {
	return Options{
		logger: zlog.New(),
		client: nil,
	}
}

// Logger getter.
//
//nolint:ireturn
func (o *Options) Logger() logger.Basic {
	return o.logger
}

// Client getter.
//
//nolint:ireturn
func (o *Options) Client() Client {
	return o.client
}

// Option function to modify underlying Options.
type Option func(options *Options)

// WithLogger puts the logger to Options.
func WithLogger(log logger.Basic) Option {
	return func(options *Options) {
		options.logger = log
	}
}

// WithClient puts the client to Options.
func WithClient(client Client) Option {
	return func(options *Options) {
		options.client = client
	}
}
