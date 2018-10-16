package http

import (
	"github.com/jphastings/corviator/pkg/tracker"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jphastings/corviator/pkg/locator"
)

func handleFocus(track *tracker.TrackerConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}
		target, err := locator.DecodeJSON(body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		select {
		case track.Targets <- target:
			w.WriteHeader(http.StatusAccepted)
		default:
			w.WriteHeader(http.StatusTooManyRequests)
		}
	})
}

func handleConfig(track *tracker.TrackerConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})
}
