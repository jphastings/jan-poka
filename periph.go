package main

import (
	"log"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/rpi"
)

func main() {
	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	motorNotSleep := rpi.P1_27
	pinDirC := rpi.P1_7
	pinStepC := rpi.P1_11

	if err := motorNotSleep.Out(gpio.High); err != nil {
		log.Fatal(err)
	}

	if err := pinDirC.Out(gpio.Low); err != nil {
		log.Fatal(err)
	}

	t := time.NewTicker(time.Second)
	for {
		if err := pinStepC.Out(gpio.Low); err != nil {
			log.Fatal(err)
		}
		<-t.C

		if err := pinStepC.Out(gpio.High); err != nil {
			log.Fatal(err)
		}

		<-time.NewTimer(200 * time.Millisecond).C
	}
}