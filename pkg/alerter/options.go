package alerter

import "time"

// AlertTimeout defines how long to wait for alert to be sent.
const AlertTimeout = time.Second * 2

// Options for Alerter.
type Options struct {
	timeout time.Duration
}

// Timeout returns timeout value stored in Options.
func (o *Options) Timeout() time.Duration {
	return o.timeout
}

// Option defines a function type to set options for Alert.
type Option func(opts *Options)

// DefaultOptions constructs Options with default values.
func DefaultOptions() *Options {
	return &Options{
		timeout: AlertTimeout,
	}
}

// WithTimeout sets a timeout for sending alert.
func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.timeout = timeout
	}
}
