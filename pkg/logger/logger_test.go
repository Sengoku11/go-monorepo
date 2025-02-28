package logger_test

import (
	"testing"
	"time"

	mockalerter "github.com/Sengoku11/go-monorepo/mocks/github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/Sengoku11/go-monorepo/pkg/logger"
	"github.com/stretchr/testify/mock"
)

func TestAlert(t *testing.T) {
	t.Parallel()

	mockAlert := mockalerter.NewMockAlerter(t)
	expectedMessage := "Hello World"
	expectedOptions := alerter.Options{ChannelSuffix: "TEST", RateLimit: time.Second}
	expectedPayload := map[string]any{"testKey": "testValue"}
	expectedEvent := alerter.Event{
		Message: expectedMessage,
		Payload: expectedPayload,
		Options: expectedOptions,
	}

	mockAlert.
		EXPECT().
		Alert(mock.Anything, expectedEvent).
		Return(nil)

	log := logger.NewZerologLogger()

	log.AddHook(mockAlert)
	log.Alert(t.Context(), expectedMessage, expectedOptions, expectedPayload)
}
