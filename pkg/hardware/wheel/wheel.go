package wheel

import (
	"fmt"
	"periph.io/x/periph/conn/gpio"
	"strings"
)

type Motor struct {
	StepChannel  chan bool
	AngleDegrees float64

	directionPin gpio.PinIO
	stepPin      gpio.PinIO
}

func New(angle float64, directionPin, stepPin gpio.PinIO) *Motor {
	m := &Motor{
		StepChannel:  make(chan bool),
		AngleDegrees: angle,

		directionPin: directionPin,
		stepPin:      stepPin,
	}

	i := int(angle / 120)
	go m.Move(string(rune(i + 65)))

	return m
}

func (m *Motor) Move(name string) {
	for {
		select {
		case isForward := <-m.StepChannel:
			if isForward {
				// directionPin = gpio.High
				fmt.Printf(name)
			} else {
				// directionPin = gpio.Low
				fmt.Printf(strings.ToLower(name))
			}
		}
	}
}
