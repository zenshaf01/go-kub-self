package checkapi

import (
	"github.com/ardanlabs/service/foundation/web"
)

func Routes(mux *web.App) {
	mux.Handle("GET /liveness", liveness)
	mux.Handle("GET /readiness", readiness)
}
