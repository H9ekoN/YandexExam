package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateEndpoint(t *testing.T) {
	srv := NewServer(&Config{Port: "8080"})
	ts := httptest.NewServer(srv.router)
	defer ts.Close()

	tests := []struct {
		name          string
		method        string
		body          interface{}
		wantStatus    int
		wantResultKey string
		wantResultVal interface{}
	}{
		{
			name:          "basic addition",
			method:        http.MethodPost,
			body:          map[string]string{"expression": "2+2"},
			wantStatus:    http.StatusOK,
			wantResultKey: "result",
			wantResultVal: float64(4),
		},
		{
			name:          "complex calculation",
			method:        http.MethodPost,
			body:          map[string]string{"expression": "(2+3)*4"},
			wantStatus:    http.StatusOK,
			wantResultKey: "result",
			wantResultVal: float64(20),
		},
		{
			name:          "division by zero",
			method:        http.MethodPost,
			body:          map[string]string{"expression": "1/0"},
			wantStatus:    http.StatusUnprocessableEntity,
			wantResultKey: "error",
			wantResultVal: "Division by zero is not allowed",
		},
		{
			name:          "invalid method",
			method:        http.MethodGet,
			body:          nil,
			wantStatus:    http.StatusMethodNotAllowed,
			wantResultKey: "error",
			wantResultVal: "Method not allowed",
		},
		{
			name:          "invalid json",
			method:        http.MethodPost,
			body:          "invalid json",
			wantStatus:    http.StatusUnprocessableEntity,
			wantResultKey: "error",
			wantResultVal: "Invalid request body",
		},
		{
			name:          "empty expression",
			method:        http.MethodPost,
			body:          map[string]string{"expression": ""},
			wantStatus:    http.StatusUnprocessableEntity,
			wantResultKey: "error",
			wantResultVal: "Expression is not valid",
		},
		{
			name:          "invalid expression",
			method:        http.MethodPost,
			body:          map[string]string{"expression": "2++2"},
			wantStatus:    http.StatusUnprocessableEntity,
			wantResultKey: "error",
			wantResultVal: "Expression is not valid",
		},
		{
			name:          "numeric overflow",
			method:        http.MethodPost,
			body:          map[string]string{"expression": "999999999999999999999999999999*999999999999999999999999999999"},
			wantStatus:    http.StatusInternalServerError,
			wantResultKey: "error",
			wantResultVal: "Internal server error",
		},
		{
			name:          "large number multiplication",
			method:        http.MethodPost,
			body:          map[string]string{"expression": "1e308*2"},
			wantStatus:    http.StatusInternalServerError,
			wantResultKey: "error",
			wantResultVal: "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyBytes []byte
			var err error

			if tt.body != nil {
				switch v := tt.body.(type) {
				case string:
					bodyBytes = []byte(v)
				default:
					bodyBytes, err = json.Marshal(v)
					if err != nil {
						t.Fatalf("Failed to marshal request body: %v", err)
					}
				}
			}

			req, err := http.NewRequest(tt.method, fmt.Sprintf("%s/api/v1/calculate", ts.URL), bytes.NewBuffer(bodyBytes))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Want status %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			var result map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			got, exists := result[tt.wantResultKey]
			if !exists {
				t.Errorf("Want key %q in response", tt.wantResultKey)
			}
			if got != tt.wantResultVal {
				t.Errorf("Want %v, got %v", tt.wantResultVal, got)
			}
		})
	}
}
