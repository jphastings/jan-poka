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
	mps := &PowerSaver{
		pin:        activePin,
		leeway:     leeway,
		resetTimer: make(chan time.Duration),
	}
	go mps.runTimer()

	return mps
}

func (mps *PowerSaver) PowerUntil(powerDownIn time.Duration) error {
	fmt.Println("Powering up motors")
	if err := mps.pin.Out(gpio.High); err != nil {
		return err
	}

	mps.resetTimer <- powerDownIn + mps.leeway
	return nil
}

func (mps *PowerSaver) runTimer() {
	deactivationTimer := time.NewTimer(time.Duration(0))

	for {
		select {
		case powerDownIn := <-mps.resetTimer:
			deactivationTimer = time.NewTimer(powerDownIn)
		case <-deactivationTimer.C:
			fmt.Println("Powering down motors")
			if err := mps.pin.Out(gpio.Low); err != nil {
				// TODO: What to do with error?
			}
		}
	}
}
