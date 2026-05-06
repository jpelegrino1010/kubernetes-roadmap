package version

import (
	"fmt"
	"net/http"
	"os"
)

const defaultVersion = "version 2"

func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v := os.Getenv("APP_VERSION")
		if v == "" {
			v = defaultVersion
		}
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, v)
	}
}
