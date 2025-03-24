package slack_test

import (
	"testing"

	"github.com/Sengoku11/go-monorepo/pkg/alerter/slack"
	"github.com/Sengoku11/go-monorepo/pkg/logger/zlog"
)

func TestDefaultOptions(t *testing.T) {
	t.Parallel()

	options := slack.DefaultOptions()
	if options.Logger() == nil {
		t.Errorf("expected non nil logger")
	}

	if options.Client() != nil {
		t.Errorf("expected nil client")
	}
}

func TestWithLogger(t *testing.T) {
	t.Parallel()

	log := zlog.New()
	options := slack.DefaultOptions()

	slack.WithLogger(log)(&options)

	if log != options.Logger() {
		t.Errorf("expected the same underlying logger")
	}
}

func TestWithClient(t *testing.T) {
	t.Parallel()

	client := slack.NewClient(slackToken)
	options := slack.DefaultOptions()

	slack.WithClient(client)(&options)

	if client != options.Client() {
		t.Errorf("expected the same underlying client")
	}

	client = slack.NewClient(slackToken)
	if client == options.Client() {
		t.Errorf("clients should be different now")
	}
}
