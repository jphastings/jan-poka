package tts

import (
	"github.com/jphastings/jan-poka/pkg/common"
	"github.com/jphastings/jan-poka/pkg/future"
	"github.com/jphastings/jan-poka/pkg/l10n"
)

func TrackedCallback(ttsEngine Engine) common.OnTracked {
	return func(details common.TrackedDetails) future.Future {
		return future.Exec(func() error {
			return ttsEngine.Speak(l10n.Phrase(details.Name, details.Bearing, details.UnobstructedDistance, details.IsFirstTrack))
		})
	}
}
