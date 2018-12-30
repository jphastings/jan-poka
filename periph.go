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

	pinDirC := rpi.P1_7
	pinStepC := rpi.P1_11
	motorActive := rpi.P1_13

	log.Println("Sleep is high")
	if err := motorActive.Out(gpio.High); err != nil {
		log.Fatal(err)
	}

	if err := pinDirC.Out(gpio.Low); err != nil {
		log.Fatal(err)
	}

	t := time.NewTicker(10 * time.Microsecond)
	for {
		log.Println("Go Low")
		if err := pinStepC.Out(gpio.Low); err != nil {
			log.Fatal(err)
		}
		<-t.C

		log.Println("Go High")
		if err := pinStepC.Out(gpio.High); err != nil {
			log.Fatal(err)
		}

		<-time.NewTimer(5 * time.Microsecond).C
	}
}