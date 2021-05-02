package main

import (
	. "github.com/jphastings/jan-poka/pkg/common"
	ads_b "github.com/jphastings/jan-poka/pkg/locator/ads-b"
	"github.com/jphastings/jan-poka/pkg/locator/google"
	"github.com/jphastings/jan-poka/pkg/locator/instagram"
)

func init() {
	configurables = append(configurables, configurable{
		"Provider: Instagram recent post locations",
		func() bool { return environment.InstagramUsername != "" && environment.InstagramPassword != "" },
		configureInstagram,
	})
	configurables = append(configurables, configurable{
		"Provider: Google search locations",
		func() bool { return environment.GoogleMapsAPIKey != "" },
		configureGoogle,
	})
	configurables = append(configurables, configurable{
		"Provider: ADS-B airplane positions",
		func() bool { return environment.Dump1090Host != "" },
		configureADSB,
	})
}

func configureInstagram() (OnTracked, error) {
	return nil, instagram.Login(environment.InstagramUsername, environment.InstagramPassword)
}

func configureGoogle() (OnTracked, error) {
	return nil, google.Login(environment.GoogleMapsAPIKey)
}

func configureADSB() (OnTracked, error) {
	return nil, ads_b.Connect(environment.Dump1090Host)
}
