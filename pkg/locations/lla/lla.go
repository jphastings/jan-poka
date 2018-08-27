package lla

import (
	. "github.com/jphastings/corviator/pkg/math"
)

const TYPE = "lla"

type locationProvider struct {
	target LLACoords
}

type params struct {
	Φ Degrees `json:"lat"`
	Λ Degrees `json:"long"`
	A Meters  `json:"alt"`
}

func NewLocationProvider() *locationProvider {
	return &locationProvider{}
}

func (lp *locationProvider) SetParams(decodeInto func(interface{}) error) error {
	loc := &params{}
	err := decodeInto(loc)
	if err == nil {
		lp.target = LLACoords{Φ: loc.Φ, Λ: loc.Λ, A: loc.A}
	}
	return err
}

func (lp *locationProvider) Location() (LLACoords, bool) {
	return lp.target, true
}
