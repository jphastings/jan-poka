// +build rpi

package main

import (
	"fmt"
	"log"

	"periph.io/x/periph/host"
	"periph.io/x/periph/host/rpi"

	"github.com/jphastings/jan-poka/pkg/hardware/motor"
	"github.com/jphastings/jan-poka/pkg/sphere"
)

func init() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	if environment.UseSteppers {
		sphereConfig := sphere.New(
			[]*motor.Motor{
				motor.New(0, rpi.P1_7, rpi.P1_11),
			},
			rpi.P1_13, environment.MotorAutoSleepLeeway,
			environment.MotorSteps,
			float64(environment.SphereDiameter/environment.OmniwheelDiameter),
			environment.MinStepInterval,
			environment.Heading,
		)
		callbacks = append(callbacks, sphereConfig.TrackerCallback)
		fmt.Println("Stepper motor tracking: on")
	} else {
		fmt.Println("Stepper motor tracking: off")
	}
}
