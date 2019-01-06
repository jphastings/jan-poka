package main

import (
	"fmt"
	"log"

	"github.com/jphastings/corviator/pkg/env"
	"github.com/jphastings/corviator/pkg/http"
	"github.com/jphastings/corviator/pkg/l10n"
	"github.com/jphastings/corviator/pkg/tracker"
)

var callbacks []tracker.OnTracked
var environment env.Config

func init() {
	var err error
	environment, err = env.ParseEnv()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if environment.UseLog {
		callbacks = append(callbacks, l10n.TrackerCallback)
		fmt.Println("Debug log tracking: on")
	} else {
		fmt.Println("Debug log tracking: off")
	}

	track := tracker.New(environment.Home, callbacks...)

	go track.Track()

	fmt.Printf("Corviator is ready. Home is (%.2f,%.2f), %.0fm above sea level.\n", environment.Home.Latitude, environment.Home.Longitude, environment.Home.Altitude)
	http.CorviatorAPI(environment.Port, track)
}
