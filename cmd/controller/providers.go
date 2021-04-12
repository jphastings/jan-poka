package main

import (
	"github.com/jphastings/jan-poka/pkg/common"
	"github.com/jphastings/jan-poka/pkg/locator/google"
	"github.com/jphastings/jan-poka/pkg/locator/instagram"
)

func init() {
	configurables = append(configurables, configurable{
		"Provider: Instagram recent post locations available.",
		func() bool { return environment.InstagramUsername != "" && environment.InstagramPassword != "" },
		configureInstagram,
	})
	configurables = append(configurables, configurable{
		"Provider: Google search locations available.",
		func() bool { return environment.GoogleMapsAPIKey != "" },
		configureGoogle,
	})
}

func configureInstagram() (common.OnTracked, error) {
	return nil, instagram.Login(environment.InstagramUsername, environment.InstagramPassword)
}

func configureGoogle() (common.OnTracked, error) {
	return nil, google.Login(environment.GoogleMapsAPIKey)
}
