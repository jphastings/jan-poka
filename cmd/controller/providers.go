package main

import (
	. "github.com/jphastings/jan-poka/pkg/common"

	"github.com/jphastings/jan-poka/pkg/locator/adsb"
	_ "github.com/jphastings/jan-poka/pkg/locator/celestial"
	_ "github.com/jphastings/jan-poka/pkg/locator/deliveroo"
	"github.com/jphastings/jan-poka/pkg/locator/google"
	"github.com/jphastings/jan-poka/pkg/locator/instagram"
	_ "github.com/jphastings/jan-poka/pkg/locator/iss"
	_ "github.com/jphastings/jan-poka/pkg/locator/lla"
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
	return nil, adsb.Connect(environment.Dump1090Host)
}
