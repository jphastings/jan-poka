package lla

import (
	"fmt"
	"log"
	"time"

	. "github.com/jphastings/jan-poka/pkg/common"
	. "github.com/jphastings/jan-poka/pkg/math"
)

const TYPE = "lla"

var _ LocationProvider = (*locationProvider)(nil)

type locationProvider struct {
	name   string
	target LLACoords
	time   time.Time
}

type params struct {
	Name      string   `json:"name"`
	Latitude  *Degrees `json:"lat"`
	Longitude *Degrees `json:"long"`
	Altitude  Meters   `json:"alt"`
}

func init() {
	Providers[TYPE] = func() LocationProvider { return &locationProvider{} }
	log.Println("✅ Provider: Latitude/Longitude positions")
}

func (lp *locationProvider) SetParams(decodeInto func(interface{}) error) error {
	loc := &params{}
	err := decodeInto(loc)
	if err != nil {
		return err
	}

	if loc.Latitude == nil {
		return fmt.Errorf("no latitude provided")
	}
	if loc.Longitude == nil {
		return fmt.Errorf("no longitude provided")
	}

	if loc.Name == "" {
		lp.name = "That location"
	} else {
		lp.name = loc.Name
	}
	lp.target = LLACoords{Latitude: *loc.Latitude, Longitude: *loc.Longitude, Altitude: loc.Altitude}
	lp.time = time.Now()
	return nil
}

func (lp *locationProvider) Location() TargetDetails {
	return TargetDetails{
		Name:       lp.name,
		Coords:     lp.target,
		AccurateAt: lp.time,
		Final:      true,
	}
}
