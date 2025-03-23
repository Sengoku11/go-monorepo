package alertchan_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Sengoku11/go-monorepo/pkg/alerter/alertchan"
)

const (
	testChannel     = alertchan.Channel("_TEST")
	testMessenger   = "WASSUP"
	testChannelName = "test-channel-for-tests"
)

func TestEnvVarName(t *testing.T) {
	t.Parallel()

	expectedName := fmt.Sprintf("%s%s", testMessenger, testChannel)
	name := alertchan.EnvVarName(testMessenger, testChannel)

	if name != expectedName {
		t.Errorf("expected %s, but got %s", expectedName, name)
	}
}

func TestFromEnv(t *testing.T) {
	envName := alertchan.EnvVarName(testMessenger, testChannel)

	t.Setenv(envName, testChannelName)

	channelName, err := alertchan.FromEnv(testMessenger, testChannel)
	if err != nil {
		t.Errorf("expected err to be nil, but got: %v", err)
	}

	if channelName != testChannelName {
		t.Errorf("expected channel name to be %s, but got %s", testChannelName, channelName)
	}
}

func TestFromEnv_ErrChannelNotFound(t *testing.T) {
	t.Parallel()

	_, err := alertchan.FromEnv(testMessenger, testChannel)
	if !errors.Is(err, alertchan.ErrChannelNotFound) {
		t.Errorf("expected err to be %v, but got: %v", alertchan.ErrChannelNotFound, nil)
	}
}
