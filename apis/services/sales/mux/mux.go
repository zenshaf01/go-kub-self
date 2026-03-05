package mux

import (
	"os"

	"github.com/ardanlabs/service/apis/services/api/mid"
	"github.com/ardanlabs/service/apis/services/sales/mux/route/sys/checkapi"
	"github.com/ardanlabs/service/foundation/logger"
	"github.com/ardanlabs/service/foundation/web"
)

func WebAPI(log *logger.Logger, shutdown chan os.Signal) *web.App {
	// Create the mux
	// The serve mux is a router implementation.
	// YOu can attach the handlers to the mux for each of the endpoints
	// The servemux's job is to:
	// - take an http request, see if there is a matching url
	// - See if it has a matching handler for the incoming url path
	// - Create a new goroutine and run that handler in that goroutine
	mux := web.NewApp(shutdown, mid.Logger(log))
	// Attach handler to mux with endpoint pattern
	checkapi.Routes(mux)
	return mux
}
