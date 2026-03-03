// Package checkapi is for system level routes that will be used by K8s
package checkapi

import (
	"context"
	"encoding/json"
	"net/http"
)

func liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	s := struct {
		Status string
	}{
		Status: "OK",
	}

	json.NewEncoder(w).Encode(s)
	return nil
}

func readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	s := struct {
		Status string
	}{
		Status: "OK",
	}

	json.NewEncoder(w).Encode(s)
	return nil
}
