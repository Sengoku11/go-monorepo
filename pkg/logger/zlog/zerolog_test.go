//nolint:testpackage
package zlog

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	mockalerter "github.com/Sengoku11/go-monorepo/mocks/github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/Sengoku11/go-monorepo/pkg/alerter"
	"github.com/Sengoku11/go-monorepo/pkg/logger"
	"github.com/stretchr/testify/mock"
)

const (
	testMessage                  = "Hello World"
	testKey1, testKey2, testKey3 = "hello_1", "hello_2", "hello_3"
	testVal1, testVal2, testVal3 = "world_1", "world_2", "world_3"
)

// newTestLogger creates a Logger that writes to a bytes.Buffer.
func newTestLogger() (*Logger, *bytes.Buffer) {
	var buf bytes.Buffer

	l := New()
	*l.logger = l.logger.Output(&buf)

	return l, &buf
}

func TestAlert(t *testing.T) {
	t.Parallel()

	mockAlert := mockalerter.NewMockAlerter(t)
	expectedPayload := map[string]any{testKey1: testVal1}

	mockAlert.
		EXPECT().
		Alert(mock.Anything, alerter.Test, testMessage, expectedPayload).
		Return(nil)

	log, _ := newTestLogger()

	// Add two "different" alerters
	log.AddHook(mockAlert)
	log.AddHook(mockAlert)
	log.Alert(t.Context(), alerter.Test, testMessage, expectedPayload)

	mockAlert.AssertNumberOfCalls(t, "Alert", 2)
}

func TestAlert_NotCalled(t *testing.T) {
	t.Parallel()

	mockAlert := mockalerter.NewMockAlerter(t)
	expectedPayload := map[string]any{testKey1: testVal1}

	log, _ := newTestLogger()

	// Alert without any hook
	log.Alert(t.Context(), alerter.Test, testMessage, expectedPayload)

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
			expectedLogMsg:   logger.ErrSendTimeout.Error(),
			expectedErrorSub: fmt.Sprintf(`"error":"%s"`, context.DeadlineExceeded),
		},
		{
			name:             "canceled",
			returnedError:    context.Canceled,
			expectedLogMsg:   logger.ErrSendAlert.Error(),
			expectedErrorSub: fmt.Sprintf(`"error":"%s"`, context.Canceled),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			expectedChannel := alerter.Test
			expectedPayload := map[string]any{testKey1: testVal1}

			log, buf := newTestLogger()
			mockAlert := mockalerter.NewMockAlerter(t)
			mockAlert.
				On("Alert", mock.Anything, expectedChannel, testMessage, expectedPayload).
				Return(testCase.returnedError).
				Once()

			log.AddHook(mockAlert)
			log.Alert(t.Context(), expectedChannel, testMessage, expectedPayload)

			out := buf.String()

			expectedSubstrings := []string{
				`"level":"error"`,
				testCase.expectedErrorSub,
				fmt.Sprintf(`"msg":"%s"`, testMessage),
				fmt.Sprintf(`"message":"%s"`, testCase.expectedLogMsg),
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
		logFunc func(l *Logger)
	}{
		{
			level: "info",
			logFunc: func(l *Logger) {
				l.Info(testMessage, testKey1, testVal1, testKey2, testVal2, testKey3, testVal3)
			},
		},
		{
			level: "warn",
			logFunc: func(l *Logger) {
				l.Warn(testMessage, testKey1, testVal1, testKey2, testVal2, testKey3, testVal3)
			},
		},
		{
			level: "error",
			logFunc: func(l *Logger) {
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

func TestZerologLogger_WithoutRateLim(t *testing.T) {
	t.Parallel()

	log, buf := newTestLogger()
	runs := 10

	for range runs {
		log.Info(testMessage, testKey1, testVal1, testKey2, testVal2, testKey3, testVal3)
	}

	out := buf.String()

	// Expecting the same number of logs we have logged.
	logEntries := strings.Split(strings.TrimSpace(out), "\n")
	if len(logEntries) != runs {
		t.Errorf("expected %d logs, but got %d", runs, len(logEntries))
	}
}

func TestZerologLogger_WithRateLim(t *testing.T) {
	t.Parallel()

	log, buf := newTestLogger()
	runs := 10
	expectedLogs := runs + 1

	info := log.WithRateLim(time.Minute, log.Info)
	// Should log just once.
	for range runs {
		info(testMessage, testKey1, testVal1, testKey2, testVal2, testKey3, testVal3)
	}

	// Run without rate lim as well, to check that there is no interfering.
	for range runs {
		log.Info(testMessage, testKey1, testVal1, testKey2, testVal2, testKey3, testVal3)
	}

	out := buf.String()

	// Even though we logged 2x10 times, the second loop should only allow one log.
	logEntries := strings.Split(strings.TrimSpace(out), "\n")
	if len(logEntries) != expectedLogs {
		t.Errorf("expected %d log, but got %d", expectedLogs, len(logEntries))
	}
}

func TestZerologLogger_WithCallStack(t *testing.T) {
	t.Parallel()

	log, buf := newTestLogger()

	expectedSubstrings := []string{
		`"level":"info"`,
		fmt.Sprintf(`"message":"%s"`, testMessage),
		fmt.Sprintf(`"%s":"%s"`, testKey1, testVal1),
		fmt.Sprintf(`"%s":"%s"`, testKey2, testVal2),
		fmt.Sprintf(`"%s":"%s"`, testKey3, testVal3),
	}

	info := log.WithCallStack(log.Info)
	info(testMessage, testKey1, testVal1, testKey2, testVal2, testKey3, testVal3)

	out := buf.String()
	for _, substring := range expectedSubstrings {
		if !strings.Contains(out, substring) {
			t.Errorf(`expected log output to contain %s, got %q`, substring, out)
		}
	}

	// Using a regular expression to check that the stack field contains the function name.
	// This regex looks for "stack":"<anything>TestZerologLogger_WithCallStack<anything>"
	re := regexp.MustCompile(`"stack":"[^"]*TestZerologLogger_WithCallStack[^"]*"`)
	if !re.MatchString(out) {
		t.Errorf("expected stack field to contain function name TestZerologLogger_WithCallStack, got %q", out)
	}
}
