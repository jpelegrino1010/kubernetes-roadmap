package health_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/app-go/internal/health"
)

func TestHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		method         string
		wantStatusCode int
		wantStatus     string
	}{
		{
			name:           "returns 200 with status ok",
			method:         http.MethodGet,
			wantStatusCode: http.StatusOK,
			wantStatus:     "ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tt.method, "/health", nil)
			w := httptest.NewRecorder()

			health.Handler().ServeHTTP(w, req)

			if w.Code != tt.wantStatusCode {
				t.Fatalf("want status %d, got %d", tt.wantStatusCode, w.Code)
			}

			var got struct {
				Status string `json:"status"`
			}
			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("decode body: %v", err)
			}

			if got.Status != tt.wantStatus {
				t.Errorf("want status %q, got %q", tt.wantStatus, got.Status)
			}
		})
	}
}
