package baller

import (
	"math"
	"sort"
	"time"
)

type Motor struct {
	Name           string
	RadialPosition float64
}

type StepInstruction struct {
	Name             string
	IsForward        bool
	WaitMilliseconds time.Duration
}

type stepTime struct {
	Name      string
	IsForward bool
	At        time.Duration
}

var motors []Motor = []Motor{
	{"A", 0},
	{"B", 120},
	{"C", 240},
}

const omniwheelToBallRatio = 200.0 / 48.0
const stepsPerRotation = 200
const minWaitBetweenSteps = 3 * time.Millisecond

// Home is at Θ = 0 (strait up)
func StepSequenceHome(currentHeading, currentElevation float64) []StepInstruction {
	oppositeHeading := 180 + currentHeading
	if oppositeHeading >= 360 {
		oppositeHeading -= 360
	}

	return StepSequenceElevation(oppositeHeading, currentElevation)
}

func StepSequenceElevation(heading, elevation float64) []StepInstruction {
	return StepSequenceΘ(heading, 90-elevation)
}

func StepSequenceΘ(heading, θ float64) []StepInstruction {
	steps := θ * omniwheelToBallRatio * stepsPerRotation / 360

	travelTime := time.Duration(steps) * minWaitBetweenSteps

	var stepTimes []stepTime
	for _, motor := range motors {
		motorSteps := -int(math.Round(math.Cos((heading-motor.RadialPosition)*math.Pi/180.0) * steps))

		if motorSteps == 0 {
			continue
		}

		isForward := motorSteps >= 0
		if !isForward {
			motorSteps = -motorSteps
		}

		pauseBetween := travelTime / time.Duration(motorSteps)

		for i := 0; i < motorSteps; i++ {
			stepTimes = append(stepTimes, stepTime{
				Name:      motor.Name,
				IsForward: isForward,
				At:        (time.Duration(i) * pauseBetween).Round(10 * time.Microsecond),
			})
		}
	}

	sort.Slice(stepTimes, func(i, j int) bool {
		return stepTimes[i].At < stepTimes[j].At
	})

	clock := time.Duration(0)
	var instructions []StepInstruction
	for _, step := range stepTimes {
		wait := step.At - clock
		clock = step.At

		instructions = append(instructions, StepInstruction{
			Name:             step.Name,
			IsForward:        step.IsForward,
			WaitMilliseconds: wait,
		})
	}

	return instructions
}
