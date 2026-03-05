package mid

import (
	"context"
	"net/http"

	"github.com/ardanlabs/service/foundation/logger"
	"github.com/ardanlabs/service/foundation/web"
)

// Logger We need a function which is a middleware function.
// A middleware function is a function which take a web handler as a parameter and return a web handler.
// The below uses a closure to pass in the logger to teh middleware function
func Logger(logger *logger.Logger) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		// See how it has wrapped the passed in handler with the new handler which does its own thing before and after calling the passed in handler.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// Do something here
			logger.Info(ctx, "request started", r.Method, "path", r.URL.Path, "remoteAddress", r.RemoteAddr)
			err := handler(ctx, w, r)
			// Do something here
			logger.Info(ctx, "request completed", r.Method, "path", r.URL.Path, "remoteAddress", r.RemoteAddr)
			return err
		}
		return h
	}
	return m
}
