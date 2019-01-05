package celestial

import (
	. "github.com/jphastings/corviator/pkg/math"
	"github.com/jphastings/corviator/pkg/math/celestial"
	"time"
)

const TYPE = "celestial"

var bodies = map[string]*locationProvider{
	"moon":    {"the Moon", celestial.Moon},
	"sun":     {"the Sun", celestial.Sun},
	"mercury": {"the planet Mercury", celestial.Mercury},
	"venus":   {"the planet Venus", celestial.Venus},
	"mars":    {"the planet Mars", celestial.Mars},
	"jupiter": {"the planet Jupiter", celestial.Jupiter},
	"saturn":  {"the planet Saturn", celestial.Saturn},
	"neptune": {"the planet Neptune", celestial.Neptune},
	"uranus":  {"the planet Uranus", celestial.Uranus},
	"pluto":   {"the dwarf planet Pluto", celestial.Pluto},
}

type locationProvider struct {
	name string
	body celestial.Body
}

type params struct {
	BodyName string `json:"body"`
}

func NewLocationProvider() *locationProvider { return &locationProvider{} }

func (lp *locationProvider) SetParams(decodeInto func(interface{}) error) error {
	req := &params{}
	err := decodeInto(req)
	if err == nil {
		// TODO: Validation
		*lp = *bodies[req.BodyName]
	}
	return err
}

func (lp *locationProvider) Location() (LLACoords, string, bool) {
	j := celestial.JulianDay(time.Now())
	return celestial.GeocentricCoordinates(lp.body, j), lp.name, true
}
