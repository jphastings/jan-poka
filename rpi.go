// +build rpi

package main

import (
	"github.com/jphastings/jan-poka/pkg/pointer/tower"
	"github.com/jphastings/jan-poka/pkg/tracker"
	"log"

	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
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
	bus, err := i2creg.Open("")

	towerConfig, err := tower.New(bus, environment.Facing)
	if err != nil {
		return nil, err
	}

	return towerConfig.TrackerCallback(), nil
}
