// Package alertchan contains alert channels and related methods to retrieve channel name.
package alertchan

import (
	"errors"
	"fmt"
	"os"
)

// Target represents an alert destination as defined via environment variables.
type Target string

// It indicates the alerting channel.
//
// Use EnvVarName() to retrieve the environment variable name or FromEnv() to fetch the channel name,
// by passing the constant to these functions.
//
// By convention, the constant's value should begin with an underscore (e.g., _ENV_NAME).
//
//	Example
//	- Constant: const Error Target = "_ERROR"
//	- Env     : SLACK_ERROR='account-microservice-errors'
const (
	Test      Target = "_TEST"
	Error     Target = "_ERROR"
	FlagError Target = "_FLAG_ERROR"
)

// ErrChannelNotFound returned when fail to find alert channel name for a given Target and other params.
var ErrChannelNotFound = errors.New("channel not found")

// EnvVarName returns environment variable's name.
func EnvVarName(messenger string, target Target) string {
	return fmt.Sprintf("%s%s", messenger, target)
}

// FromEnv returns alert channel name for a given messenger and Target suffix.
// Returns error if the channel is not found.
func FromEnv(messenger string, target Target) (string, error) {
	env := EnvVarName(messenger, target)

	channel := os.Getenv(env)
	if channel == "" {
		return "", ErrChannelNotFound
	}

	return channel, nil
}
