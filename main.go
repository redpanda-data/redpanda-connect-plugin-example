package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	_ "github.com/Jeffail/benthos-plugin-example/condition"
	_ "github.com/Jeffail/benthos-plugin-example/input"
	_ "github.com/Jeffail/benthos-plugin-example/processor"
	"github.com/Jeffail/benthos/lib/api"
	"github.com/Jeffail/benthos/lib/config"
	"github.com/Jeffail/benthos/lib/log"
	"github.com/Jeffail/benthos/lib/manager"
	"github.com/Jeffail/benthos/lib/metrics"
	"github.com/Jeffail/benthos/lib/stream"
)

//------------------------------------------------------------------------------

var configPath = flag.String("c", "", "Path to a Benthos config file")

func main() {
	flag.Parse()

	conf := config.New()
	lints := []string{}

	if len(*configPath) > 0 {
		var err error
		if lints, err = config.Read(*configPath, true, &conf); err != nil {
			fmt.Fprintf(os.Stderr, "Configuration file read error: %v\n", err)
			os.Exit(1)
		}
	}

	logger := log.New(os.Stdout, conf.Logger)
	for _, lint := range lints {
		logger.Infoln(lint)
	}

	// Create our metrics type.
	stats, err := metrics.New(conf.Metrics, metrics.OptSetLogger(logger))
	if err != nil {
		logger.Errorf("Failed to connect to metrics aggregator: %v\n", err)
		os.Exit(1)
	}
	defer stats.Close()

	// Create HTTP API.
	httpServer, err := api.New("", "", conf.HTTP, conf, logger, stats)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create API: %v\n", err)
		os.Exit(1)
	}

	// Create resource manager.
	resourceMgr, err := manager.New(conf.Manager, httpServer, logger, stats)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create resource: %v\n", err)
		os.Exit(1)
	}

	dataStreamClosedChan := make(chan struct{})

	// Create stream pipeline.
	dataStream, err := stream.New(
		conf.Config,
		stream.OptSetManager(resourceMgr),
		stream.OptSetLogger(logger),
		stream.OptSetStats(stats),
		stream.OptOnClose(func() {
			close(dataStreamClosedChan)
		}),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Service closing due to: %v\n", err)
		os.Exit(1)
	}
	logger.Infoln("Launching a Benthos instance, use CTRL+C to close.")

	// Start HTTP server.
	httpServerClosedChan := make(chan struct{})
	go func() {
		logger.Infof(
			"Listening for HTTP requests at: %v\n",
			"http://"+conf.HTTP.Address,
		)
		httpErr := httpServer.ListenAndServe()
		if httpErr != nil && httpErr != http.ErrServerClosed {
			logger.Errorf("HTTP Server error: %v\n", httpErr)
		}
		close(httpServerClosedChan)
	}()

	var exitTimeout time.Duration
	if tout := conf.SystemCloseTimeout; len(tout) > 0 {
		var err error
		if exitTimeout, err = time.ParseDuration(tout); err != nil {
			logger.Errorf("Failed to parse shutdown timeout period string: %v\n", err)
			os.Exit(1)
		}
	}

	// Defer clean up.
	defer func() {
		go func() {
			httpServer.Shutdown(context.Background())
			select {
			case <-httpServerClosedChan:
			case <-time.After(exitTimeout / 2):
				logger.Warnln("Service failed to close HTTP server gracefully in time.")
			}
		}()

		go func() {
			<-time.After(exitTimeout + time.Second)
			logger.Warnln(
				"Service failed to close cleanly within allocated time." +
					" Exiting forcefully and dumping stack trace to stderr.",
			)
			pprof.Lookup("goroutine").WriteTo(os.Stderr, 1)
			os.Exit(1)
		}()

		if err := dataStream.Stop(exitTimeout); err != nil {
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for termination signal
	select {
	case <-sigChan:
		logger.Infoln("Received SIGTERM, the service is closing.")
	case <-dataStreamClosedChan:
		logger.Infoln("Pipeline outputs have terminated. Shutting down the service.")
	case <-httpServerClosedChan:
		logger.Infoln("HTTP Server has terminated. Shutting down the service.")
	}
}

//------------------------------------------------------------------------------
