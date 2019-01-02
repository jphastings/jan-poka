package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/rpi"

	. "github.com/jphastings/corviator/pkg/math"
)

var (
	degreesPerStep Degrees = 0.45
	motorMinLow            = 150 * time.Microsecond
	motorMinHigh           = 10 * time.Microsecond

	motorActive = rpi.P1_13
	pinStepC    = rpi.P1_11
	pinDirC     = rpi.P1_7
)

func main() {
	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	motorState(false)

	if err := pinStepC.Out(gpio.Low); err != nil {
		log.Fatal(err)
	}

	// Gogo!
	for {
		degs := Degrees(rand.Float64() * 720)

		if err := pinDirC.Out(gpio.Low); err != nil {
			log.Fatal(err)
		}

		log.Println("Step ", degs, "º")
		step(degs)

		<-time.NewTimer(1 * time.Second).C

		log.Println("…and back")
		if err := pinDirC.Out(gpio.High); err != nil {
			log.Fatal(err)
		}

		step(degs)
		<-time.NewTimer(10 * time.Second).C
	}
}

func step(deg Degrees) {
	steps := int(math.Floor(float64(deg / degreesPerStep)))
	log.Println("Going steps:", steps)

	motorState(true)
	<-time.NewTimer(time.Millisecond).C
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
	<-time.NewTimer(time.Millisecond).C
	motorState(false)
}

func motorState(on bool) {
	var level gpio.Level

	if on {
		level = gpio.High
	} else {
		level = gpio.Low
	}

	if err := motorActive.Out(level); err != nil {
		log.Fatal(err)
	}
}
