// Package channels provides type safety for selecting channel names on alert.
package channels

import (
	"errors"
	"fmt"
	"os"
)

// Prefixes of environment variables that contains channels names.
const (
	TEST      = "TEST"
	Error     = "ERROR"
	FlagError = "FLAG_ERROR"
)

// ErrChannelNotFound returned when fail to find a channel for a given channel prefix and other params.
var ErrChannelNotFound = errors.New("channel not found")

// FromENV returns channel name for a given messenger and channel prefix, e.g. "SLACK_TEST".
// Returns error if the channel is not found.
func FromENV(messenger, channelPrefix string) (string, error) {
	env := fmt.Sprintf("%s_%s", messenger, channelPrefix)

	channel := os.Getenv(env)
	if channel == "" {
		return "", ErrChannelNotFound
	}

	return channel, nil
}
