// Package alertchan contains alert channels and related methods to retrieve channel name.
package alertchan

import (
	"errors"
	"fmt"
	"os"
)

// Channel represents a constant for channel's suffix.
type Channel string

// Channel of environment variables that contains channels names.
//
// Real channel name is different for each messenger, so constants are used to make alerting more generic.
const (
	Test      Channel = "TEST"
	Error     Channel = "ERROR"
	FlagError Channel = "FLAG_ERROR"
)

// ErrChannelNotFound returned when fail to find a channel for a given channel prefix and other params.
var ErrChannelNotFound = errors.New("channel not found")

// FromENV returns real channel name for a given messenger and channel prefix, e.g. "SLACK_TEST".
// Returns error if the channel is not found.
func FromENV(messenger string, suffix Channel) (string, error) {
	env := fmt.Sprintf("%s_%s", messenger, suffix)

	channel := os.Getenv(env)
	if channel == "" {
		return "", ErrChannelNotFound
	}

	return channel, nil
}
