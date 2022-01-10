package main

import (
	"fmt"
	"github.com/jphastings/jan-poka/pkg/output/mqtt"
	"log"

	"github.com/jphastings/jan-poka/pkg/common"
	"github.com/jphastings/jan-poka/pkg/env"
	"github.com/jphastings/jan-poka/pkg/future"
	"github.com/jphastings/jan-poka/pkg/http"
	"github.com/jphastings/jan-poka/pkg/l10n"
	"github.com/jphastings/jan-poka/pkg/output/mapper"
	"github.com/jphastings/jan-poka/pkg/output/webmapper"
	"github.com/jphastings/jan-poka/pkg/shutdown"
	"github.com/jphastings/jan-poka/pkg/tracker"
)

var environment env.Config

type configurable struct {
	name      string
	toggle    func() bool
	configure func() (common.OnTracked, error)
}

var configurables = []configurable{
	{"Tracker: Logging", func() bool { return environment.UseLog }, configureLogging},
	{"Tracker: MQTT", func() bool { return true }, configureMQTT},
	{"Tracker: Mapper", func() bool { return environment.UseMapper }, configureMapper},
}

func init() {
	var err error
	environment, err = env.ParseEnv()
	if err != nil {
		log.Fatalf("🛑 Could not prepare environment: %s\n", err)
	}
}

func main() {
	callbacks := configureModules()
	track := tracker.New(environment.Home, callbacks)

	go track.Track()

	fmt.Printf("Jan Poka is ready. Home is (%.2f,%.2f), %.0fm above sea level.\n", environment.Home.Latitude, environment.Home.Longitude, environment.Home.Altitude)
	http.WebAPI(environment.Port, track, environment.UseMapper)

	shutdown.Await()
}

func configureLogging() (common.OnTracked, error) {
	return l10n.TrackerCallback, nil
}

// configureMapper needs to be called before any output method that uses mapper details (eg. MQTT)
func configureMapper() (common.OnTracked, error) {
	m, err := mapper.New(environment.Persistence)
	if err != nil {
		return nil, err
	}

	return func(details common.TrackedDetails) future.Future {
		return future.All(m.TrackerCallback(details),
			webmapper.TrackerCallback(details))
	}, nil
}

func configureMQTT() (common.OnTracked, error) {
	pub, err := mqtt.New(environment.MQTTPort, environment.Persistence)
	if err != nil {
		return nil, err
	}

	return pub.TrackerCallback, nil
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
				log.Fatalf("🛑 %s: %v\n", conf.name, err)
			}
		} else {
			log.Printf("✋ %s\n", conf.name)
		}
	}
	return callbacks
}
