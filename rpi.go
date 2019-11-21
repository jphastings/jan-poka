// +build rpi

package main

import (
	"github.com/jphastings/jan-poka/pkg/pointer/tower"
	"github.com/jphastings/jan-poka/pkg/tracker"
	"log"
	"os"
	"os/signal"

	"periph.io/x/periph/host"
	_ "periph.io/x/periph/host/rpi"
)

func init() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	configurables = append(configurables, configurable{
		"Tower tracking",
		func() bool { return environment.UseTower },
		configureTower,
	})
}

func configureTower() (tracker.OnTracked, error) {
	towerConfig, err := tower.New(environment.Facing)
	if err != nil {
		return nil, err
	}

	// Hack. Fix this.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		towerConfig.Shutdown()
		os.Exit(0)
	}()

	return towerConfig.TrackerCallback, nil
}
