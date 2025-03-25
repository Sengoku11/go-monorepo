package alerter_test

import (
	"context"
	"testing"
	"time"

	mockalerter "github.com/Sengoku11/go-monorepo/mocks/github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/Sengoku11/go-monorepo/pkg/alerter/alertchan"
	"github.com/stretchr/testify/mock"
)

const (
	testMessage = "Hello World"
	testKey1    = "hello_1"
	testVal1    = "world_1"
)

func newTestLogger(hooks ...alerter.Alerter) *alerter.Client {
	return alerter.New(hooks...)
}

func TestAlert(t *testing.T) {
	t.Parallel()

	mockAlert := mockalerter.NewMockAlerter(t)
	options := alerter.DefaultOptions()

	alerter.WithTimeout(7 * time.Second)(options)

	expectedEvent := alerter.Event{
		Chan: alertchan.Test,
		Msg:  testMessage,
		Args: map[string]any{testKey1: testVal1},
	}

	mockAlert.
		EXPECT().
		Alert(mock.Anything, expectedEvent, *options).
		Return(nil)

	// Construct with two "different" alerters
	alert := newTestLogger(mockAlert, mockAlert)

	alert.Alert(t.Context(), expectedEvent, alerter.WithTimeout(7*time.Second))

	mockAlert.AssertNumberOfCalls(t, "Alert", 2)
}

func TestAlert_NotCalled(t *testing.T) {
	t.Parallel()

	mockAlert := mockalerter.NewMockAlerter(t)
	alert := newTestLogger()

	expectedEvent := alerter.Event{
		Chan: alertchan.Test,
		Msg:  testMessage,
		Args: map[string]any{testKey1: testVal1},
	}

	// Alert without any hook
	alert.Alert(t.Context(), expectedEvent)

	mockAlert.AssertNumberOfCalls(t, "Alert", 0)
}

func TestAlert_Errors(t *testing.T) {
	t.Parallel()

	expectedEvent := alerter.Event{
		Chan: alertchan.Test,
		Msg:  testMessage,
		Args: map[string]any{testKey1: testVal1},
	}
	expectedOptions := alerter.DefaultOptions()

	tests := []struct {
		name          string
		returnedError error
	}{
		{
			name:          "timeout",
			returnedError: context.DeadlineExceeded,
		},
		{
			name:          "canceled",
			returnedError: context.Canceled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockAlert := mockalerter.NewMockAlerter(t)
			mockAlert.
				On("Alert", mock.Anything, expectedEvent, *expectedOptions).
				Return(tt.returnedError).
				Once()

			mockAlert.
				EXPECT().
				HandleError(tt.returnedError, expectedEvent).
				Once()

			alert := newTestLogger(mockAlert)
			alert.Alert(t.Context(), expectedEvent)

			mockAlert.AssertNumberOfCalls(t, "Alert", 1)
		})
	}
}
