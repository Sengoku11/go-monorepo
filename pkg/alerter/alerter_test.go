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
	testMessage                  = "Hello World"
	testKey1, testKey2, testKey3 = "hello_1", "hello_2", "hello_3"
	testVal1, testVal2, testVal3 = "world_1", "world_2", "world_3"
)

func newTestLogger() *alerter.Client {
	return alerter.New()
}

func TestAlert(t *testing.T) {
	t.Parallel()

	mockAlert := mockalerter.NewMockAlerter(t)
	expectedPayload := map[string]any{testKey1: testVal1}

	mockAlert.
		EXPECT().
		Alert(mock.Anything, alertchan.Test, testMessage, expectedPayload).
		Return(nil)

	alert := newTestLogger()

	// Add two "different" alerters
	alert.AddHook(mockAlert)
	alert.AddHook(mockAlert)
	alert.Alert(t.Context(), alertchan.Test, testMessage, expectedPayload)

	mockAlert.AssertNumberOfCalls(t, "Alert", 2)
}

func TestAlert_NotCalled(t *testing.T) {
	t.Parallel()

	mockAlert := mockalerter.NewMockAlerter(t)
	expectedPayload := map[string]any{testKey1: testVal1}

	alert := newTestLogger()

	// Alert without any hook
	alert.Alert(t.Context(), alertchan.Test, testMessage, expectedPayload)

	mockAlert.AssertNumberOfCalls(t, "Alert", 0)
}

func TestAlert_Errors(t *testing.T) {
	t.Parallel()

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

			expectedChannel := alertchan.Test
			expectedPayload := map[string]any{testKey1: testVal1}

			alert := newTestLogger()
			mockAlert := mockalerter.NewMockAlerter(t)
			mockAlert.
				On("Alert", mock.Anything, expectedChannel, testMessage, expectedPayload).
				Return(testCase.returnedError).
				Once()

			mockAlert.
				EXPECT().
				HandleError(testCase.returnedError, testMessage, expectedPayload).
				Once()

			alert.AddHook(mockAlert)
			alert.Alert(t.Context(), expectedChannel, testMessage, expectedPayload)

			mockAlert.AssertNumberOfCalls(t, "Alert", 1)
		})
	}
}
