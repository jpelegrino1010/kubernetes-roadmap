package user

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// lister is the behavior the handler needs from a user store.
// Declared here at the consumer side, not on the implementer side.
type lister interface {
	All() []User
}

// Handler returns an http.HandlerFunc that writes the full user list as JSON.
func Handler(store lister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := store.All()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(users); err != nil {
			// At this point we have already committed the status header;
			// surface the error as a structured log entry at the boundary.
			fmt.Printf("user.Handler: encode: %v\n", err)
		}
	}
}
