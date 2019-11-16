package motor

import (
	"fmt"
	. "github.com/jphastings/jan-poka/pkg/math"
	"periph.io/x/periph/conn/gpio"
	"time"
)

type Motor struct {
	StepChannel chan bool
	Angle       Degrees

	directionPin gpio.PinOut
	stepPin      gpio.PinOut
}

var (
	pulsePeriod     = 100 * time.Microsecond
	directionLevels = map[bool]gpio.Level{true: gpio.High, false: gpio.Low}
)

func New(angle Degrees, directionPin, stepPin gpio.PinIO) *Motor {
	m := &Motor{
		StepChannel: make(chan bool),
		Angle:       angle,

		directionPin: directionPin,
		stepPin:      stepPin,
	}

	i := int(angle / 120)
	go m.move(string(rune(i + 65)))

	return m
}

func (m *Motor) move(name string) {
	for {
		select {
		case isForward := <-m.StepChannel:
			// TODO: Looks like we're sleeping before the roundtrip is done
			if err := m.directionPin.Out(directionLevels[isForward]); err == nil {
				if err := m.stepPin.Out(gpio.High); err == nil {
					<-time.NewTimer(pulsePeriod).C
					if err := m.stepPin.Out(gpio.Low); err != nil {
						fmt.Printf("could not turn stepper controller %s low: %v\n", name, err)
					}
				}
			} else {
				fmt.Printf("could not turn stepper controller %s high: %v\n", name, err)
			}
		}
	}
}
