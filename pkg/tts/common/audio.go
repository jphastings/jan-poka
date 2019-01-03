package common

import (
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"io"
	"time"
)

func Play(wavdata io.ReadCloser) error {
	s, format, err := wav.Decode(wavdata)
	if err != nil {
		return err
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(s)

	return nil
}
