package tts

import (
	"github.com/jphastings/jan-poka/pkg/future"
	"github.com/jphastings/jan-poka/pkg/l10n"
	"github.com/jphastings/jan-poka/pkg/math"
)

func TrackedCallback(ttsEngine Engine) func(string, math.AERCoords, math.Meters, bool) future.Future {
	return func(name string, bearing math.AERCoords, distance math.Meters, isFirstTrack bool) future.Future {
		f := future.New()
		go func() {
			err := ttsEngine.Speak(l10n.Phrase(name, bearing, distance, isFirstTrack))
			if err == nil {
				f.Succeed()
			} else {
				f.Fail(err)
			}
		}()
		return f
	}
}
