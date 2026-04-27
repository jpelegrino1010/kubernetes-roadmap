package user_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/app-go/internal/user"
)

func TestStore_All(t *testing.T) {
	t.Parallel()

	store := user.New()
	got := store.All()

	if len(got) == 0 {
		t.Fatal("want at least one user, got none")
	}

	// Mutation of the returned slice must not affect the store.
	got[0].Name = "mutated"
	fresh := store.All()
	if fresh[0].Name == "mutated" {
		t.Error("store.All returned a reference to internal state; want a copy")
	}
}

func TestHandler_Users(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		wantStatusCode int
		wantMinUsers   int
	}{
		{
			name:           "returns 200 with user list",
			wantStatusCode: http.StatusOK,
			wantMinUsers:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			store := user.New()
			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			w := httptest.NewRecorder()

			user.Handler(store).ServeHTTP(w, req)

			if w.Code != tt.wantStatusCode {
				t.Fatalf("want status %d, got %d", tt.wantStatusCode, w.Code)
			}

			var users []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				Role string `json:"role"`
			}
			if err := json.NewDecoder(w.Body).Decode(&users); err != nil {
				t.Fatalf("decode body: %v", err)
			}

			if len(users) < tt.wantMinUsers {
				t.Errorf("want at least %d user(s), got %d", tt.wantMinUsers, len(users))
			}
		})
	}
}
