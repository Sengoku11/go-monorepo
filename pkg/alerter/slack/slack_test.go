package slack_test

import (
	"testing"

	mockslack "github.com/Sengoku11/go-monorepo/mocks/github.com/Sengoku11/go-monorepo/pkg/alerter/slack"
	"github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/Sengoku11/go-monorepo/pkg/alerter/alertchan"
	"github.com/Sengoku11/go-monorepo/pkg/alerter/slack"
	"github.com/stretchr/testify/mock"
)

const (
	slackToken  = "anything"
	testMessage = "hello"
	testChannel = "_TESTING"
	testKey1    = "key1"
	testVal1    = "val1"
)

func testAlerter(t *testing.T) (*slack.Alerter, *mockslack.MockClient) {
	t.Helper()
	mockClient := mockslack.NewMockClient(t)

	return slack.New(slackToken, slack.WithClient(mockClient)), mockClient
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	client := slack.NewClient(slackToken)
	if client == nil {
		t.Errorf("expected non nil client")
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	alert := slack.New(slackToken)
	if alert == nil {
		t.Errorf("expected non nil alerter")
	}
}

func TestAlerter_Alert(t *testing.T) {
	channelEnv := alertchan.EnvVarName(slack.MessengerPrefix, testChannel)
	channel := "some-channel"
	t.Setenv(channelEnv, channel)

	expectedPayload := map[string]any{testKey1: testVal1}

	event := alerter.Event{
		Chan: testChannel,
		Msg:  testMessage,
		Args: expectedPayload,
	}

	alert, mockClient := testAlerter(t)

	mockClient.
		EXPECT().
		PostMessageContext(mock.Anything, channel, mock.Anything).
		Return("", "", nil).
		Once()

	if err := alert.Alert(t.Context(), event, alerter.Options{}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
