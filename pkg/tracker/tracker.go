package tracker

import (
	"github.com/jphastings/corviator/pkg/locator"
	"github.com/jphastings/corviator/pkg/math"
	. "github.com/jphastings/corviator/pkg/math"
)

type Config struct {
	home      math.LLACoords
	callbacks []OnTracked
	Targets   chan *locator.TargetInstructions
}

type OnTracked func(string, AERCoords, bool) chan error

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

			for _, callback := range track.callbacks {
				// TODO: deal with errors
				_ = callback(details.Name, bearing, isFirstTrack)
			}

			isFirstTrack = false
		}
	}
}
