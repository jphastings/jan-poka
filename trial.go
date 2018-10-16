package main

import (
	"time"

	"github.com/jphastings/corviator/pkg/hardware/wheel"
	"github.com/jphastings/corviator/pkg/http"
	"github.com/jphastings/corviator/pkg/sphere"
	"github.com/jphastings/corviator/pkg/tracker"
	"github.com/jphastings/corviator/pkg/tts"

	. "github.com/jphastings/corviator/pkg/math"
)

var home = LLACoords{
	Φ: 51.498842,
	Λ: -0.084357,
	A: 10,
}

const facing = Degrees(120)

func main() {
	motors := []*wheel.Motor{
		wheel.New(0, nil, nil),
		wheel.New(120, nil, nil),
		wheel.New(240, nil, nil),
	}

	sphereConfig := sphere.New(
		motors, 200, 200.0/48.0,
		3*time.Millisecond, facing)

	ttsEngine, err := tts.NewGoogle()
	if err != nil {
		panic(err)
	}

	track := tracker.New(home, ttsEngine, sphereConfig)

	go track.Track()
	http.CorviatorAPI(2678, track)
}
