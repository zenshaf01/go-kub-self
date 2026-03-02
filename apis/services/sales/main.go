package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/ardanlabs/service/apis/services/api/debug"
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

	// ---------------------------------------------------------------------------------------------
	// Configuration
	// Import conf package in main and set up the configuration for the project
	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout        time.Duration `conf:"default:5s"`
			WriteTimeout       time.Duration `conf:"default:5s"`
			IdleTimeout        time.Duration `conf:"default:120s"`
			ShutdownTimeout    time.Duration `conf:"default:20s"`
			APIHost            string        `conf:"default:0.0.0.0:3000"`
			DebugHost          string        `conf:"default:0.0.0.0:3010"`
			CORSAllowedOrigins []string      `conf:"default:*,mask"` // use mask / noprint if you want to hide the value while printing config
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "Sales",
		},
	}

	// Parse config
	const prefix = "SALES"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}
	log.Info(ctx, "service starting with version: ", cfg.Build)
	// ---------------------------------------------------------------------------------------------
	// Start debug service
	// As a general rule, the parent goroutine should not terminate before child goroutines.
	// You should not have orphan goroutines.
	// The main goroutine is the parent and the one below is child.
	go func() {
		log.Info(ctx, "starting api service", "debug v1 router started")

		// The listenAndServe call is blocking
		// The mux has the url to handler mapping
		if err = http.ListenAndServe(cfg.Web.DebugHost, debug.Mux()); err != nil {
			log.Error(ctx, "shutdown", "status", "http debug server closed", "error", err)
		}
	}()

	// ---------------------------------------------------------------------------------------------
	// Build a string from config to be written to stdout
	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config: %w", err)
	}
	log.Info(ctx, "startup config: ", out)
	expvar.NewString("build").Set(cfg.Build)

	// create a channel for os.Signal. This channel will only take OS signals
	shutdown := make(chan os.Signal, 1) // buffered

	// Notify, sends the SIGINT (interrupt) and SIGTERM (terminate) signals to the shutdown channel
	// It notifies the shutdown channel of signals from the OS.
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// We need to block this program from exiting
	// We block on the channel.
	// This will only fire when we get an interrupt or terminate signal from the OS.
	// The two signals can also happen when we do CTRL C
	//sig := <-shutdown
	//
	//log.Info(ctx, "Receive shutdown signal from OS. Shutdown started", "signal", sig)
	//defer log.Info(ctx, "Receive shutdown signal from OS. Shutdown completed", "signal", sig)

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      nil,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)

		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
