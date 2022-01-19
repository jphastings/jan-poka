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

	mdns.Register("API", "_http._tcp", int(port))

	shutdown.Ensure("Webserver", func() error {
		if webserver == nil {
			return nil
		}
		return webserver.Shutdown(context.Background())
	})
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
