package motor

import (
	"fmt"
	"periph.io/x/periph/conn/gpio"
	"time"
)

type PowerSaver struct {
	pin    gpio.PinOut
	leeway time.Duration

	resetTimer chan time.Duration
}

func NewPowerSaver(activePin gpio.PinOut, leeway time.Duration) *PowerSaver {
	return &PowerSaver{
		pin:        activePin,
		leeway:     leeway,
		resetTimer: make(chan time.Duration),
	}
}

func (mps *PowerSaver) PowerOn() error {
	fmt.Println("Powering up motors")
	return mps.pin.Out(gpio.High)
}

func (mps *PowerSaver) PowerOff() error {
	fmt.Println("Powering down motors")
	return mps.pin.Out(gpio.Low)
}
