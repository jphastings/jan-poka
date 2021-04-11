package main

import (
	"github.com/jphastings/jan-poka/pkg/common"
	"github.com/jphastings/jan-poka/pkg/locator/instagram"
)

func init() {
	configurables = append(configurables, configurable{
		"Provider: Instagram recent post locations available.",
		func() bool { return environment.InstagramUsername != "" && environment.InstagramPassword != "" },
		configureInstagram,
	})
}

func configureInstagram() (common.OnTracked, error) {
	return nil, instagram.Login(environment.InstagramUsername, environment.InstagramPassword)
}
