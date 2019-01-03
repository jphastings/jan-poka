package tts

import (
	"github.com/jphastings/corviator/pkg/l10n"
	. "github.com/jphastings/corviator/pkg/math"
)

type Engine interface {
	Speak(string) error
}

func TrackedCallback(ttsEngine Engine) func(string, AERCoords, bool) error {
	return func(name string, bearing AERCoords, isFirstTrack bool) error {
		return ttsEngine.Speak(l10n.Phrase(name, bearing, isFirstTrack))
	}
}
