package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIEndpoints(t *testing.T) {
	handler, err := newHandler()
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedState  string
	}{
		{
			name:           "Health endpoint",
			path:           "/api/health",
			expectedStatus: http.StatusOK,
			expectedState:  "ok",
		},
		{
			name:           "Status endpoint",
			path:           "/api/status",
			expectedStatus: http.StatusOK,
			expectedState:  "online",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(
				http.MethodGet,
				test.path,
				nil,
			)

			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)

			if response.Code != test.expectedStatus {
				t.Fatalf(
					"expected status %d, got %d",
					test.expectedStatus,
					response.Code,
				)
			}

			var result apiResponse

			if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
				t.Fatalf("invalid JSON response: %v", err)
			}

			if result.Service != "lc-share" {
				t.Errorf("unexpected service: %s", result.Service)
			}

			if result.Version != appVersion {
				t.Errorf("unexpected version: %s", result.Version)
			}

			if result.Status != test.expectedState {
				t.Errorf(
					"expected state %s, got %s",
					test.expectedState,
					result.Status,
				)
			}
		})
	}
}

func TestAPIRejectsUnsupportedMethods(t *testing.T) {
	handler, err := newHandler()
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	endpoints := []string{
		"/api/health",
		"/api/status",
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			request := httptest.NewRequest(
				http.MethodPost,
				endpoint,
				nil,
			)

			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)

			if response.Code != http.StatusMethodNotAllowed {
				t.Errorf(
					"expected 405, got %d",
					response.Code,
				)
			}
		})
	}
}
