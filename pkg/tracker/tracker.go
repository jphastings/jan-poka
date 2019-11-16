package tracker

import (
	"fmt"
	"github.com/jphastings/jan-poka/pkg/future"
	"github.com/jphastings/jan-poka/pkg/locator"
	"github.com/jphastings/jan-poka/pkg/math"
	"sync"
)

type Config struct {
	home      math.LLACoords
	callbacks []OnTracked
	Targets   chan *locator.TargetInstructions
}

type OnTracked func(string, math.AERCoords, math.Meters, bool) future.Future

func New(home math.LLACoords, callbacks ...OnTracked) *Config {
	return &Config{
		home:      home,
		callbacks: callbacks,
		Targets:   make(chan *locator.TargetInstructions, 1),
	}
}

func (track *Config) Track() {
	var (
		target       *locator.TargetInstructions
		tracker      <-chan locator.TargetDetails
		isFirstTrack bool
	)

	for {
		select {
		case target = <-track.Targets:
			isFirstTrack = true
			tracker = target.Poll()
		case details := <-tracker:
			bearing := track.home.DirectionTo(details.Coords, 0)
			distance := track.home.GreatCircleDistance(details.Coords)

			var wg sync.WaitGroup

			for _, callback := range track.callbacks {
				wg.Add(1)
				go func(cb OnTracked) {
					result := <-cb(details.Name, bearing, distance, isFirstTrack)
					if !result.IsOK() {
						fmt.Printf("could not present location: %v\n", result.Err)
					}
					wg.Done()
				}(callback)
			}

			wg.Wait()

			isFirstTrack = false
		}
	}
}
