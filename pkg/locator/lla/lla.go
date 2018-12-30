package lla

import (
	. "github.com/jphastings/corviator/pkg/math"
)

const TYPE = "lla"

type locationProvider struct {
	name   string
	target LLACoords
}

type params struct {
	Name string  `json:"name"`
	Φ    Degrees `json:"lat"`
	Λ    Degrees `json:"long"`
	A    Meters  `json:"alt"`
}

func NewLocationProvider() *locationProvider {
	return &locationProvider{}
}

func (lp *locationProvider) SetParams(decodeInto func(interface{}) error) error {
	loc := &params{}
	err := decodeInto(loc)
	if err == nil {
		if loc.Name == "" {
			lp.name = "That location"
		} else {
			lp.name = loc.Name
		}
		lp.target = LLACoords{Φ: loc.Φ, Λ: loc.Λ, A: loc.A}
	}
	return err
}

func (lp *locationProvider) Location() (LLACoords, string, bool) {
	return lp.target, lp.name, true
}
