//nolint:testpackage
package logger

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	mockalerter "github.com/Sengoku11/go-monorepo/mocks/github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/stretchr/testify/mock"
)

const (
	testMessage                  = "Hello World"
	testKey1, testKey2, testKey3 = "hello_1", "hello_2", "hello_3"
	testVal1, testVal2, testVal3 = "world_1", "world_2", "world_3"
)

// newTestLogger creates a ZerologLogger that writes to a bytes.Buffer.
func newTestLogger() (*ZerologLogger, *bytes.Buffer) {
	var buf bytes.Buffer

	l := NewZerologLogger()
	*l.logger = l.logger.Output(&buf)

	return l, &buf
}

func TestAlert(t *testing.T) {
	t.Parallel()

	mockAlert := mockalerter.NewMockAlerter(t)
	expectedOptions := alerter.Options{ChannelSuffix: "TEST", RateLimit: time.Second}
	expectedPayload := map[string]any{testKey1: testVal1}
	expectedEvent := alerter.Event{
		Message: testMessage,
		Payload: expectedPayload,
		Options: expectedOptions,
	}

	mockAlert.
		EXPECT().
		Alert(mock.Anything, expectedEvent).
		Return(nil)

	log, _ := newTestLogger()

	// Add two "different" alerters
	log.AddHook(mockAlert)
	log.AddHook(mockAlert)
	log.Alert(t.Context(), testMessage, expectedOptions, expectedPayload)

	mockAlert.AssertNumberOfCalls(t, "Alert", 2)
}

func TestAlert_NotCalled(t *testing.T) {
	t.Parallel()

	mockAlert := mockalerter.NewMockAlerter(t)
	expectedOptions := alerter.Options{ChannelSuffix: "TEST", RateLimit: time.Second}
	expectedPayload := map[string]any{testKey1: testVal1}

	log, _ := newTestLogger()

	// Alert without any hook
	log.Alert(t.Context(), testMessage, expectedOptions, expectedPayload)

	mockAlert.AssertNumberOfCalls(t, "Alert", 0)
}

func TestAlert_Errors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		returnedError    error
		expectedLogMsg   string
		expectedErrorSub string
	}{
		{
			name:             "timeout",
			returnedError:    context.DeadlineExceeded,
			expectedLogMsg:   ErrSendTimeout.Error(),
			expectedErrorSub: fmt.Sprintf(`"error":"%s"`, context.DeadlineExceeded),
		},
		{
			name:             "canceled",
			returnedError:    context.Canceled,
			expectedLogMsg:   ErrSendAlert.Error(),
			expectedErrorSub: fmt.Sprintf(`"error":"%s"`, context.Canceled),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			expectedChannel := "TEST_ALERT"
			expectedOptions := alerter.Options{ChannelSuffix: expectedChannel, RateLimit: time.Second}
			expectedPayload := map[string]any{testKey1: testVal1}

			log, buf := newTestLogger()
			mockAlert := mockalerter.NewMockAlerter(t)
			log.AddHook(mockAlert)

			mockAlert.
				On("Alert", mock.Anything, mock.Anything).
				Return(testCase.returnedError).
				Once()

			log.Alert(t.Context(), testMessage, expectedOptions, expectedPayload)

			out := buf.String()

			expectedSubstrings := []string{
				`"level":"error"`,
				fmt.Sprintf(`"msg":"%s"`, testMessage),
				fmt.Sprintf(`"message":"%s"`, testCase.expectedLogMsg),
				testCase.expectedErrorSub,
				fmt.Sprintf(`"channel":"%s"`, expectedChannel),
			}

			for _, substr := range expectedSubstrings {
				if !strings.Contains(out, substr) {
					t.Errorf("expected log output to contain %s, got %q", substr, out)
				}
			}

			mockAlert.AssertNumberOfCalls(t, "Alert", 1)
		})
	}
}

