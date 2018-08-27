package wheel

import (
	"fmt"
	. "github.com/jphastings/corviator/pkg/math"
	"periph.io/x/periph/conn/gpio"
	"strings"
)

type Motor struct {
	StepChannel chan bool
	Angle       Degrees

	directionPin gpio.PinIO
	stepPin      gpio.PinIO
}

func New(angle float64, directionPin, stepPin gpio.PinIO) *Motor {
	m := &Motor{
		StepChannel: make(chan bool),
		Angle:       angle,

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
