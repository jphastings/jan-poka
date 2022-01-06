package tracker

import (
	"fmt"
	"sync"
	"time"

	"github.com/jphastings/jan-poka/pkg/common"
	"github.com/jphastings/jan-poka/pkg/locator"
	"github.com/jphastings/jan-poka/pkg/math"
	"github.com/jphastings/jan-poka/pkg/temporal"
)

const (
	presenterTimeout              = 1 * time.Second
	greatCircleCutoff math.Meters = 5000
)

type Config struct {
	home      math.LLACoords
	callbacks map[string]common.OnTracked
	Targets   chan *locator.TargetInstructions
}

func New(home math.LLACoords, callbacks map[string]common.OnTracked) *Config {
	return &Config{
		home:      home,
		callbacks: callbacks,
		Targets:   make(chan *locator.TargetInstructions, 1),
	}
}

func (track *Config) Track() {
	var (
		target       *locator.TargetInstructions
		tracker      <-chan common.TargetDetails
		isFirstTrack bool
	)

	for {
		select {
		case target = <-track.Targets:
			isFirstTrack = true
			tracker = target.Poll()
		case details := <-tracker:
			bearing := track.home.DirectionTo(details.Coords, 0)

			unobstructedDistance := bearing.Range
			if details.Coords.Altitude < greatCircleCutoff {
				unobstructedDistance = track.home.GreatCircleDistance(details.Coords)
			}

			trackedDetails := common.TrackedDetails{
				Name:                 details.Name,
				AccurateAt:           details.AccurateAt,
				Target:               details.Coords,
				Bearing:              bearing,
				UnobstructedDistance: unobstructedDistance,
				IsFirstTrack:         isFirstTrack,
			}

			localTime, skyChanges, err := temporal.LocalTimeAndSkiesAt(details.Coords)
			if err == nil {
				trackedDetails.LocalTime = localTime
				trackedDetails.SkyChanges = skyChanges
			}

			callbacks := make(map[string]common.OnTracked)
			for name, callback := range track.callbacks {
				callbacks[name] = callback
			}
			if target.Requester != nil {
				callbacks["Requester"] = target.Requester
			}

			var wg sync.WaitGroup
			for presenter, callback := range callbacks {
				wg.Add(1)
				go func(presenter string, callback common.OnTracked, trackedDetails common.TrackedDetails) {
					select {
					case <-time.After(presenterTimeout):
						fmt.Printf("⚠️ timed out while trying to present to %s\n", presenter)
					case result := <-callback(trackedDetails):
						if !result.IsOK() {
							fmt.Printf("⚠️ could not present location with %s: %v\n", presenter, result.Err)
						}
					}

					wg.Done()
				}(presenter, callback, trackedDetails)
			}

			wg.Wait()

			isFirstTrack = false
		}
	}
}
