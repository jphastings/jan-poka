package http

import (
	"encoding/json"
	"github.com/jphastings/jan-poka/pkg/common"
	"github.com/jphastings/jan-poka/pkg/future"
	"github.com/jphastings/jan-poka/pkg/l10n"
	"github.com/jphastings/jan-poka/pkg/locator"
	"github.com/jphastings/jan-poka/pkg/tracker"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const resultTimeout = time.Second

func handleFocus(track *tracker.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		target, err := locator.DecodeJSON(body)
		if err != nil {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			_, _ = w.Write([]byte(`{"error":"couldn't decode JSON request"}`))
			return
		}

		onTracked, cb := onTrackedChannel()
		target.Requester = cb

		select {
		case track.Targets <- target:
			select {
			case <-time.After(resultTimeout):
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write([]byte(`{"message":"Location calculation took a while, calculating in the background"}`))
			case loc := <-onTracked:
				enc, err := json.Marshal(loc)
				if err != nil {
					w.WriteHeader(http.StatusAccepted)
					_, _ = w.Write([]byte(`{"message":"Location couldn't be encoded for response, but succeeded."}`))
				} else {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write(enc)
				}
			}

		default:
			w.WriteHeader(http.StatusTooManyRequests)
		}
		// No need to process future locations for HTTP responses
		onTracked = nil
	})
}

type trackedResponse struct {
	Name                 string    `json:"name"`
	Latitude             float64   `json:"latitude"`
	Longitude            float64   `json:"longitude"`
	Altitude             float64   `json:"altitude"`
	UnobstructedDistance float64   `json:"unobstructedDistance"`
	AccurateAt           time.Time `json:"accurateAt"`
	Summary              string    `json:"summary"`
}

func onTrackedChannel() (chan trackedResponse, common.OnTracked) {
	c := make(chan trackedResponse)
	return c, func(details common.TrackedDetails) future.Future {
		return future.Exec(func() error {
			if c == nil {
				return nil
			}

			summary := l10n.Phrase(details.Name, details.Bearing, details.UnobstructedDistance, false)

			resp := trackedResponse{
				Name:                 details.Name,
				Latitude:             float64(details.Target.Latitude),
				Longitude:            float64(details.Target.Longitude),
				Altitude:             float64(details.Target.Altitude),
				UnobstructedDistance: float64(details.UnobstructedDistance),
				AccurateAt:           details.AccurateAt,
				Summary:              summary,
			}

			// Allow for it to have been nilled
			select {
			case c <- resp:
			default:
			}

			return nil
		})
	}
}

func handleConfig(track *tracker.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})
}
