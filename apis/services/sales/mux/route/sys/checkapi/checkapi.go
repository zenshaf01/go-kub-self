// Package checkapi is for system level routes that will be used by K8s
package checkapi

import (
	"context"
	"net/http"

	"github.com/ardanlabs/service/foundation/web"
)

func liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	s := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, s, http.StatusOK)
}

func readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	s := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, s, http.StatusOK)
}
