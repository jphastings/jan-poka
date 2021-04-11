package tts

import (
	"github.com/jphastings/jan-poka/pkg/future"
	"github.com/jphastings/jan-poka/pkg/l10n"
	"github.com/jphastings/jan-poka/pkg/math"
	"github.com/jphastings/jan-poka/pkg/tracker"
)

func TrackedCallback(ttsEngine Engine) tracker.OnTracked {
	return func(name string, _ math.LLACoords, bearing math.AERCoords, distance math.Meters, isFirstTrack bool) future.Future {
		return future.Exec(func() error {
			return ttsEngine.Speak(l10n.Phrase(name, bearing, distance, isFirstTrack))
		})
	}
}
