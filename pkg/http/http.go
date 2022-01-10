package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jphastings/jan-poka/pkg/mdns"
	"github.com/jphastings/jan-poka/pkg/output/webmapper"
	"github.com/jphastings/jan-poka/pkg/shutdown"
	"github.com/jphastings/jan-poka/pkg/tracker"
)

const (
	readTimeout     = 5 * time.Second
	writeTimeout    = 10 * time.Second
	idleTimeout     = 15 * time.Second
	shutdownTimeout = 20 * time.Second
)

func WebAPI(port uint16, track *tracker.Config, includeMapper bool) {
	router := http.NewServeMux()

	router.Handle("/focus", handleFocus(track))
	router.Handle("/config", handleConfig(track))
	if includeMapper && false {
		router.Handle("/", webmapper.Handler())
	}

	webserver := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	shutdownMDNS, _ := mdns.Register("_http._tcp", int(port))
	webserver.RegisterOnShutdown(func() {
		if err := shutdownMDNS(); err != nil {
			log.Printf("⚠️ Failed to shutdown mDNS server for HTTP: %v", err)
		}
	})

	shutdown.Ensure("Webserver", func() error { return webserver.Shutdown(context.Background()) })
	go func() {
		for {
			err := webserver.ListenAndServe()
			if err == http.ErrServerClosed {
				break
			}
			log.Printf("😱 Webserver died, attempting to restart: %v\n", err)
		}
	}()
}
