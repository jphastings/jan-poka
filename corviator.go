package main

import (
	"fmt"
	"github.com/jphastings/corviator/pkg/env"
	"github.com/jphastings/corviator/pkg/hardware/motor"
	"github.com/jphastings/corviator/pkg/l10n"
	"github.com/jphastings/corviator/pkg/sphere"
	"github.com/jphastings/corviator/pkg/tts"
	"github.com/jphastings/corviator/pkg/tts/googletts"
	"log"
	"periph.io/x/periph/host/rpi"

	"github.com/jphastings/corviator/pkg/http"
	"github.com/jphastings/corviator/pkg/tracker"
	"periph.io/x/periph/host"
)

func main() {
	environment, err := env.ParseEnv()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	var callbacks []tracker.OnTracked

	if environment.UseLog {
		callbacks = append(callbacks, l10n.TrackerCallback)
		fmt.Println("Tracking with log")
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
		fmt.Println("Tracking with stepper motors")
	}

	if environment.UseTTS {
		ttsEngine, err := googletts.New()
		if err != nil {
			log.Fatal(err)
		}

		callbacks = append(callbacks, tts.TrackedCallback(ttsEngine))
		fmt.Println("Tracking with text-to-speech engine")
	}

	track := tracker.New(environment.Home, callbacks...)

	go track.Track()

	fmt.Printf("Corviator is ready. Home is (%.2f,%.2f), %.0fm above sea level.\n", environment.Home.Latitude, environment.Home.Longitude, environment.Home.Altitude)
	http.CorviatorAPI(environment.Port, track)
}
