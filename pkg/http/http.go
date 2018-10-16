package http

import (
	"fmt"
	"github.com/jphastings/corviator/pkg/tracker"
	"net/http"
	"time"
)

func CorviatorAPI(port uint16, track *tracker.TrackerConfig) {
	router := http.NewServeMux()

	router.Handle("/focus", handleFocus(track))
	router.Handle("/config", handleConfig(track))

	webserver := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	webserver.ListenAndServe()
}
