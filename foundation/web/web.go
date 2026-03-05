package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	*http.ServeMux
	shutdown chan os.Signal
	mw       []MidHandler
}

func NewApp(shutdown chan os.Signal, mw ...MidHandler) *App {
	a := &App{
		ServeMux: http.NewServeMux(),
		shutdown: shutdown,
		mw:       mw,
	}
	return a
}

// Handle This handle function is our own handler which wraps the serve mux HandleFunc function
// This allows us to have our custom handlers to do any logic before or after the ServeMux
// HandleFunc. This is also necessary if we want to change the signature of our handlers.
func (a *App) Handle(path string, handler Handler, mw ...MidHandler) {
	// since we have 2 sources of the middlewares we need to call this twice.
	// after both calls, the handler becomes nested under the middleware call chain.
	handler = wrapMiddleware(mw, handler)
	handler = wrapMiddleware(a.mw, handler)

	// Define the serve mux handler
	h := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			TraceID: uuid.NewString(),
			Tracer:  nil,
			Now:     time.Now().UTC(),
		}
		ctx := setValues(r.Context(), &v)

		if err := handler(ctx, w, r); err != nil {
			// error handling
			fmt.Println(err)
			return
		}
	}
	a.ServeMux.HandleFunc(path, h)
}
