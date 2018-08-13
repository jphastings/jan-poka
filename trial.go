package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/jphastings/corviator/pkg/baller"
	"github.com/jphastings/corviator/pkg/locations"
	"github.com/jphastings/corviator/pkg/transforms"
)

const facing float64 = 110

var currentHeading float64 = 0
var currentElevation float64 = 90

var home = transforms.LLACoords{
	Φ: 51.498842,
	Λ: -0.084357,
	A: 10,
}

func main() {
	locations, _ := locations.DecodeJSON([]byte(`{
		"poll": 5,
		"target": [
			{ "type": "iss" }
		]
	}`))

	for {
		select {
		case location := <-locations:
			admire(location)
		}
	}
}

func admire(target transforms.LLACoords) {
	distance, heading, elevation := transforms.RelativeDirection(home, target, facing)

	pointAt(heading, elevation)
	fmt.Printf("\nLooking %.0fm facing %.1fº up %.1fº\n\n", distance, heading, elevation)
}

func pointAt(heading, elevation float64) {
	var seq []baller.StepInstruction
	if heading == currentHeading {
		// Rotate the extra elevation to reach new direction, inverting the elevation direction
		dElevation := elevation - currentElevation
		seq = baller.StepSequenceΘ(heading, -dElevation)
	} else {
		seq = baller.StepSequenceHome(currentHeading, currentElevation)
		seq = append(seq, baller.StepSequenceElevation(heading, elevation)...)
	}

	for _, step := range seq {
		executeStepInstruction(step)
	}

	currentHeading = heading
	currentElevation = elevation
}

func executeStepInstruction(step baller.StepInstruction) {
	if step.IsForward {
		fmt.Printf(step.Name)
	} else {
		fmt.Printf(strings.ToLower(step.Name))
	}
	time.Sleep(step.WaitMilliseconds)
}
