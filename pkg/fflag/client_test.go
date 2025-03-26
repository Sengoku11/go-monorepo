package fflag_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	mockfflag "github.com/Sengoku11/go-monorepo/mocks/github.com/Sengoku11/go-monorepo/pkg/fflag"
	"github.com/Sengoku11/go-monorepo/pkg/fflag"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/stretchr/testify/mock"
)

const (
	namespace = "testing"
	testFlag  = "test_flag"
	testUser  = "test_user"
)

var errTest = errors.New("hey toster")

func newMock(t *testing.T) (*fflag.Client, *mockfflag.MockMockClient) {
	t.Helper()

	client := mockfflag.NewMockMockClient(t)

	return fflag.WithMockClient(client), client
}

func TestNew(t *testing.T) {
	t.Parallel()

	provider := mockfflag.NewMockMockProvider(t)
	provider.
		On("Metadata").
		Return(openfeature.Metadata{Name: namespace})

	client, err := fflag.New(namespace, provider)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if client == nil {
		t.Errorf("expected client to be non nil")
	}
}

//nolint:funlen
func TestClient_WatchBoolFlag(t *testing.T) {
	t.Parallel()

	parentCtx := t.Context()

	tests := []struct {
		name       string
		defaultVal bool
		firstVal   bool
		secondVal  bool
		targetVal  bool
		calledBack bool
		firstErr   error
		secondErr  error
	}{
		{
			name:       "successfully update",
			defaultVal: false,
			firstVal:   true,
			secondVal:  false,
			targetVal:  false,
			calledBack: true,
			firstErr:   nil,
			secondErr:  nil,
		},
		{
			name:       "not called back",
			defaultVal: true,
			firstVal:   true,
			secondVal:  true,
			targetVal:  true,
			calledBack: false,
			firstErr:   nil,
			secondErr:  nil,
		},
		{
			name:       "callback error",
			defaultVal: true,
			firstVal:   true,
			secondVal:  true,
			targetVal:  true,
			calledBack: true,
			firstErr:   errTest,
			secondErr:  nil,
		},
		{
			name:       "second update error",
			defaultVal: false,
			firstVal:   true,
			secondVal:  false,
			targetVal:  true,
			calledBack: true,
			firstErr:   nil,
			secondErr:  errTest,
		},
		{
			name:       "did not update flag because of errors",
			defaultVal: false,
			firstVal:   true,
			secondVal:  true,
			targetVal:  false,
			calledBack: true,
			firstErr:   errTest,
			secondErr:  errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(parentCtx)
			defer cancel()

			client, mockedClient := newMock(t)
			flag := fflag.BooleanFlag{
				Name:         testFlag,
				DefaultValue: tt.defaultVal,
				EvalCtx:      openfeature.NewEvaluationContext(testUser, nil),
				Options:      nil,
			}

			// First call to get value
			mockedClient.
				On("BooleanValue", ctx, flag.Name, flag.DefaultValue, flag.EvalCtx, mock.Anything).
				Return(tt.firstVal, tt.firstErr).
				Once()

			// Second and thereafter calls
			mockedClient.
				On("BooleanValue", ctx, flag.Name, flag.DefaultValue, flag.EvalCtx, mock.Anything).
				Return(tt.secondVal, tt.secondErr)

			lastValue := new(atomic.Bool)
			lastValue.Store(tt.defaultVal)

			callCount := new(atomic.Uint32)
			errCount := new(atomic.Uint32)

			go client.WatchBoolFlag(ctx, flag, time.NewTicker(time.Nanosecond), func(val bool, err error) {
				callCount.Add(1)

				if callCount.Load() > 2 {
					return
				}

				if err != nil {
					t.Logf("%s: %v", tt.name, err.Error())
					errCount.Add(1)
				} else {
					lastValue.Store(val)
				}
			})

			select {
			case <-ctx.Done():
			case <-time.After(time.Millisecond):
			}

			cancel()

			if lastValue.Load() != tt.targetVal {
				t.Errorf("expected value: %v, but the last value: %v", tt.targetVal, lastValue.Load())
			}

			var expectedErrors uint32
			if tt.firstErr != nil {
				expectedErrors++
			}

			if tt.secondErr != nil {
				expectedErrors++
			}

			if expectedErrors != errCount.Load() {
				t.Errorf("expected %d errors, got %d", expectedErrors, errCount.Load())
			}

			calls := int(callCount.Load())
			if tt.calledBack != (calls > 0) {
				t.Errorf("expected call back: %v, but called %v times", tt.calledBack, calls)
			}

			mockedClient.AssertExpectations(t)
		})
	}
}
