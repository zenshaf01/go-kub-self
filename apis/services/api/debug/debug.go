package debug

import (
	"expvar"
	"net/http"
	"net/http/pprof"

	"github.com/arl/statsviz"
)

func Mux() *http.ServeMux {
	// The purpose of the ServeMux is to define the endpoints and their handlers.
	// Never use http.DefaultServeMux in prod. Any imported package can bind to the http.DefaultServeMux.
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler()) // gives us metrics

	// Package statsviz allows visualizing Go runtime metrics data in real time in your browser.
	err := statsviz.Register(mux)
	if err != nil {
		return nil
	}

	return mux
}
