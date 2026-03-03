package mux

import (
	"encoding/json"
	"net/http"
)

func WebAPI() *http.ServeMux {
	// Create the mux
	// The serve mux is a router implementation.
	// YOu can attach the handlers to the mux for each of the endpoints
	mux := http.NewServeMux()

	// Create the handler
	h := func(w http.ResponseWriter, r *http.Request) {
		s := struct {
			Status string
		}{
			Status: "OK",
		}

		json.NewEncoder(w).Encode(s)
	}

	// Attach handler to mux with endpoint pattern
	mux.HandleFunc("GET /test", h)

	return mux
}
