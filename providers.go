package main

import (
	"github.com/jphastings/jan-poka/pkg/locator/instagram"
	"github.com/jphastings/jan-poka/pkg/tracker"
)

func init() {
	configurables = append(configurables, configurable{
		"Provider: Instagram recent post locations available.",
		func() bool { return environment.InstagramUsername != "" && environment.InstagramPassword != "" },
		configureInstagram,
	})
}

func configureInstagram() (tracker.OnTracked, error) {
	return nil, instagram.Login(environment.InstagramUsername, environment.InstagramPassword)
}
