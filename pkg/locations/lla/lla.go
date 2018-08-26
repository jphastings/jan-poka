package lla

import (
	"github.com/jphastings/corviator/pkg/transforms"
)

const TYPE = "lla"

type locationProvider struct {
	target transforms.LLACoords
}

type params struct {
	Φ float64 `json:"lat"`
	Λ float64 `json:"long"`
	A float64 `json:"alt"`
}

func NewLocationProvider() *locationProvider {
	return &locationProvider{}
}

func (lp *locationProvider) SetParams(decodeInto func(interface{}) error) error {
	loc := &params{}
	err := decodeInto(loc)
	if err == nil {
		lp.target = transforms.LLACoords{Φ: loc.Φ, Λ: loc.Λ, A: loc.A}
	}
	return err
}

func (lp *locationProvider) Location() (transforms.LLACoords, bool) {
	return lp.target, true
}
