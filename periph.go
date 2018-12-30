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
	motorMinLow = 1 * time.Millisecond
	motorMinHigh = motorMinLow

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
	for {
		log.Println("Step 90º")
		step(90)

		log.Println("Waiting")
		<-time.NewTimer(5 * time.Second).C
	}
}

func step(deg Degrees) {
	steps := math.Floor(float64(deg / degreesPerStep))
	log.Println("Going steps:", steps)
	for ; steps > 0; steps-- {
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