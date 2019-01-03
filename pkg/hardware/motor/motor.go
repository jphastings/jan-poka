package motor

import (
	. "github.com/jphastings/corviator/pkg/math"
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
	pulsePeriod     = 25 * time.Microsecond
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
	go m.Move(string(rune(i + 65)))

	return m
}

func (m *Motor) Move(name string) {
	for {
		select {
		case isForward := <-m.StepChannel:
			if err := m.directionPin.Out(directionLevels[isForward]); err == nil {
				if err := m.stepPin.Out(gpio.High); err == nil {
					<-time.NewTimer(pulsePeriod).C
					if err := m.stepPin.Out(gpio.Low); err != nil {
						// TODO: What happens with an error?
					}
				}
			} else {
				// TODO: What happens with an error?
			}
		}
	}
}
