package tracker

import (
	"fmt"
	"github.com/jphastings/corviator/pkg/locator"
	"github.com/jphastings/corviator/pkg/math"
	"github.com/jphastings/corviator/pkg/sphere"
	"github.com/jphastings/corviator/pkg/tts"
)

type TrackerConfig struct {
	home      math.LLACoords
	sphere    *sphere.Config
	ttsEngine tts.TTSEngine
	Targets   chan *locator.TargetInstructions
}

func New(home math.LLACoords, ttsEngine tts.TTSEngine, sphereConfig *sphere.Config) *TrackerConfig {
	return &TrackerConfig{
		home:      home,
		sphere:    sphereConfig,
		ttsEngine: ttsEngine,
		Targets:   make(chan *locator.TargetInstructions, 1),
	}
}

func (track *TrackerConfig) Track() {
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

			desc := tts.Phrase(details.Name, bearing, isFirstTrack)
			fmt.Println(desc)

			err := track.ttsEngine.Speak(desc)
			if err != nil {
				panic(err)
			}

			//track.sphere.StepToDirection(bearing)
			isFirstTrack = false
		}
	}
}
