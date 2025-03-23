package alerter_test

import (
	"context"
	"testing"

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

	expectedEvent := alerter.Event{
		Chan: alertchan.Test,
		Msg:  testMessage,
		Args: map[string]any{testKey1: testVal1},
	}

	mockAlert.
		EXPECT().
		Alert(mock.Anything, expectedEvent).
		Return(nil)

	// Add two "different" alerters
	alert := newTestLogger(mockAlert, mockAlert)

	alert.Alert(t.Context(), expectedEvent)

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

	testCases := []struct {
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

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mockAlert := mockalerter.NewMockAlerter(t)
			mockAlert.
				On("Alert", mock.Anything, expectedEvent).
				Return(testCase.returnedError).
				Once()

			mockAlert.
				EXPECT().
				HandleError(testCase.returnedError, expectedEvent).
				Once()

			alert := newTestLogger(mockAlert)
			alert.Alert(t.Context(), expectedEvent)

			mockAlert.AssertNumberOfCalls(t, "Alert", 1)
		})
	}
}
