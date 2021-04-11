// +build rpi

package main

import (
	"github.com/jphastings/jan-poka/pkg/hardware/stepper"
	"github.com/jphastings/jan-poka/pkg/locator/common"
	"github.com/jphastings/jan-poka/pkg/pointer/tower"
	"log"
	"os"
	"os/signal"
	"syscall"

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

func configureTower() (common.OnTracked, error) {
	stepper.SetStateDir(environment.TowerStatePath)
	towerConfig, err := tower.New(environment.Facing)
	if err != nil {
		return nil, err
	}

	// Hack. Fix this.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		towerConfig.Shutdown()
		os.Exit(0)
	}()

	return towerConfig.TrackerCallback, nil
}
