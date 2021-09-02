package http

import (
	"context"
	"fmt"
	"github.com/jphastings/jan-poka/pkg/output/webmapper"
	"github.com/jphastings/jan-poka/pkg/tracker"
	"log"
	"net/http"
	"time"
)

const (
	readTimeout     = 5 * time.Second
	writeTimeout    = 10 * time.Second
	idleTimeout     = 15 * time.Second
	shutdownTimeout = 20 * time.Second
)

func WebAPI(port uint16, track *tracker.Config, includeMapper bool) func() {
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

	webserver.RegisterOnShutdown(announce(int(port)))

	go func() {
		for {
			err := webserver.ListenAndServe()
			if err == http.ErrServerClosed {
				break
			}
			log.Printf("😱 Webserver died, attempting to restart: %v\n", err)
		}
	}()

	return func() {
		log.Printf("🥱 Shutting down web server…")
		ctx, _ := context.WithTimeout(context.Background(), shutdownTimeout)
		_ = webserver.Shutdown(ctx)
	}
}
