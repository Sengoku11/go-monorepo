package alerter_test

import (
	"testing"
	"time"

	"github.com/Sengoku11/go-monorepo/pkg/alerter"
)

func TestDefaultOptions(t *testing.T) {
	t.Parallel()

	opts := alerter.DefaultOptions()
	if opts == nil {
		t.Error("expected options, but got nil")

		return
	}

	if timeout := opts.Timeout(); timeout != alerter.AlertTimeout {
		t.Errorf("expected timeout to be %v, but got %v", alerter.AlertTimeout, timeout)
	}
}

func TestWithTimeout(t *testing.T) {
	t.Parallel()

	expectedTimeout := time.Minute

	opts := alerter.DefaultOptions()
	if opts == nil {
		t.Error("expected options, but got nil")

		return
	}

	alerter.WithTimeout(expectedTimeout)(opts)

	if timeout := opts.Timeout(); timeout != expectedTimeout {
		t.Errorf("expected timeout to be %v, but got %v", expectedTimeout, timeout)
	}
}
