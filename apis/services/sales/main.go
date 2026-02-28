package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/ardanlabs/service/foundation/logger"
)

// This is used for determining what code was used to build a docker image / container. This helps us with debugging
var build = "develop"

// application entry point
func main() {
	// Create a logger
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT *******")
		},
	}

	traceIDFn := func(ctx context.Context) string {
		return "" //web.GetTraceID(ctx)
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "SALES", traceIDFn, events)

	// -------------------------------------------------------------------------------------------------

	// Create a context
	ctx := context.Background()

	// Call run
	if err := run(ctx, log); err != nil {
		// if run returns an error, terminate app
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	// It's always good to log the th num CPU's you have
	log.Info(ctx, "starting service", "GOMAXPROCS", runtime.GOMAXPROCS(0), "build", build)

	// We need to block this program from exiting

	// create a channel for os.Signal. This channel will only take OS signals
	shutdown := make(chan os.Signal, 1) // buffered

	// Notify, sends the SIGINT (interrupt) and SIGTERM (terminate) signals to the shutdown channel
	// It notifies the shutdown channel of signals from the OS.
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// We block on the channel.
	// This will only fire when we get an interrupt or terminate signal from the OS.
	// The two signals can also happen when we do CTRL C
	sig := <-shutdown

	log.Info(ctx, "Receive shutdown signal from OS. Shutdown started", "signal", sig)
	defer log.Info(ctx, "Receive shutdown signal from OS. Shutdown completed", "signal", sig)

	return nil
}
