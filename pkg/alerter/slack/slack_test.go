package slack_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	mockslack "github.com/Sengoku11/go-monorepo/mocks/github.com/Sengoku11/go-monorepo/pkg/alerter/slack"
	mocklog "github.com/Sengoku11/go-monorepo/mocks/github.com/Sengoku11/go-monorepo/pkg/logger"
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
	channel     = "some-channel"
)

var errExpected = errors.New("hey")

func testAlerter(t *testing.T, opts ...slack.Option) (*slack.Alerter, *mockslack.MockClient) {
	t.Helper()
	mockClient := mockslack.NewMockClient(t)

	opts = append(opts, slack.WithClient(mockClient))

	return slack.New(slackToken, opts...), mockClient
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	if c := slack.NewClient(slackToken); c == nil {
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

func TestAlerter_Alert_UndefinedChannel(t *testing.T) {
	t.Parallel()

	expectedPayload := map[string]any{testKey1: testVal1}
	event := alerter.Event{
		Chan: testChannel,
		Msg:  testMessage,
		Args: expectedPayload,
	}

	alert, mockClient := testAlerter(t)

	err := alert.Alert(t.Context(), event, alerter.Options{})
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}

	if !strings.Contains(err.Error(), "alert channel is undefined") {
		t.Errorf("error must relate to undefined channel, but got: %v", err)
	}

	mockClient.AssertNumberOfCalls(t, "PostMessageContext", 0)
}

func TestAlerter_Alert_MarshalError(t *testing.T) {
	channelEnv := alertchan.EnvVarName(slack.MessengerPrefix, testChannel)
	t.Setenv(channelEnv, channel)

	expectedPayload := map[string]any{"key": func() {}}
	event := alerter.Event{
		Chan: testChannel,
		Msg:  testMessage,
		Args: expectedPayload,
	}

	alert, mockClient := testAlerter(t)

	err := alert.Alert(t.Context(), event, alerter.Options{})
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}

	if !strings.Contains(err.Error(), "fail to marshal alert payload") {
		t.Errorf("error must relate to marshalling, but got: %v", err)
	}

	mockClient.AssertNumberOfCalls(t, "PostMessageContext", 0)
}

func TestAlerter_HandleError(t *testing.T) {
	t.Parallel()

	mockedLogger := mocklog.NewMockLogger(t)
	expectedPayload := map[string]any{testKey1: testVal1}

	event := alerter.Event{
		Chan: testChannel,
		Msg:  testMessage,
		Args: expectedPayload,
	}

	alert, _ := testAlerter(t, slack.WithLogger(mockedLogger))

	mockedLogger.
		EXPECT().
		Error(errExpected.Error(), "message", testMessage, expectedPayload).
		Once()

	alert.HandleError(errExpected, event)
}

func TestAlerter_WithRateLimit(t *testing.T) {
	channelEnv := alertchan.EnvVarName(slack.MessengerPrefix, testChannel)
	t.Setenv(channelEnv, channel)

	expectedPayload := map[string]any{testKey1: testVal1}

	options := alerter.DefaultOptions()
	alerter.WithRateLimit(time.Minute)(options)

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

	if err := alert.Alert(t.Context(), event, *options); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Shouldn't call PostMessageContext method.
	if err := alert.Alert(t.Context(), event, *options); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAlerter_WithRateLimit_Resend(t *testing.T) {
	channelEnv := alertchan.EnvVarName(slack.MessengerPrefix, testChannel)
	t.Setenv(channelEnv, channel)

	expectedPayload := map[string]any{testKey1: testVal1}
	event := alerter.Event{
		Chan: testChannel,
		Msg:  testMessage,
		Args: expectedPayload,
	}

	options := alerter.DefaultOptions()
	alerter.WithRateLimit(time.Minute)(options)

	alert, mockClient := testAlerter(t)

	mockClient.
		On("PostMessageContext", mock.Anything, channel, mock.Anything).
		Return("", "", errExpected)

	// When fail to alert, rate limit shouldn't count
	if err := alert.Alert(t.Context(), event, *options); err == nil {
		t.Errorf("expected errorr, but received nil")
	}

	// This should try to alert again
	_ = alert.Alert(t.Context(), event, *options)

	mockClient.AssertNumberOfCalls(t, "PostMessageContext", 2)
}

func TestToMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		payload map[string]any
		want    string
		wantErr bool
	}{
		{
			name:    "nil payload",
			payload: nil,
			want:    "null",
			wantErr: false,
		},
		{
			name:    "empty map",
			payload: map[string]any{},
			want:    "{}",
			wantErr: false,
		},
		{
			name:    "valid payload",
			payload: map[string]any{"key": "value"},
			want:    `{"key":"value"}`,
			wantErr: false,
		},
		{
			name: "invalid payload",
			// functions are not JSON serializable.
			payload: map[string]any{"key": func() {}},
			want:    "",
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got, err := slack.ToMessage(testCase.payload)
			if (err != nil) != testCase.wantErr {
				t.Errorf("ToMessage() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}

			if got != testCase.want {
				t.Errorf("ToMessage() got = %v, want %v", got, testCase.want)
			}
		})
	}
}
