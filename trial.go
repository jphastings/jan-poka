package main

import (
	"fmt"
	"time"

	"github.com/jphastings/corviator/pkg/hardware/wheel"
	"github.com/jphastings/corviator/pkg/locations"
	"github.com/jphastings/corviator/pkg/math"
	"github.com/jphastings/corviator/pkg/sphere"
	"github.com/jphastings/corviator/pkg/transforms"
)

var home = math.LLACoords{
	Φ: 51.498842,
	Λ: -0.084357,
	A: 10,
}

func main() {
	target, _ := locations.DecodeJSON([]byte(`{
		"poll": 10,
		"target": [
			{"type": "iss", "lat": 47.9520658, "long": 7.9562333, "alt": 0}
		]
	}`))

	motors := []*wheel.Motor{
		wheel.New(0, nil, nil),
		wheel.New(120, nil, nil),
		wheel.New(240, nil, nil),
	}

	s := sphere.New(
		motors,
		200,
		200.0/48.0,
		3*time.Millisecond,
		0)

	locations := target.Poll()

	for {
		select {
		case location := <-locations:
			admire(s, location)
		}
	}
}

func admire(s *sphere.Config, target math.LLACoords) {
	distance, heading, elevation := transforms.RelativeDirection(home, target, s.Facing)

	fmt.Printf("\n\nLooking at (%.2f, %.2f) which is: %.0fm facing %.1fº up %.1fº\n", target.Φ, target.Λ, distance, heading, elevation)
	_ = pointAt(s, heading, elevation)
}

func pointAt(s *sphere.Config, heading, elevation float64) time.Duration {
	return s.StepToElevation(heading, elevation)
}
