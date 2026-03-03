package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	*http.ServeMux
	shutdown chan os.Signal
}

func NewApp(shutdown chan os.Signal) *App {
	a := &App{
		ServeMux: http.NewServeMux(),
		shutdown: shutdown,
	}
	return a
}

// Handle This handle function is our own handler which wraps the serve mux HandleFunc function
// This allows us to have our custom handlers to do any logic before or after the ServeMux
// HandleFunc. This is also necessary if we want to change the signature of our handlers.
func (a *App) Handle(path string, handler Handler) {
	// Define the serve mux handler
	h := func(w http.ResponseWriter, r *http.Request) {

		//Put logic before our handler. middleware

		if err := handler(r.Context(), w, r); err != nil {
			// error handling
			fmt.Println(err)
			return
		}

		// put logic after our handler middleware
	}
	a.ServeMux.HandleFunc(path, h)
}
