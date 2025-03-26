package middleware_test

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	mocklog "github.com/Sengoku11/go-monorepo/mocks/github.com/Sengoku11/go-monorepo/pkg/logger"
	mw "github.com/Sengoku11/go-monorepo/pkg/middleware"
	chimware "github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/mock"
)

const (
	// URLs and HTTP method.
	httpURL    = "http://example.com/foo?bar=baz"
	httpsURL   = "https://example.com/foo?bar=baz"
	httpMethod = http.MethodGet

	// Expected parts of the URL.
	expectedPath  = "/foo"
	expectedQuery = "bar=baz"

	// Test Request IDs.
	httpRequestID   = "test-req-id"
	httpsRequestID  = "https-req-id"
	middlewareReqID = "middleware-req-id"

	// Field keys for log entries.
	fieldMethod    = "method"
	fieldPath      = "path"
	fieldQuery     = "query"
	fieldRequestID = "request_id"
	fieldDuration  = "duration_ms"
)

func TestBuildLogEntry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		url       string
		requestID string
		duration  int64
		useTLS    bool
	}{
		{
			name:      "HTTP",
			url:       httpURL,
			requestID: httpRequestID,
			duration:  150,
			useTLS:    false,
		},
		{
			name:      "HTTPS",
			url:       httpsURL,
			requestID: httpsRequestID,
			duration:  200,
			useTLS:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest(httpMethod, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.useTLS {
				req.TLS = new(tls.ConnectionState)
			}

			req.RequestURI = req.URL.RequestURI()
			// Inject a test request ID into the context.
			req = req.WithContext(context.WithValue(req.Context(), chimware.RequestIDKey, tt.requestID))

			msg, fields := mw.BuildLogEntry(req, tt.duration)
			expectedMsg := "GET " + tt.url + " HTTP/1.1"

			if msg != expectedMsg {
				t.Errorf("Expected msg %q, got %q", expectedMsg, msg)
			}

			expectedFields := []interface{}{
				fieldMethod, httpMethod,
				fieldPath, expectedPath,
				fieldQuery, expectedQuery,
				fieldRequestID, tt.requestID,
				fieldDuration, tt.duration,
			}

			if len(fields) != len(expectedFields) {
				t.Errorf("Expected %d fields, got %d", len(expectedFields), len(fields))
			}

			for i, v := range expectedFields {
				if fields[i] != v {
					t.Errorf("At index %d, expected %v, got %v", i, v, fields[i])
				}
			}
		})
	}
}

func TestLogRequestMiddleware(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		url              string
		requestID        string
		expectedResponse string
		useTLS           bool
	}{
		{
			name:             "HTTP",
			url:              httpURL,
			requestID:        middlewareReqID,
			expectedResponse: "OK",
			useTLS:           false,
		},
		{
			name:             "HTTPS",
			url:              httpsURL,
			requestID:        httpsRequestID,
			expectedResponse: "Secure OK",
			useTLS:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create the mock logger.
			mockedLog := mocklog.NewMockLogger(t)

			nextCalled := false
			// Stub next handler that writes a response.
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				nextCalled = true

				w.WriteHeader(http.StatusOK)

				_, err := w.Write([]byte(tt.expectedResponse))
				if err != nil {
					t.Fatal(err)
				}
			})

			// Wrap the next handler with the LogRequest middleware.
			handler := mw.LogRequest(mockedLog)(nextHandler)

			req, err := http.NewRequest(httpMethod, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.useTLS {
				req.TLS = new(tls.ConnectionState)
			}

			req.RequestURI = req.URL.RequestURI()
			// Inject a test request ID.
			req = req.WithContext(context.WithValue(req.Context(), chimware.RequestIDKey, tt.requestID))

			expectedMsg := "GET " + tt.url + " HTTP/1.1"
			mockedLog.
				On("Info",
					expectedMsg,
					fieldMethod, httpMethod,
					fieldPath, expectedPath,
					fieldQuery, expectedQuery,
					fieldRequestID, tt.requestID,
					fieldDuration, mock.AnythingOfType("int64")).
				Return().
				Once()

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			// Verify that the next handler was executed.
			if !nextCalled {
				t.Error("expected next handler to be called")
			}

			// Verify the HTTP response.
			if rr.Code != http.StatusOK {
				t.Errorf("expected status code 200, got %d", rr.Code)
			}

			if rr.Body.String() != tt.expectedResponse {
				t.Errorf("expected response body %q, got %q", tt.expectedResponse, rr.Body.String())
			}

			// Assert that the expectations were met.
			mockedLog.AssertExpectations(t)
		})
	}
}
