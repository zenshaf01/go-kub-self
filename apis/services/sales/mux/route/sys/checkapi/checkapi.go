// Package checkapi is for system level routes that will be used by K8s
package checkapi

import (
	"encoding/json"
	"net/http"
)

func liveness(w http.ResponseWriter, r *http.Request) {
	s := struct {
		Status string
	}{
		Status: "OK",
	}

	json.NewEncoder(w).Encode(s)
}

func readiness(w http.ResponseWriter, r *http.Request) {
	s := struct {
		Status string
	}{
		Status: "OK",
	}

	json.NewEncoder(w).Encode(s)
}
