package mux

import (
	"net/http"

	"github.com/ardanlabs/service/apis/services/sales/mux/route/sys/checkapi"
)

func WebAPI() *http.ServeMux {
	// Create the mux
	// The serve mux is a router implementation.
	// YOu can attach the handlers to the mux for each of the endpoints
	mux := http.NewServeMux()
	// Attach handler to mux with endpoint pattern
	checkapi.Routes(mux)
	return mux
}