func TestZerologLogger(t *testing.T) {
	t.Parallel()

	expectedSubstrings := []string{
		fmt.Sprintf(`"message":"%s"`, testMessage),
		fmt.Sprintf(`"%s":"%s"`, testKey1, testVal1),
		fmt.Sprintf(`"%s":"%s"`, testKey2, testVal2),
		fmt.Sprintf(`"%s":"%s"`, testKey3, testVal3),
	}

	testCases := []struct {
		level   string
		logFunc func(l *ZerologLogger)
	}{
		{
			level: "info",
			logFunc: func(l *ZerologLogger) {
				l.Info(testMessage, testKey1, testVal1, testKey2, testVal2, testKey3, testVal3)
			},
		},
		{
			level: "warn",
			logFunc: func(l *ZerologLogger) {
				l.Warn(testMessage, testKey1, testVal1, testKey2, testVal2, testKey3, testVal3)
			},
		},
		{
			level: "error",
			logFunc: func(l *ZerologLogger) {
				l.Error(testMessage, testKey1, testVal1, testKey2, testVal2, testKey3, testVal3)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.level, func(t *testing.T) {
			t.Parallel()

			log, buf := newTestLogger()
			testCase.logFunc(log)

			out := buf.String()
			expectedLevel := fmt.Sprintf(`"level":"%s"`, testCase.level)

			if !strings.Contains(out, expectedLevel) {
				t.Errorf("expected log to contain level %q, got %q", testCase.level, out)
			}

			for _, substr := range expectedSubstrings {
				if !strings.Contains(out, substr) {
					t.Errorf("expected log output to contain %s, got %q", substr, out)
				}
			}
		})
	}
}

func TestZerologLogger_Panic(t *testing.T) {
	t.Parallel()

	log, buf := newTestLogger()

	expectedSubstrings := []string{
		`"level":"panic"`,
		fmt.Sprintf(`"message":"%s"`, testMessage),
		fmt.Sprintf(`"%s":"%s"`, testKey1, testVal1),
		fmt.Sprintf(`"%s":"%s"`, testKey2, testVal2),
		fmt.Sprintf(`"%s":"%s"`, testKey3, testVal3),
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`did not panic`)
		}

		out := buf.String()
		for _, substring := range expectedSubstrings {
			if !strings.Contains(out, substring) {
				t.Errorf(`expected log output to contain %s, got %q`, substring, out)
			}
		}
	}()

	log.Panic(testMessage, testKey1, testVal1, testKey2, testVal2, testKey3, testVal3)
}

func TestZerologLogger_Debug(t *testing.T) {
	t.Parallel()

	log, buf := newTestLogger()

	log.Debug(testMessage, testKey1, testVal1, testKey2, testVal2, testKey3, testVal3)

	// By default, debug mode is turned off
	if out := buf.String(); len(out) > 0 {
		t.Errorf("Expected empty buffer but got %q", out)
	}
}

func TestZerologLogger_EnableDebugMode(t *testing.T) {
	t.Parallel()

	log, buf := newTestLogger()

	expectedSubstrings := []string{
		`"level":"debug"`,
		fmt.Sprintf(`"message":"%s"`, testMessage),
		fmt.Sprintf(`"%s":"%s"`, testKey1, testVal1),
		fmt.Sprintf(`"%s":"%s"`, testKey2, testVal2),
		fmt.Sprintf(`"%s":"%s"`, testKey3, testVal3),
	}

	log.EnableDebugMode()
	log.Debug(testMessage, testKey1, testVal1, testKey2, testVal2, testKey3, testVal3)

	out := buf.String()
	for _, substring := range expectedSubstrings {
		if !strings.Contains(out, substring) {
			t.Errorf(`expected log output to contain %s, got %q`, substring, out)
		}
	}
}

func TestZerologLogger_DisableDebugMode(t *testing.T) {
	t.Parallel()

	log, buf := newTestLogger()

	log.EnableDebugMode()
	log.DisableDebugMode()
	log.Debug(testMessage, testKey1, testVal1, testKey2, testVal2, testKey3, testVal3)

	if out := buf.String(); len(out) > 0 {
		t.Errorf("Expected empty buffer but got %q", out)
	}
}

//nolint:paralleltest,tparallel
func TestZerologLogger_WithRateLim(t *testing.T) {
	t.Parallel()

	log, buf := newTestLogger()

	expectedSubstrings := []string{
		`"level":"info"`,
		fmt.Sprintf(`"message":"%s"`, testMessage),
		fmt.Sprintf(`"%s":"%s"`, testKey1, testVal1),
		fmt.Sprintf(`"%s":"%s"`, testKey2, testVal2),
		fmt.Sprintf(`"%s":"%s"`, testKey3, testVal3),
	}

	runs := 10
	testCases := []struct {
		name         string
		infoFunc     LogMethod
		expectedLogs int
	}{
		{
			name:         "without rate limit",
			infoFunc:     log.Info,
			expectedLogs: runs,
		},
		{
			name:         "with rate limit",
			infoFunc:     log.WithRateLim(time.Minute, log.Info),
			expectedLogs: 1,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			buf.Reset()

			for range runs {
				testCase.infoFunc(testMessage, testKey1, testVal1, testKey2, testVal2, testKey3, testVal3)
			}

			out := buf.String()
			for _, substring := range expectedSubstrings {
				if !strings.Contains(out, substring) {
					t.Errorf(`expected log output to contain %s, got %q`, substring, out)
				}
			}

			logEntries := strings.Split(strings.TrimSpace(out), "\n")

			logCount := len(logEntries)
			if logCount != testCase.expectedLogs {
				t.Errorf("expected %d logs, but printed: %d\n", testCase.expectedLogs, logCount)
			}
		})
	}
}
