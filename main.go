package main

import (
	"fmt"
	"github.com/jphastings/jan-poka/pkg/common"
	"github.com/jphastings/jan-poka/pkg/l10n"
	"log"

	"github.com/jphastings/jan-poka/pkg/env"
	"github.com/jphastings/jan-poka/pkg/http"
	"github.com/jphastings/jan-poka/pkg/tracker"
)

var callbacks []common.OnTracked
var environment env.Config

type configurable struct {
	name      string
	toggle    func() bool
	configure func() (common.OnTracked, error)
}

var configurables = []configurable{
	{"Logging", func() bool { return environment.UseLog }, loggingCallback},
}

func init() {
	var err error
	environment, err = env.ParseEnv()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	callbacks := configureModules()
	track := tracker.New(environment.Home, callbacks)

	go track.Track()

	fmt.Printf("Jan Poka is ready. Home is (%.2f,%.2f), %.0fm above sea level.\n", environment.Home.Latitude, environment.Home.Longitude, environment.Home.Altitude)
	http.WebAPI(environment.Port, track)
}

func loggingCallback() (common.OnTracked, error) {
	return l10n.TrackerCallback, nil
}

func configureModules() map[string]common.OnTracked {
	callbacks := make(map[string]common.OnTracked)
	for _, conf := range configurables {
		if conf.toggle() {
			callback, err := conf.configure()

			if err == nil {
				if callback != nil {
					callbacks[conf.name] = callback
				}
				log.Printf("✅ %s\n", conf.name)
			} else {
				log.Fatalf("🛑 %s: \n%v", conf.name, err)
			}
		} else {
			log.Printf("✋ %s\n", conf.name)
		}
	}
	return callbacks
}
