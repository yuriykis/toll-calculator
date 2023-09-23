package main

import (
	"net"
	"net/http"
	"os"

	log "github.com/go-kit/log"
	"github.com/yuriykis/tolling/go-kit-example/aggsvc/aggendpoint"
	"github.com/yuriykis/tolling/go-kit-example/aggsvc/aggservice"
	"github.com/yuriykis/tolling/go-kit-example/aggsvc/aggtransport"
)

func main() {
	var (
		logger      = log.NewLogfmtLogger(os.Stderr)
		service     = aggservice.New()
		endpoints   = aggendpoint.New(service, logger)
		httpHandler = aggtransport.NewHTTPHandler(endpoints, logger)
	)

	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	// The HTTP listener mounts the Go kit HTTP handler we created.
	httpListener, err := net.Listen("tcp", ":4000")
	if err != nil {
		logger.Log("transport", "HTTP", "during", "Listen", "err", err)
		os.Exit(1)
	}
	err = http.Serve(httpListener, httpHandler)
	if err != nil {
		logger.Log("transport", "HTTP", "during", "Serve", "err", err)
		os.Exit(1)
	}
}
