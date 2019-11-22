// +build rpi
package main

import (
	"fmt"
	"github.com/jphastings/jan-poka/pkg/hardware/stepper"
	"github.com/jphastings/jan-poka/pkg/math"
	"log"
	"periph.io/x/periph/host"
	_ "periph.io/x/periph/host/rpi"
	"time"
)

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	steppers := stepper.Pi2Quad(stepper.Motors["28BYJ-48"])
	s := steppers[0]
	s.SetSpeed(20)

	s.Off()
	log.Println("Off")
	<-time.NewTimer(1 * time.Second).C

	angle := math.Degrees(0)
	for {
		angle = math.ModDeg(angle + 180)
		start := time.Now()
		s.SetAngle(angle)
		duration := time.Now().Sub(start)
		fmt.Printf("Took: %v\n", duration)
		<-time.NewTimer(1 * time.Second).C
	}
}
