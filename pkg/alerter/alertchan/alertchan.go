// Package alertchan contains alert channels and related methods to retrieve channel name.
package alertchan

import (
	"errors"
	"fmt"
	"os"
)

// Channel represents an alert destination as defined via environment variables.
type Channel string

// It indicates the alerting channel.
//
// Use EnvVarName() to retrieve the environment variable name or FromEnv() to fetch the channel name,
// by passing the constant to these functions.
//
// By convention, the constant's value should begin with an underscore (e.g., _ENV_NAME).
//
//	Example
//	- Constant: const Error Channel = "_ERROR"
//	- Env     : SLACK_ERROR='account-microservice-errors'
const (
	Test      Channel = "_TEST"
	Error     Channel = "_ERROR"
	FlagError Channel = "_FLAG_ERROR"
)

// ErrChannelNotFound returned when fail to find alert channel name for a given Channel and other params.
var ErrChannelNotFound = errors.New("channel not found")

// EnvVarName returns environment variable's name.
func EnvVarName(messenger string, suffix Channel) string {
	return fmt.Sprintf("%s%s", messenger, suffix)
}

// FromEnv returns alert channel name for a given messenger and Channel suffix.
// Returns error if the channel is not found.
func FromEnv(messenger string, suffix Channel) (string, error) {
	env := EnvVarName(messenger, suffix)

	channel := os.Getenv(env)
	if channel == "" {
		return "", ErrChannelNotFound
	}

	return channel, nil
}
