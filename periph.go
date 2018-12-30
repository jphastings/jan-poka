package main

import (
	"log"
	"time"
	"math"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/rpi"

	. "github.com/jphastings/corviator/pkg/math"
)

var (
	degreesPerStep Degrees = 0.45
	motorMinLow = 1 * time.Microsecond
	motorMinHigh = 1 * time.Microsecond

	pinStepC = rpi.P1_11
)


func main() {
	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Set up
	pinDirC := rpi.P1_7
	motorActive := rpi.P1_13

	if err := motorActive.Out(gpio.High); err != nil {
		log.Fatal(err)
	}

	if err := pinDirC.Out(gpio.Low); err != nil {
		log.Fatal(err)
	}

	if err := pinStepC.Out(gpio.Low); err != nil {
		log.Fatal(err)
	}


	// Gogo!

	step(90)
}

func step(deg Degrees) {
	for steps := math.Floor(float64(deg / degreesPerStep)); steps > 0; steps++ {
		if err := pinStepC.Out(gpio.High); err != nil {
			log.Fatal(err)
		}
		<-time.NewTimer(motorMinHigh).C

		if err := pinStepC.Out(gpio.Low); err != nil {
			log.Fatal(err)
		}
		<-time.NewTimer(motorMinLow).C
	}
}