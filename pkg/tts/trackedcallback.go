package tts

import (
	"github.com/jphastings/corviator/pkg/future"
	"github.com/jphastings/corviator/pkg/l10n"
	"github.com/jphastings/corviator/pkg/math"
)

func TrackedCallback(ttsEngine Engine) func(string, math.AERCoords, bool) future.Future {
	return func(name string, bearing math.AERCoords, isFirstTrack bool) future.Future {
		f := future.New()
		go func() {
			err := ttsEngine.Speak(l10n.Phrase(name, bearing, isFirstTrack))
			if err == nil {
				f.Succeed()
			} else {
				f.Fail(err)
			}
		}()
		return f
	}
}
