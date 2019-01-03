package main

import (
	"fmt"
	"github.com/jphastings/corviator/pkg/env"
	"github.com/jphastings/corviator/pkg/hardware/motor"
	"github.com/jphastings/corviator/pkg/l10n"
	"github.com/jphastings/corviator/pkg/sphere"
	//"github.com/jphastings/corviator/pkg/tts"
	//"github.com/jphastings/corviator/pkg/tts/googletts"
	"log"
	"periph.io/x/periph/host/rpi"

	"github.com/jphastings/corviator/pkg/http"
	"github.com/jphastings/corviator/pkg/tracker"
	"periph.io/x/periph/host"

	. "github.com/jphastings/corviator/pkg/math"
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
		callbacks = append(callbacks, setupLogger())
	}
	if environment.UseSteppers {
		callbacks = append(callbacks, setupSphere(environment))
	}
	//if environment.UseTTS {
	//	callbacks = append(callbacks, setupTTS())
	//}

	track := tracker.New(environment.Home, callbacks...)

	go track.Track()

	log.Printf("Corviator is ready. Home is (%.2f,%.2f) %.0fm above sea level.\n", environment.Home.Latitude, environment.Home.Longitude, environment.Home.Altitude)
	http.CorviatorAPI(environment.Port, track)
}

func setupLogger() tracker.OnTracked {
	return func(name string, bearing AERCoords, _ bool) chan error {
		promise := make(chan error)
		fmt.Printf(l10n.Phrase(name, bearing, false))
		promise <- nil
		return promise
	}
}

func setupSphere(env env.Config) tracker.OnTracked {
	return sphere.New(
		[]*motor.Motor{
			motor.New(0, rpi.P1_7, rpi.P1_11),
		},
		rpi.P1_13, env.MotorAutoSleepLeeway,
		env.MotorSteps,
		float64(env.SphereDiameter/env.OmniwheelDiameter),
		env.MinStepInterval,
		env.Heading,
	).TrackerCallback
}

//func setupTTS() tracker.OnTracked {
//	ttsEngine, err := googletts.New()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	return tts.TrackedCallback(ttsEngine)
//}
