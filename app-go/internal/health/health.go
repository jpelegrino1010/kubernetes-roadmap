// Package health provides the HTTP handler for liveness checks.
package health

import (
	"encoding/json"
	"net/http"
)

// response is the JSON shape returned by the health endpoint.
type response struct {
	Status string `json:"status"`
}

// Handler returns an http.HandlerFunc that replies with {"status":"ok"}.
func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := response{Status: "ok"}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			// Encoding to a ResponseWriter rarely fails; nothing useful to do here
			// beyond having already written the header. Log at the boundary instead.
			_ = err
		}
	}
}
