package wheel

import "periph.io/x/periph/conn/gpio"

type Motor struct {
	StepChannel  chan bool
	AngleDegrees float64

	directionPin gpio.PinIO
	stepPin      gpio.PinIO
}

func New(angle float64, directionPin, stepPin gpio.PinIO) *Motor {
	return &Motor{
		StepChannel:  make(chan bool),
		AngleDegrees: angle,

		directionPin: directionPin,
		stepPin:      stepPin,
	}
}
